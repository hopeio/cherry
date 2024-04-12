package kafka

import (
	"github.com/IBM/sarama"
)

type Config struct {
	Addrs []string
	*sarama.Config
}

func (c *Config) Default() {
}
func (c *Config) Init() {
	c.Config.Version = sarama.V3_1_0_0
}
