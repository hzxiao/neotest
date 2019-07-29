package neotest

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

//Source parse source file into commands
type Source struct {
	buf     *bufio.Scanner
	curLine int
	varType map[string]string
}

func NewSource(filename string) (*Source, error) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return newSourceByBytes(bs), nil
}

func newSourceByBytes(data []byte) *Source {
	src := &Source{}
	src.buf = bufio.NewScanner(bytes.NewBuffer(data))
	src.buf.Split(splitCmd)
	src.varType = make(map[string]string)
	return src
}

//Parse parse all command line by line
func (src *Source) Parse() ([]Commander, error) {
	var cmds []Commander
	for src.buf.Scan() {
		text := src.buf.Text()
		src.curLine++
		line := strings.Count(text, "\n")
		text = strings.Trim(text, " \t\n")
		if text == "" { //empty content line
			continue
		}
		if strings.HasPrefix(text, "#") { //annotation
			continue
		}
		cmd, err := src.ParseCmd(text)
		if err != nil {
			return nil, fmt.Errorf("line: %v, err: %v", src.curLine, err)
		}

		cmds = append(cmds, cmd)
		src.curLine += line
	}
	return cmds, nil
}

//ParseCmd parse a special cmd by one-line string
func (src *Source) ParseCmd(text string) (Commander, error) {
	scan := bufio.NewScanner(strings.NewReader(text))
	scan.Split(splitExpr)

	var cmdName string
	if scan.Scan() {
		cmdName = scan.Text()
	}
	switch cmdName {
	case "echo":
		return src.parseEchoCmd(src.curLine, scan)
	case "let":
		return src.parseLetCmd(src.curLine, scan)
	default:
		return nil, fmt.Errorf("unknown cmd: %v", cmdName)
	}
	return nil, nil
}

//parseEchoCmd parse 'echo' command
func (src *Source) parseEchoCmd(line int, buf *bufio.Scanner) (*EchoCmd, error) {
	echo := NewEchoCmd(line)
	for buf.Scan() {
		text := buf.Text()
		expr, err := src.ParseExpr(text)
		if err != nil {
			return nil, err
		}

		echo.exprList = append(echo.exprList, expr)
	}

	if len(echo.exprList) == 0 {
		return nil, fmt.Errorf("no enough arguments")
	}
	return echo, nil
}

//parseLetCmd parse 'let' command
func (src *Source) parseLetCmd(line int, buf *bufio.Scanner) (*LetCmd, error) {
	let := NewLetCmd(line)
	for buf.Scan() {
		text := buf.Text()
		expr, err := src.ParseExpr(text)
		if err != nil {
			return nil, err
		}

		let.exprList = append(let.exprList, expr)
	}

	if len(let.exprList) != 2 {
		return nil, fmt.Errorf("num of expr must be 2, but it is %v", len(let.exprList))
	}
	if let.exprList[0].Type() != Identity {
		return nil, fmt.Errorf("invaild cmd syntax: the first argument should be @ID")
	}

	second := let.exprList[1]
	var vType string
	switch second.Type() {
	case Bool, Float, String:
		if second.Type() == Bool {
			vType = "bool"
		} else if second.Type() == Float {
			vType = "float"
		} else {
			vType = "string"
		}
	case SubCommand:
		vType = second.(Resultant).ResultType()
	default:
		return nil, fmt.Errorf("invalid cmd syntax: invalid second argument type")
	}

	//check variable exist
	IDs := second.(Variate).Variables()
	for _, id := range IDs {
		if _, ok := src.varType[id]; !ok {
			return nil, fmt.Errorf("%v: %v", ErrVariableUndefine.Error(), id)
		}
	}

	//record variable type on source-parsing stage
	ID := let.exprList[0].(*IDExpr).ID
	src.varType[ID] = vType
	return let, nil
}

//ParseExpr parse expression
func (src *Source) ParseExpr(text string) (ExprNode, error) {
	text = strings.Trim(text, " \t\n")
	if text == "" {
		return nil, fmt.Errorf("empty expression")
	}

	var expr ExprNode
	switch {
	case text == "true" || text == "false": // bool expr
		expr = newBoolExpr(text)
	case strings.HasPrefix(text, "\"") && strings.HasSuffix(text, "\""): //string expr one-line
		return src.parseStringExpr(strings.Trim(text, "\""))
	case strings.HasPrefix(text, "'") && strings.HasSuffix(text, "'"): // string expr multi line
		return src.parseStringExpr(strings.Trim(text, "'"))
	case strings.HasPrefix(text, "`") && strings.HasSuffix(text, "`"): // sub command expr
		return src.parseSubCmdExpr(strings.Trim(text, "`"))
	case strings.HasPrefix(text, "$(") && strings.HasSuffix(text, ")"): //var
		ID := text[2 : len(text)-1]
		if !ValidID(ID) {
			return nil, fmt.Errorf("invalid variable: %v", text)
		}
		return src.parseExprByVqr(text)
	case strings.HasPrefix(text, "@"): //ID expr
		return src.parseIDExpr(text)
	default:
		//check valid float format
		_, err := strconv.ParseFloat(text, 64)
		if err == nil { // float expr
			return newFloatExpr(text), nil
		}
		return nil, fmt.Errorf("invalid expression: %v", text)
	}

	return expr, nil
}

func (src *Source) parseIDExpr(ID string) (ExprNode, error) {
	ID = strings.TrimPrefix(ID, "@")
	if !ValidID(ID) {
		return nil, fmt.Errorf("invalid ID: %v", ID)
	}

	return &IDExpr{ID: ID}, nil
}

func (src *Source) parseExprByVqr(IDFull string) (ExprNode, error) {
	ID, _ := isVar(IDFull)
	Type, exist := src.varType[ID]
	if !exist {
		return nil, fmt.Errorf("%v: %v", ErrVariableUndefine.Error(), ID)
	}

	var expr ExprNode
	switch Type {
	case "bool":
		expr = newBoolExpr(IDFull)
	case "string":
		expr = newStringExpr(IDFull)
	case "float":
		expr = newFloatExpr(IDFull)
	default:
		return nil, fmt.Errorf("wrong variable type: %v", Type)
	}
	return expr, nil
}

//parseStringExpr check var if text contains
func (src *Source) parseStringExpr(text string) (ExprNode, error) {
	all := regexp.MustCompile(`\$\(.*?\)`).FindAllString(text, -1)
	for _, v := range all {
		ID, yes := isVar(v)
		if !yes {
			return nil, fmt.Errorf("invalid variable: %v", v)
		}
		_, ok := src.varType[ID]
		if !ok {
			return nil, fmt.Errorf("%v: %v", ErrVariableUndefine.Error(), ID)
		}
	}

	return newStringExpr(text), nil
}

func (src *Source) parseSubCmdExpr(text string) (ExprNode, error) {
	scan := bufio.NewScanner(strings.NewReader(text))
	scan.Split(splitExpr)

	var cmdName string
	if scan.Scan() {
		cmdName = scan.Text()
	}
	switch cmdName {
	case "env":
		return src.parseEnvSubCmd(src.curLine, scan)
	default:
		return nil, fmt.Errorf("unknown cmd: %v", cmdName)
	}
	return nil, nil
}

func (src *Source) parseEnvSubCmd(line int, buf *bufio.Scanner) (*EnvSubCmd, error) {
	env := NewEnvSubCmd(line)
	for buf.Scan() {
		text := buf.Text()
		expr, err := src.ParseExpr(text)
		if err != nil {
			return nil, err
		}

		env.exprList = append(env.exprList, expr)
	}

	if len(env.exprList) != 1 {
		return nil, fmt.Errorf("num of expr must be 1, but it is %v", len(env.exprList))
	}

	return env, nil
}

func splitExpr(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanWords(data, atEOF)
	if err != nil {
		return
	}

	var findEndpoint = func(c byte) (int, []byte, bool) {
		if i := bytes.IndexByte(token, c); i >= 0 {
			for j := i + 1; j < len(data); j++ {
				if data[j] == c {
					return j + 1, data[i : j+1], true
				}
			}

			if atEOF {
				return len(data), data, true
			}

			return 0, nil, true
		}
		return advance, token, false
	}
	//find sub cmd expr
	var ok bool
	advance, token, ok = findEndpoint('`')
	if ok {
		return
	}

	advance, token, ok = findEndpoint('\'')
	if ok {
		return
	}

	//find string expr
	advance, token, _ = findEndpoint('"')
	return
}

func splitCmd(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanLines(data, atEOF)
	if err != nil {
		return
	}

	var i = 0
	if i = bytes.IndexByte(token, '\''); i < 0 {
		return
	}

	for j := i + 1; j < len(data)-1; j++ {
		if data[j] == '\'' && data[j+1] == '\n' {
			return j + 2, data[0 : j+1], nil
		}
	}

	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func ValidID(ID string) bool {
	ok, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*`, ID)
	return ok
}
