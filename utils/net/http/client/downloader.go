package client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hopeio/cherry/utils/io/fs"
	"github.com/hopeio/cherry/utils/log"
	httpi "github.com/hopeio/cherry/utils/net/http"
	urli "github.com/hopeio/cherry/utils/net/url"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type DownloadMode uint16

const (
	DModeOverwrite DownloadMode = iota
	DModeNotExist
	DModeContinue          = DModeNotExist << 1
	DModeMultipartDownload = DModeNotExist << 2 // TODO
)

// TODO: Range Status(206) PartialContent 下载
type Downloader struct {
	ctx                context.Context
	client             *http.Client
	defaultClient      bool
	proxyUrl           string
	timeout            time.Duration
	authUser, authPass string
	header             http.Header
	mode               DownloadMode // 模式，0-强制覆盖，1-不存在下载，2-断续下载
	responseHandler    func(response []byte) ([]byte, error)
	retryTimes         int
	retryInterval      time.Duration
	requestOptions     []RequestOption
}

func NewDownloader() *Downloader {
	return &Downloader{
		retryTimes:    3,
		retryInterval: time.Second,
	}
}

func (d *Downloader) HttpClient(c *http.Client) *Downloader {
	d.client = c
	return d
}

func (d *Downloader) Options(opts ...RequestOption) *Downloader {
	d.requestOptions = append(d.requestOptions, opts...)
	return d
}

func (d *Downloader) ResponseHandler(responseHandler func(response []byte) ([]byte, error)) *Downloader {
	d.responseHandler = responseHandler
	return d
}

func (d *Downloader) Header(header http.Header) *Downloader {
	d.header = header
	return nil
}

func (d *Downloader) AddHeader(header, value string) *Downloader {
	d.header.Add(header, value)
	return d
}

func (d *Downloader) Range(begin, end string) *Downloader {
	d.header.Add(httpi.HeaderRange, "bytes="+begin+"-"+end)
	return d
}

func (d *Downloader) Mode(mode DownloadMode) *Downloader {
	d.mode = mode
	return d
}

func (d *Downloader) GetMode() DownloadMode {
	return d.mode
}

// 如果文件已存在，不下载覆盖
func (d *Downloader) ExistsSkipMode() *Downloader {
	d.mode |= DModeNotExist
	return d
}

func (d *Downloader) GetResponse(url string) (*http.Response, error) {
	if d.client == nil {
		d.client = DefaultHttpClient
		d.defaultClient = true
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if d.header != nil {
		req.Header = d.header
	}
	// 如果自己设置了接受编码，http库不会自动gzip解压，需要自己处理，不加Accept-Encoding和Range头会自动设置gzip
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set(httpi.HeaderAcceptLanguage, "zh-CN,zh;q=0.9;charset=utf-8")
	req.Header.Set(httpi.HeaderConnection, "keep-alive")
	req.Header.Set(httpi.HeaderUserAgent, UserAgentChrome117)

	for _, opt := range d.requestOptions {
		opt(req)
	}

	var resp *http.Response
	if d.retryTimes == 0 {
		d.retryTimes = 1
	}
	for i := range d.retryTimes {
		if i > 0 {
			time.Sleep(d.retryInterval)
		}
		resp, err = d.client.Do(req)
		if err != nil {
			log.Warn(err, "Url:", req.URL.Path)
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
			return nil, fmt.Errorf("返回错误,状态码:%d,Url:%s", resp.StatusCode, req.URL.Path)
		} else {
			return resp, nil
		}
	}
	return nil, err
}

func (d *Downloader) Download(filepath, url string) error {
	if d.mode&DModeNotExist != 0 && fs.Exist(filepath) {
		return nil
	}

	if d.mode&DModeContinue != 0 {
		return d.ContinuationDownload(filepath, url)
	}

	resp, err := d.GetResponse(url)
	if err != nil {
		return err
	}
	var reader io.Reader = resp.Body
	if d.responseHandler != nil {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		data, err = d.responseHandler(data)
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

func (d *Downloader) ContinuationDownload(filepath, url string) error {
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
		d.header.Set(httpi.HeaderRange, "bytes="+strconv.FormatInt(offset, 10)+"-")

		resp, err := d.GetResponse(url)
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
func (d *Downloader) ConcurrencyDownload(filepath string, concurrencyNum int) error {
	if d.mode == 1 && fs.Exist(filepath) {
		return nil
	}
	panic("TODO")
	return nil
}

func GetFile(url string) (io.ReadCloser, error) {
	return GetFileWithReqOption(url, nil)
}

func GetFileWithReqOption(url string, opts ...RequestOption) (io.ReadCloser, error) {
	d := NewDownloader()

	resp, err := d.Options(opts...).GetResponse(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func Download(filepath, url string) error {
	d := NewDownloader()
	return d.Download(filepath, url)
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

func DownloadToDir(dir, url string) error {
	d := NewDownloader()
	return d.Download(dir+fs.PathSeparator+urli.PathBase(url), url)
}
