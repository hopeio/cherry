package etcd

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var ConfigCenter = &Etcd{}

type Etcd struct {
	Conf   Config
	Client *clientv3.Client
}

type Config struct {
	clientv3.Config
	Key string
}

func (e *Etcd) Type() string {
	return "etcd"
}

func (cc *Etcd) Config() any {
	return &cc.Conf
}

// TODO: 监听更改
func (e *Etcd) Handle(handle func([]byte)) error {
	var err error
	if e.Client == nil {
		e.Client, err = clientv3.New(e.Conf.Config)
		if err != nil {
			return err
		}
	}

	resp, err := e.Client.Get(context.Background(), e.Conf.Key)
	if err != nil {
		return err
	}
	handle(resp.Kvs[0].Value)
	return nil
}

func (cc *Etcd) Close() error {
	return cc.Client.Close()
}
