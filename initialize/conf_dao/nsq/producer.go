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

func (c *ProducerConfig) Build() *nsq.Producer {
	c.Init()
	producer, err := nsq.NewProducer(c.Addr, c.Config)
	if err != nil {
		panic(err)
	}

	return producer
}

type Producer struct {
	*nsq.Producer
	Conf ProducerConfig
}

func (p *Producer) Config() any {
	p.Conf.Config = nsq.NewConfig()
	return &p.Conf
}

func (p *Producer) Set() {
	p.Producer = p.Conf.Build()
}

func (p *Producer) Close() error {
	if p.Producer != nil {
		p.Producer.Stop()
	}
	return nil
}
