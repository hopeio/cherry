package kafka

import (
	"github.com/IBM/sarama"
	"github.com/hopeio/cherry/utils/log"
)

type ConsumerConfig Config

func (c *ConsumerConfig) InitBeforeInject() {
}

func (c *ConsumerConfig) InitAfterInject() {
	(*Config)(c).Init()
}

func (c *ConsumerConfig) Build() sarama.Consumer {
	c.InitAfterInject()
	consumer, err := sarama.NewConsumer(c.Addrs, c.Config)
	if err != nil {
		log.Fatal(err)
	}
	return consumer

}

type Consumer struct {
	sarama.Consumer
	Conf ConsumerConfig
}

func (c *Consumer) Config() any {
	c.Conf.Config = sarama.NewConfig()
	return &c.Conf
}

func (c *Consumer) SetEntity() {
	c.Consumer = c.Conf.Build()
}

func (c *Consumer) Close() error {
	if c.Consumer == nil {
		return nil
	}
	return c.Consumer.Close()
}
