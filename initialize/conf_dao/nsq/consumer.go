package nsq

import (
	"github.com/nsqio/go-nsq"
)

type ConsumerConfig struct {
	NSQLookupdAddrs []string
	NSQdAddrs       []string
	Topic           string
	Channel         string
	*nsq.Config
}

func (c *ConsumerConfig) InitBeforeInject() {
}

func (c *ConsumerConfig) Init() {
}

func (c *ConsumerConfig) Build() (*nsq.Consumer, error) {
	c.Init()
	consumer, err := nsq.NewConsumer(c.Topic, c.Channel, c.Config)
	if err != nil {
		return nil, err
	}

	if len(c.NSQLookupdAddrs) > 0 {
		if err := consumer.ConnectToNSQLookupds(c.NSQLookupdAddrs); err != nil {
			return consumer, err
		}
	}
	if len(c.NSQdAddrs) > 0 {
		if err = consumer.ConnectToNSQDs(c.NSQdAddrs); err != nil {
			return consumer, err
		}

	}
	return consumer, nil

}

type Consumer struct {
	*nsq.Consumer
	Conf ConsumerConfig
}

func (c *Consumer) Config() any {
	c.Conf.Config = nsq.NewConfig()
	return &c.Conf
}

func (c *Consumer) Init() error {
	var err error
	c.Consumer, err = c.Conf.Build()
	return err
}

func (c *Consumer) Close() error {
	if c.Consumer != nil {
		c.Consumer.Stop()
	}
	return nil
}
