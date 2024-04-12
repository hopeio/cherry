package client

import "fmt"

type ResponseBody struct {
	Status  int    `json:"status"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

func CommonResponse(response interface{}) ResponseBodyCheck {
	return &ResponseBody{Data: response}
}

func (res *ResponseBody) CheckError() error {
	if res.Status != 0 {
		return fmt.Errorf("status:%d,message:%s", res.Status, res.Message)
	}
	return nil
}

type ResponseBody2 struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

func CommonResponse2(response interface{}) ResponseBodyCheck {
	return &ResponseBody2{Data: response}
}

func (res *ResponseBody2) CheckError() error {
	if res.Code != 0 {
		return fmt.Errorf("status:%d,message:%s", res.Code, res.Msg)
	}
	return nil
}

var (
	ErrNotFound = fmt.Errorf("not found")
)
