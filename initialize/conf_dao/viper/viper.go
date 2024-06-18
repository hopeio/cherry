package viper

import (
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"gopkg.in/ini.v1"
	"os"
)

// 全局变量,只一个实例,只提供config
type Config struct {
	Debug             bool
	Watch             bool
	ConfigName        string
	ConfigFile        string
	ConfigType        string
	ConfigPermissions os.FileMode
	EnvPrefix         string
	RemoteProvider    []RemoteProvider
	AllowEmptyEnv     bool
	IniLoadOptions    ini.LoadOptions
	EnvVars           []string
}

type RemoteProvider struct {
	Provider      string
	Endpoint      string
	Path          string
	SecretKeyring string
}

func (c *Config) InitBeforeInject() {

}
func (c *Config) Init() {
	if c.ConfigType == "" {
		c.ConfigType = "toml"
	}
}

func (c *Config) InitAfterInject() {
	c.Init()
	c.build(viper.GetViper())
}

func (c *Config) Build() (*viper.Viper, error) {
	c.Init()
	var runtimeViper = viper.New()
	return runtimeViper, c.build(runtimeViper)
}

func (c *Config) build(runtimeViper *viper.Viper) error {
	if c.Debug {
		runtimeViper.Debug()
	}
	runtimeViper.SetConfigType(c.ConfigType) // because there is no file extension in a stream of bytes, supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "Env", "dotenv"
	if len(c.RemoteProvider) > 0 {
		var err error
		for _, conf := range c.RemoteProvider {
			if conf.SecretKeyring != "" {
				err = runtimeViper.AddSecureRemoteProvider(conf.Provider, conf.Endpoint, conf.Path, conf.SecretKeyring)
			} else {
				err = runtimeViper.AddRemoteProvider(conf.Provider, conf.Endpoint, conf.Path)
			}
			if err != nil {
				return err
			}
		}

		// read from remote Config the first time.
		err = runtimeViper.ReadRemoteConfig()
		if err != nil {
			return err
		}
		if c.Watch {
			err = runtimeViper.WatchRemoteConfig()
			if err != nil {
				return err
			}
		}

	} else {

		runtimeViper.SetConfigFile(c.ConfigFile)
		if c.ConfigPermissions > 0 {
			runtimeViper.SetConfigPermissions(c.ConfigPermissions)
		}
		err := runtimeViper.ReadInConfig()
		if err != nil {
			return err
		}
		if c.Watch {
			runtimeViper.WatchConfig()
		}
	}
	runtimeViper.AllowEmptyEnv(c.AllowEmptyEnv)
	runtimeViper.SetEnvPrefix(c.EnvPrefix)
	if len(c.EnvVars) > 0 {
		err := runtimeViper.BindEnv(c.EnvVars...)
		if err != nil {
			return err
		}
	}

	// open a goroutine to watch remote changes forever
	//这段实现不够优雅
	/*	go func() {
		for {
			time.Sleep(time.Second * 5) // delay after each request

			// currently, only tested with etcd support
			err := runtime_viper.WatchRemoteConfig()
			if err != nil {
				log.Errorf("unable to read remote Config: %v", err)
				continue
			}
			vconf :=runtime_viper.AllSettings()
			log.Debug(vconf)
			// unmarshal new Config into our runtime Config struct. you can also use channel
			// to implement a signal to notify the system of the changes
			runtime_viper.Unmarshal(cCopy)
			refresh(cCopy, dCopy)
			log.Debug(cCopy)
		}
	}()*/
	return nil
}

// 不建议使用,请使用viper全局变量
type Viper struct {
	*viper.Viper
	Conf Config
}

func (v *Viper) Config() any {
	return &v.Conf
}

func (v *Viper) Init() error {
	var err error
	v.Viper, err = v.Conf.Build()
	return err
}

func (v *Viper) Close() error {
	return v.Close()
}
