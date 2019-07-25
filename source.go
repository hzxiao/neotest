package neotest

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
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

//ParseExpr parse expression
func (src *Source) ParseExpr(text string) (ExprNode, error) {
	text = strings.Trim(text, " \t\n")
	if text == "" {
		return nil, fmt.Errorf("empty expression")
	}

	var expr ExprNode
	switch {
	case text == "true" || text == "false":
		expr = &boolExpr{val: text}
	case strings.EqualFold(text, "\"") && strings.HasSuffix(text, "\""):
		expr = &stringExpr{val: strings.Trim(text, "\"")}
	case strings.HasPrefix(text, "'") && strings.HasSuffix(text, "'"):
		expr = &stringExpr{val: strings.Trim(text, "'")}
	case strings.HasPrefix(text, "`") && strings.HasSuffix(text, "`"):
		//TODO: sub cmd expr

	case strings.HasPrefix(text, "$(") && strings.HasSuffix(text, ")"): //var

	case strings.HasPrefix(text, "@"): //ID
		return src.parseIDExpr(text)
	default:
		//TODO: check valid float format
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

	var i int
	if i = bytes.IndexByte(token, '`'); i < 0 {
		return
	}

	for j := i + 1; j < len(data); j++ {
		if data[j] == '`' {
			return j + 1, data[i : j+1], nil
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
