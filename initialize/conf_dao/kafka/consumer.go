package kafka

import (
	"github.com/IBM/sarama"
)

type ConsumerConfig Config

func (c *ConsumerConfig) InitBeforeInject() {
}

func (c *ConsumerConfig) Init() {
	(*Config)(c).Init()
}

func (c *ConsumerConfig) Build() (sarama.Consumer, error) {
	c.Init()
	return sarama.NewConsumer(c.Addrs, c.Config)

}

type Consumer struct {
	sarama.Consumer
	Conf ConsumerConfig
}

func (c *Consumer) Config() any {
	c.Conf.Config = sarama.NewConfig()
	return &c.Conf
}

func (c *Consumer) Init() error {
	var err error
	c.Consumer, err = c.Conf.Build()
	return err
}

func (c *Consumer) Close() error {
	if c.Consumer == nil {
		return nil
	}
	return c.Consumer.Close()
}
