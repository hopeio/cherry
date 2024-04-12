package kibana

import (
	"github.com/elastic/go-elasticsearch/v8"
	"go.uber.org/zap"
	"net/url"
)

type Kibana struct {
	Es    *elasticsearch.Client
	Index string
}

func (k *Kibana) Write(b []byte) (n int, err error) {
	return
}

func (k *Kibana) Sync() error {
	return nil
}

func (k *Kibana) Close() error {
	return nil
}

// kibana://${token}?index=${index}
func RegisterSink() {
	_ = zap.RegisterSink("kibana", func(url *url.URL) (sink zap.Sink, e error) {
		k := new(Kibana)
		k.Es, e = elasticsearch.NewClient(elasticsearch.Config{})
		k.Index = url.Query().Get("index")
		return k, e
	})
}

func New(es *elasticsearch.Client, index string) zap.Sink {
	return &Kibana{es, index}
}
