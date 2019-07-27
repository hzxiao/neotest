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
	src.buf.Split(bufio.ScanLines)
	src.varType = make(map[string]string)
	return src
}

func (src *Source) Parse() ([]Commander, error) {
	var cmds []Commander
	for src.buf.Scan() {
		text := src.buf.Text()
		src.curLine++
		text = strings.Trim(text, " \t\n")
		if text == "" { //empty content line
			continue
		}
		cmd, err := src.ParseCmd(text)
		if err != nil {
			return nil, err
		}

		cmds = append(cmds, cmd)
	}
	return cmds, nil
}

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
		return nil, fmt.Errorf("line %v: unknown cmd: %v", src.curLine, cmdName)
	}
	return nil, nil
}

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

	ID := let.exprList[0].(*IDExpr).ID

	second := let.exprList[1]
	var vType string
	switch second.Type() {
	case Bool, Float, String:
		IDs := second.(Variate).Variables()
		for _, id := range IDs {
			if _, ok := src.varType[id]; !ok {
				return nil, fmt.Errorf("%v: %v", ErrVariableUndefine.Error(), ID)
			}
		}
		if second.Type() == Bool {
			vType = "bool"
		} else if second.Type() == Float {
			vType = "float"
		} else {
			vType = "string"
		}
	case SubCommand:
		//TODO: handle sub cmd expr
	default:
		return nil, fmt.Errorf("invalid cmd syntax: invalid second argument type")
	}

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
	case text == "true" || text == "false":
		expr = newBoolExpr(text)
	case strings.HasPrefix(text, "\"") && strings.HasSuffix(text, "\""):
		return src.parseStringExpr(strings.Trim(text, "\""))
	case strings.HasPrefix(text, "'") && strings.HasSuffix(text, "'"):
		return src.parseStringExpr(strings.Trim(text, "'"))
	case strings.HasPrefix(text, "`") && strings.HasSuffix(text, "`"):
		//TODO: sub cmd expr

	case strings.HasPrefix(text, "$(") && strings.HasSuffix(text, ")"): //var
		ID := text[2 : len(text)-1]
		if !ValidID(ID) {
			return nil, fmt.Errorf("invalid variable: %v", text)
		}
		return src.parseExprByVqr(text)
	case strings.HasPrefix(text, "@"): //ID
		return src.parseIDExpr(text)
	default:
		//check valid float format
		_, err := strconv.ParseFloat(text, 64)
		if err == nil {
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

func (src *Source) ReadString(delim byte) (string, error) {
	builder := strings.Builder{}
	for src.buf.Scan() {
		src.curLine++
		text := strings.Trim(src.buf.Text(), " \t\n")
		builder.WriteString(text)
		if strings.HasSuffix(text, string(delim)) {
			break
		}
	}
	return builder.String(), nil
}

func splitExpr(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanWords(data, atEOF)
	if err != nil {
		return
	}

	var findEndpoint = func(c byte) (int, []byte, bool){
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
	advance,token, ok = findEndpoint('`')
	if ok {
		return
	}

	//find string expr
	advance, token, _ = findEndpoint('"')
	return
}

func ValidID(ID string) bool {
	ok, _ := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*`, ID)
	return ok
}
