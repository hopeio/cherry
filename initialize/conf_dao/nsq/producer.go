package nsq

import "github.com/nsqio/go-nsq"

type ProducerConfig struct {
	Addr string
	*nsq.Config
}

func (c *ProducerConfig) InitBeforeInject() {
}
func (c *ProducerConfig) Init() {
}

func (c *ProducerConfig) Build() (*nsq.Producer, error) {
	c.Init()
	return nsq.NewProducer(c.Addr, c.Config)
}

type Producer struct {
	*nsq.Producer
	Conf ProducerConfig
}

func (p *Producer) Config() any {
	p.Conf.Config = nsq.NewConfig()
	return &p.Conf
}

func (p *Producer) Init() error {
	var err error
	p.Producer, err = p.Conf.Build()
	return err
}

func (p *Producer) Close() error {
	if p.Producer != nil {
		p.Producer.Stop()
	}
	return nil
}
