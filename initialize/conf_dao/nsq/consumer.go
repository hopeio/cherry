package nsq

import (
	"log"

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

func (c *ConsumerConfig) InitAfterInject() {
}

func (c *ConsumerConfig) Build() *nsq.Consumer {
	c.InitAfterInject()
	consumer, err := nsq.NewConsumer(c.Topic, c.Channel, c.Config)
	if err != nil {
		log.Fatal(err)
	}

	if len(c.NSQLookupdAddrs) > 0 {
		if err := consumer.ConnectToNSQLookupds(c.NSQLookupdAddrs); err != nil {
			log.Fatal(err)
		}
	}
	if len(c.NSQdAddrs) > 0 {
		if err = consumer.ConnectToNSQDs(c.NSQdAddrs); err != nil {
			log.Fatal(err)
		}

	}
	return consumer

}

type Consumer struct {
	*nsq.Consumer
	Conf ConsumerConfig
}

func (c *Consumer) Config() any {
	c.Conf.Config = nsq.NewConfig()
	return &c.Conf
}

func (c *Consumer) SetEntity() {
	c.Consumer = c.Conf.Build()
}

func (c *Consumer) Close() error {
	if c.Consumer != nil {
		c.Consumer.Stop()
	}
	return nil
}
