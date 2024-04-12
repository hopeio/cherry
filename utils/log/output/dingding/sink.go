package dingding

import (
	"github.com/hopeio/cherry/utils/sdk/dingding"
	"go.uber.org/zap"
	"net/url"
)

type TalkHook dingding.DingRobotConfig

func (th *TalkHook) Write(b []byte) (n int, err error) {
	return
}

func (th *TalkHook) Sync() error {
	return nil
}

func (th *TalkHook) Close() error {
	return nil
}

// dingding://${token}?sercret=${sercret}
func RegisterSink() {
	_ = zap.RegisterSink("dingding", func(url *url.URL) (sink zap.Sink, e error) {
		th := new(TalkHook)
		return th, nil
	})
}
