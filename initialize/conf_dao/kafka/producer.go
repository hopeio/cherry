package kafka

import (
	"github.com/IBM/sarama"
	"github.com/hopeio/cherry/utils/log"
)

type ProducerConfig Config

func (c *ProducerConfig) InitBeforeInject() {
}
func (c *ProducerConfig) Init() {
	(*Config)(c).Init()
}

func (c *ProducerConfig) Build() sarama.SyncProducer {
	c.Init()
	// 使用给定代理地址和配置创建一个同步生产者
	producer, err := sarama.NewSyncProducer(c.Addrs, c.Config)
	if err != nil {
		log.Fatal(err)
	}
	return producer

}

type Producer struct {
	sarama.SyncProducer
	Conf ProducerConfig
}

func (p *Producer) Config() any {
	p.Conf.Config = sarama.NewConfig()
	return &p.Conf
}

func (p *Producer) Set() {
	p.SyncProducer = p.Conf.Build()
}

func (p *Producer) Close() error {
	if p.SyncProducer == nil {
		return nil
	}
	return p.SyncProducer.Close()
}
