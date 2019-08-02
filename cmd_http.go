package neotest

import (
	"fmt"
	"net/http"
)

//ReqCmd http request command
type ReqCmd struct {
	*Cmd
}

func NewReqCmd(line int) *ReqCmd {
	return &ReqCmd{
		Cmd: NewCmd("req", "req <http-method> <url>", line),
	}
}

func (c *ReqCmd) Exec(vm *VM) error {
	if len(c.exprList) != 2 {
		return fmt.Errorf("num of expr must be 2, but it is %v", len(c.exprList))
	}

	if vm.CurHttpReq != nil {
		return fmt.Errorf("there is already a http request but do not be handled")
	}

	method, err := c.exprList[0].Run(vm)
	if err != nil {
		return err
	}
	url, err := c.exprList[1].Run(vm)
	if err != nil {
		return err
	}

	//just check
	_, err = http.NewRequest(method.(string), url.(string), nil)
	if err != nil {
		return err
	}

	//store req
	vm.CurHttpReq = &HttpRequest{
		Method: method.(string),
		URL:    url.(string),
	}
	return nil
}

func (c *ReqCmd) CheckExpr(varType map[string]string) error {
	if len(c.exprList) != 2 {
		return fmt.Errorf("num of expr must be 2, but it is %v", len(c.exprList))
	}

	if c.exprList[0].Type() != String {
		return fmt.Errorf("first arg type must be string")
	}

	if c.exprList[1].Type() != String {
		return fmt.Errorf("second arg type must be string")
	}
	return nil
}

//BodyCmd http request body command
type BodyCmd struct {
	*Cmd
}

func NewBodyCmd(line int) *BodyCmd {
	return &BodyCmd{
		Cmd: NewCmd("body", "body <json-data>", line),
	}
}

func (c *BodyCmd) Exec(vm *VM) error {
	if len(c.exprList) != 1 {
		return fmt.Errorf("num of expr must be 1, but it is %v", len(c.exprList))
	}

	if vm.CurHttpReq == nil {
		return fmt.Errorf("there is not a http request, please input 'req' command before")
	}

	body, err := c.exprList[0].Run(vm)
	if err != nil {
		return nil
	}

	vm.CurHttpReq.Body = body.(string)

	return nil
}

func (c *BodyCmd) CheckExpr(varType map[string]string) error {
	if len(c.exprList) != 1 {
		return fmt.Errorf("num of expr must be 1, but it is %v", len(c.exprList))
	}

	if c.exprList[0].Type() != String {
		return fmt.Errorf("arg type must be string")
	}

	return nil
}

//RetCmd return http response
type RetCmd struct {
	*Cmd
}

func NewRetCmd(line int) *RetCmd {
	return &RetCmd{
		Cmd: NewCmd("ret", "ret [<status-code>]", line),
	}
}

func (c *RetCmd) Exec(vm *VM) error {
	if len(c.exprList) != 1 {
		return fmt.Errorf("num of expr must be 1, but it is %v", len(c.exprList))
	}

	if vm.CurHttpReq == nil {
		return fmt.Errorf("there is not a http request, please input 'req' command before")
	}

	err := vm.SendHttp()
	if err != nil {
		return err
	}

	if len(c.exprList) > 0 {
		v, err := c.exprList[0].Run(vm)
		if err != nil {
			return err
		}

		actual, _ := vm.FloatV("resp.code")
		if v.(float64) != actual {
			//TODO record
		}
	}
	return nil
}

func (c *RetCmd) CheckExpr(varType map[string]string) error {
	if len(c.exprList) > 1 {
		return fmt.Errorf("num of expr must be 1 or 0, but it is %v", len(c.exprList))
	}

	if len(c.exprList) > 0 && c.exprList[0].Type() != Float {
		return fmt.Errorf("status code type must be number")
	}

	return nil
}
