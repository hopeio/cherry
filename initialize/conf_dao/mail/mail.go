package mail

import (
	"github.com/hopeio/cherry/utils/log"
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

func (c *Config) InitAfterInject() {
}

func (c *Config) Build() smtp.Auth {
	c.InitAfterInject()
	if strings.ToUpper(c.AuthType) == "PLAINAUTH" {
		pc := c.PlainAuth
		return smtp.PlainAuth(pc.Identity, pc.From, pc.Password, pc.Host)
	}
	if strings.ToUpper(c.AuthType) == "CRAMMD5AUTH" {
		cc := c.CRAMMD5Auth
		return smtp.CRAMMD5Auth(cc.UserName, cc.Secret)
	}
	log.Fatal("邮箱配置AuthType必填")
	return nil
}

type Mail struct {
	smtp.Auth
	Conf Config
}

func (m *Mail) Config() any {
	return &m.Conf
}

func (m *Mail) SetEntity() {
	m.Auth = m.Conf.Build()
}

func (m *Mail) Close() error {
	return nil
}
