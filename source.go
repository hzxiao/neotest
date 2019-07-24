package neotest

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
)

type Source struct {
	buf *bufio.Scanner
	curLine int
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

	return src
}

func (src *Source) Parse(text string) ([]Commander, error) {
	var cmds []Commander
	for src.buf.Scan() {
		text := src.buf.Text()
		src.curLine++
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
		return parseEchoCmd(src.curLine, scan)
	default:
		return nil, fmt.Errorf("line %v: unknown cmd: %v", src.curLine, cmdName)
	}
	return nil, nil
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

	for j :=i+1; j < len(data); j++ {
		if data[j] == '`' {
			return j+1, data[i:j+1], nil
		}
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}
