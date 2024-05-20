package client

import (
	"fmt"
	"github.com/hopeio/cherry/utils/net/http/client"
)

type ResponseBody[RES any] struct {
	Status  int    `json:"status"`
	Data    RES    `json:"data"`
	Message string `json:"message"`
}

func CommonResponse[RES any]() client.ResponseBodyCheck {
	return &ResponseBody[RES]{}
}

func (res *ResponseBody[RES]) CheckError() error {
	if res.Status != 0 {
		return fmt.Errorf("status: %d, message: %s", res.Status, res.Message)
	}
	return nil
}

type ResponseBody2[RES any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data RES    `json:"data"`
}

func CommonResponse2[RES any]() client.ResponseBodyCheck {
	return &ResponseBody2[RES]{}
}

func (res *ResponseBody2[RES]) CheckError() error {
	if res.Code != 0 {
		return fmt.Errorf("code: %d, msg: %s", res.Code, res.Msg)
	}
	return nil
}

func (res *ResponseBody2[RES]) GetData() *RES {
	return &res.Data
}
