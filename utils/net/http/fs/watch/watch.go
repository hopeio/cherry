package watch

import (
	"crypto/md5"
	"github.com/hopeio/cherry/utils/log"
	http_fs "github.com/hopeio/cherry/utils/net/http/fs"
	"net/http"
	"time"
)

type Watch struct {
	interval time.Duration
	done     chan struct{}
	handler  Handler
}

type Callback struct {
	req         *http.Request
	lastModTime time.Time
	callback    func(file *http_fs.FileInfo)
	md5value    [16]byte
}

type Handler map[string]*Callback

func New(interval time.Duration) *Watch {
	w := &Watch{
		interval: interval,
		done:     make(chan struct{}, 1),
		//1.map和数组做取舍
		handler: make(map[string]*Callback),
		//handler:  make(map[string]map[fsnotify.Op]func()),
		//2.提高时间复杂度，用event做key，然后每次事件循环取值
		//handler:  make(map[fsnotify.Event]func()),
	}

	go w.run()

	return w
}

func (w *Watch) Add(req *http.Request, callback func(file *http_fs.FileInfo)) error {
	c := &Callback{
		req:      req,
		callback: callback,
	}

	c.Do(req)
	w.handler[req.RequestURI] = c
	return nil
}

func (w *Watch) AddGet(url string, callback func(file *http_fs.FileInfo)) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	return w.Add(req, callback)
}

func (w *Watch) Remove(url string) error {
	delete(w.handler, url)
	return nil
}

func (w *Watch) run() {
	timer := time.NewTicker(w.interval)
OuterLoop:
	for {
		select {
		case <-timer.C:
			for _, callback := range w.handler {
				callback.Do(callback.req)
			}
		case <-w.done:
			break OuterLoop
		}
	}
	timer.Stop()
}

func (w *Watch) Close() {
	w.done <- struct{}{}
	close(w.done)
}

func (c *Callback) Do(r *http.Request) {
	file, err := http_fs.FetchFile(r)
	if err != nil {
		log.Error(err)
		return
	}
	if !file.ModTime().IsZero() {
		if file.ModTime().After(c.lastModTime) {
			c.lastModTime = file.ModTime()
			c.callback(file)
		}
		return
	}
	md5value := md5.Sum(file.Binary)
	if md5value != c.md5value {
		c.md5value = md5value
		c.lastModTime = file.ModTime()
		c.callback(file)
	}
}

func (w *Watch) Update(interval time.Duration) {
	w.done <- struct{}{}
	w.interval = interval
	go w.run()
}
