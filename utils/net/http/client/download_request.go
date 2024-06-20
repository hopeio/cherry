package client

import (
	"bytes"
	"context"
	"fmt"
	ioi "github.com/hopeio/cherry/utils/io"
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

var DefaultDownloader = NewDownloader()

type DownloadMode uint16

const (
	DModeOverwrite DownloadMode = 1 << iota
	DModeContinue
	DModeMultipart   // TODO 分块下载后合并
	DModeMultiThread // TODO 暂时没找到并发写文件的方法，可以并发下载,顺序写入
)

const RangeFormat = "bytes=%d-%d/%d"
const RangeFormat2 = "bytes=%d-%d/*"

type DownloadReq struct {
	Url        string
	downloader *Downloader
	ctx        context.Context
	headers    Header       //请求级请求头
	mode       DownloadMode // 模式，0-强制覆盖，1-不存在下载，2-断续下载
}

func NewDownloadReq(url string) *DownloadReq {
	return &DownloadReq{
		ctx:        context.Background(),
		Url:        url,
		downloader: DefaultDownloader,
	}
}

func (req *DownloadReq) WithDownloader(c *Downloader) *DownloadReq {
	req.downloader = c
	req.downloader.req = req
	return req
}

func (req *DownloadReq) SetDownloader(set func(c *Downloader)) *DownloadReq {
	req.downloader = NewDownloader()
	req.downloader.req = req
	set(req.downloader)
	return req
}

func (req *DownloadReq) Downloader() *Downloader {
	req.downloader = NewDownloader()
	req.downloader.req = req
	return req.downloader
}

func (req *DownloadReq) AddHeader(k, v string) *DownloadReq {
	req.headers.Set(k, v)
	return req
}

func (c *DownloadReq) Mode(mode DownloadMode) *DownloadReq {
	c.mode = mode
	return c
}

func (c *DownloadReq) GetMode() DownloadMode {
	return c.mode
}

// 如果文件已存在，强制覆盖
func (c *DownloadReq) OverwriteMode() *DownloadReq {
	c.mode |= DModeOverwrite
	return c
}

func (c *DownloadReq) GetResponse() (*http.Response, error) {
	d := c.downloader
	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, c.Url, nil)
	if err != nil {
		return nil, err
	}

	// 如果自己设置了接受编码，http库不会自动gzip解压，需要自己处理，不加Accept-Encoding和Range头会自动设置gzip
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set(httpi.HeaderAcceptLanguage, "zh-CN,zh;q=0.9;charset=utf-8")
	req.Header.Set(httpi.HeaderConnection, "keep-alive")
	req.Header.Set(httpi.HeaderUserAgent, UserAgentChrome117)
	for i := 0; i+1 < len(c.headers); i += 2 {
		req.Header.Set(c.headers[i], c.headers[i+1])
	}
	for _, opt := range d.httpRequestOptions {
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
		resp, err = d.httpClient.Do(req)
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
			if resp.StatusCode == http.StatusRequestedRangeNotSatisfiable {
				return resp, nil
			}
			return nil, fmt.Errorf("返回错误,状态码:%d,Url:%s", resp.StatusCode, req.URL.Path)
		} else {
			return resp, nil
		}
	}
	return nil, err
}

func (c *DownloadReq) GetReader() (io.ReadCloser, error) {
	resp, err := c.GetResponse()
	if err != nil {
		return nil, err
	}
	d := c.downloader
	reader := resp.Body
	if d.resDataHandler != nil {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		data, err = d.resDataHandler(data)
		if err != nil {
			return nil, err
		}
		err = resp.Body.Close()
		if err != nil {
			return nil, err
		}
		reader = ioi.WarpCloser(bytes.NewBuffer(data))
	}
	return reader, nil
}

func (c *DownloadReq) Download(filepath string) error {
	if c.mode&DModeOverwrite == 0 && fs.Exist(filepath) {
		return nil
	}

	if c.mode&DModeContinue != 0 {
		return c.ContinuationDownload(filepath)
	}
	reader, err := c.GetReader()
	if err != nil {
		return err
	}
	err = fs.Download(filepath, reader)
	err1 := reader.Close()
	if err1 != nil {
		log.Warn("Close Reader", err1)
	}
	return err
}

func (c *DownloadReq) ContinuationDownload(filepath string) error {
	f, err := fs.OpenFile(filepath+DownloadKey, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	fileinfo, err := f.Stat()
	if err != nil {
		return err
	}

	offset := fileinfo.Size()

	c.headers = append(c.headers, httpi.HeaderRange, "bytes="+strconv.FormatInt(offset, 10)+"-")

	reader, err := c.GetReader()
	if err != nil {
		return err
	}

	written, err := io.Copy(f, reader)

	err1 := reader.Close()
	if err1 != nil {
		log.Warn("close reader error:", err1)
	}

	if err != nil && err != io.EOF {
		return err
	}
	offset += written
	err = f.Close()
	if err != nil {
		return err
	}
	return os.Rename(filepath+DownloadKey, filepath)

}

// bytes xxx-xxx/xxxx
const defaultRange = "bytes=0-8388608" // 1024*1024*8

// TODO: 利用简单任务调度实现
func (c *DownloadReq) ConcurrencyDownload(filepath string, url string, concurrencyNum int) error {
	if c.mode&DModeOverwrite == 0 && fs.Exist(filepath) {
		return nil
	}
	panic("TODO")
	return nil
}

func GetReader(url string) (io.ReadCloser, error) {
	return GetReaderWithReqOption(url, nil)
}

func GetReaderWithReqOption(url string, opts ...HttpRequestOption) (io.ReadCloser, error) {

	resp, err := NewDownloader().Options(opts...).DownloadReq(url).GetResponse()
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func Download(filepath, url string) error {
	return NewDownloadReq(url).Download(filepath)
}

func GetImage(url string) (io.ReadCloser, error) {
	return GetReaderWithReqOption(url, ImageOption)
}

func DownloadImage(filepath, url string) error {
	reader, err := GetReaderWithReqOption(url, ImageOption)
	if err != nil {
		return err
	}
	return fs.Download(filepath, reader)
}

func ImageOption(req *http.Request) {
	req.Header.Set(httpi.HeaderAccept, "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8")
}

func DownloadToDir(dir, url string) error {
	return NewDownloadReq(url).Download(dir + fs.PathSeparator + urli.PathBase(url))
}
