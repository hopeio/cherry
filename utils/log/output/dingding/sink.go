package dingding

import (
	"github.com/hopeio/cherry/utils/sdk/dingding"
	"go.uber.org/zap"
	"net/url"
)

type DingDing dingding.RobotConfig

// TODO
func (th *DingDing) Write(b []byte) (n int, err error) {
	return
}

func (th *DingDing) Sync() error {
	return nil
}

func (th *DingDing) Close() error {
	return nil
}

// dingding://${token}?sercret=${sercret}
func RegisterSink() {
	_ = zap.RegisterSink("dingding", func(url *url.URL) (sink zap.Sink, e error) {
		th := new(DingDing)
		return th, nil
	})
}
