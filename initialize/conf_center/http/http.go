package http

import (
	http_fs "github.com/hopeio/cherry/utils/net/http/fs"
	http_fs_watch "github.com/hopeio/cherry/utils/net/http/fs/watch"
	"net/http"
	"time"
)

var ConfigCenter = &Http{}

type Http struct {
	Interval  time.Duration
	Url       string
	AuthBasic string
	Headers   map[string]string
}

func (cc *Http) Type() string {
	return "http"
}

// 本地配置
func (cc *Http) HandleConfig(handle func([]byte)) error {

	req, _ := http.NewRequest(http.MethodGet, cc.Url, nil)
	if cc.AuthBasic != "" {
		req.Header.Add("Authorization", cc.AuthBasic)
	}
	if cc.Headers != nil {
		for k, v := range cc.Headers {
			req.Header.Add(k, v)
		}
	}

	if cc.Interval == 0 {
		file, err := http_fs.FetchFile(req)
		if err != nil {
			return err
		}
		handle(file.Binary)
		return nil
	}

	watch := http_fs_watch.New(time.Second * cc.Interval)

	callback := func(hfile *http_fs.FileInfo) {
		handle(hfile.Binary)
	}

	return watch.Add(req, callback)
}
