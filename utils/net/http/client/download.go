package client

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hopeio/cherry/utils/io/fs"
	"github.com/hopeio/cherry/utils/log"
	httpi "github.com/hopeio/cherry/utils/net/http"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type DownloadMode uint8

const (
	DModeForceOverwrite DownloadMode = iota
	DModeNotExistDownload
	DModeContinueDownload
	DModeMultipartDownload // TODO
)

// TODO: Range Status(206) PartialContent 下载
type Downloader struct {
	Client          *http.Client
	Request         *http.Request
	Mode            DownloadMode // 模式，0-强制覆盖，1-不存在下载，2-断续下载
	ResponseHandler func(response []byte) ([]byte, error)
}

func NewDownloader(url string) (*Downloader, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// 如果自己设置了接受编码，http库不会自动gzip解压，需要自己处理，不加Accept-Encoding和Range头会自动设置gzip
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set(httpi.HeaderAcceptLanguage, "zh-CN,zh;q=0.9;charset=utf-8")
	req.Header.Set(httpi.HeaderConnection, "keep-alive")
	req.Header.Set(httpi.HeaderUserAgent, UserAgentChrome117)
	return &Downloader{
		Client:  defaultClient,
		Request: req,
	}, nil
}

func (d *Downloader) WithClient(c *http.Client) *Downloader {
	d.Client = c
	return d
}

func (d *Downloader) SetClient(set func(*http.Client)) *Downloader {
	set(d.Client)
	return d
}

func (d *Downloader) WithRequest(c *http.Request) *Downloader {
	d.Request = c
	return d
}

func (d *Downloader) SetRequest(set RequestOption) *Downloader {
	set(d.Request)
	return d
}

func (d *Downloader) WithOptions(opts ...RequestOption) *Downloader {
	for _, opt := range opts {
		opt(d.Request)
	}
	return d
}

func (d *Downloader) WithResponseHandler(responseHandler func(response []byte) ([]byte, error)) *Downloader {
	d.ResponseHandler = responseHandler
	return d
}

func (d *Downloader) SetHeader(header Header) *Downloader {
	return d.WithOptions(SetHeader(header))
}

func (d *Downloader) AddHeader(header, value string) *Downloader {
	return d.WithOptions(SetHeader(Header{header, value}))
}

func (d *Downloader) WithRange(begin, end string) *Downloader {
	d.Request.Header.Set(httpi.HeaderRange, "bytes="+begin+"-"+end)
	return d
}

func (d *Downloader) WithMode(mode DownloadMode) *Downloader {
	d.Mode = mode
	return d
}

// 如果文件已存在，不下载覆盖
func (d *Downloader) ExistsSkipMode() *Downloader {
	d.Mode = DModeNotExistDownload
	return d
}

func (d *Downloader) GetResponse() (*http.Response, error) {
	if d.Request == nil {
		return nil, errors.New("client 或 request 为 nil")
	}
	if d.Client == nil {
		d.Client = defaultClient
	}

	var resp *http.Response
	var err error
	for i := 0; i < 3; i++ {
		if i > 0 {
			time.Sleep(time.Second)
		}
		resp, err = d.Client.Do(d.Request)
		if err != nil {
			log.Warn(err, "url:", d.Request.URL.Path)
			if strings.HasPrefix(err.Error(), "dial tcp: lookup") {
				return nil, err
			}
			continue
		}
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			resp.Body.Close()
			if resp.StatusCode == http.StatusNotFound {
				return nil, ErrNotFound
			}
			return nil, fmt.Errorf("返回错误,状态码:%d,url:%s", resp.StatusCode, d.Request.URL.Path)
		} else {
			return resp, nil
		}
	}
	return nil, err
}

func (d *Downloader) DownloadFile(filepath string) error {
	if d.Mode == DModeNotExistDownload && fs.Exist(filepath) {
		return nil
	}

	if d.Mode == DModeContinueDownload {
		return d.ContinuationDownloadFile(filepath)
	}

	resp, err := d.GetResponse()
	if err != nil {
		return err
	}
	var reader io.Reader = resp.Body
	if d.ResponseHandler != nil {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		data, err = d.ResponseHandler(data)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(data)
	}
	err = fs.DownloadFile(filepath, reader)
	err1 := resp.Body.Close()
	if err1 != nil {
		log.Warn("Close Reader", err1)
	}
	return err
}

const DownloadKey = fs.DownloadKey

func (d *Downloader) ContinuationDownloadFile(filepath string) error {
	f, err := fs.OpenFile(filepath+DownloadKey, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	fileinfo, err := f.Stat()
	if err != nil {
		return err
	}

	offset := fileinfo.Size()
	for {
		d.Request.Header.Set(httpi.HeaderRange, "bytes="+strconv.FormatInt(offset, 10)+"-")

		resp, err := d.GetResponse()
		if err != nil {
			return err
		}

		written, err := io.Copy(f, resp.Body)

		err1 := resp.Body.Close()
		if err1 != nil {
			log.Warn("close reader error:", err1)
		}

		if err != nil {
			log.Warn("copy error:", err, ",will go on")
			offset += written
		} else {
			err = f.Close()
			if err != nil {
				return err
			}
			return os.Rename(filepath+DownloadKey, filepath)
		}

	}

}

// bytes xxx-xxx/xxxx
const defaultRange = "bytes=0-8388608" // 1024*1024*8

// TODO: 利用简单任务调度实现
func (d *Downloader) ConcurrencyDownloadFile(filepath string, concurrencyNum int) error {
	if d.Mode == 1 && fs.Exist(filepath) {
		return nil
	}
	panic("TODO")
	return nil
}

func GetFile(url string) (io.ReadCloser, error) {
	return GetFileWithReqOption(url, nil)
}

func GetFileWithReqOption(url string, opts ...RequestOption) (io.ReadCloser, error) {
	d, err := NewDownloader(url)
	if err != nil {
		return nil, err
	}
	resp, err := d.WithOptions(opts...).GetResponse()
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func DownloadFile(filepath, url string) error {
	d, err := NewDownloader(url)
	if err != nil {
		return err
	}
	return d.DownloadFile(filepath)
}

func GetImage(url string) (io.ReadCloser, error) {
	return GetFileWithReqOption(url, ImageOption)
}

func DownloadImage(filepath, url string) error {
	reader, err := GetFileWithReqOption(url, ImageOption)
	if err != nil {
		return err
	}
	return fs.DownloadFile(filepath, reader)
}

func ImageOption(req *http.Request) {
	req.Header.Set(httpi.HeaderAccept, "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
}
