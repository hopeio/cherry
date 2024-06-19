package luosimao

import (
	"errors"
	"github.com/hopeio/cherry/utils/net/http/client"
	"net/http"
)

var Error = errors.New("人机识别验证失败")

type Result struct {
	Error int    `json:"error"`
	Res   string `json:"res"`
	Msg   string `json:"msg"`
}

func (l *Result) CheckError() error {
	if l.Res != "success" {
		return Error
	}
	return nil
}

// Verify 对前端的验证码进行验证
func Verify(reqURL, apiKey, response string) error {
	if reqURL == "" || apiKey == "" {
		// 没有配置LuosimaoAPIKey的话，就没有验证码功能
		return nil
	}
	if response == "" {
		return Error
	}

	req := struct {
		ApiKey   string `json:"api_key"`
		Response string `json:"response"`
	}{
		ApiKey:   apiKey,
		Response: response,
	}
	result := new(Result)

	err := client.NewRequest().
		ContentType(client.ContentTypeForm).Do(http.MethodPost, reqURL, &req, result)
	if err != nil {
		return err
	}
	return nil
}
