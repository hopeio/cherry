package mail

import (
	"errors"
	"net/smtp"
	"strings"
)

type Config struct {
	AuthType    string `comment:"CRAMMD5Auth,PLAINAuth"`
	PlainAuth   PlainAuth
	CRAMMD5Auth CRAMMD5Auth
}

type PlainAuth struct {
	Identity string
	Host     string
	Port     string
	From     string
	Password string
}

type CRAMMD5Auth struct {
	UserName string
	Secret   string
}

func (c *Config) InitBeforeInject() {

}

func (c *Config) Init() {
}

func (c *Config) Build() (smtp.Auth, error) {
	c.Init()
	if strings.ToUpper(c.AuthType) == "PLAINAUTH" {
		pc := c.PlainAuth
		return smtp.PlainAuth(pc.Identity, pc.From, pc.Password, pc.Host), nil
	}
	if strings.ToUpper(c.AuthType) == "CRAMMD5AUTH" {
		cc := c.CRAMMD5Auth
		return smtp.CRAMMD5Auth(cc.UserName, cc.Secret), nil
	}

	return nil, errors.New("邮箱配置AuthType必填")
}

type Mail struct {
	smtp.Auth
	Conf Config
}

func (m *Mail) Config() any {
	return &m.Conf
}

func (m *Mail) Init() error {
	var err error
	m.Auth, err = m.Conf.Build()
	return err
}

func (m *Mail) Close() error {
	return nil
}
