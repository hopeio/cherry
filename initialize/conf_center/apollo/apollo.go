package apollo

import (
	"encoding/json"
	"github.com/hopeio/cherry/utils/configor/apollo"
)

var ConfigCenter = &Apollo{}

type Apollo struct {
	apollo.Config
}

func (e *Apollo) Type() string {
	return "apollo"
}

// TODD: 更改监听
func (e *Apollo) HandleConfig(handle func([]byte)) error {
	client, err := e.NewClient()
	if err != nil {
		return err
	}
	config, err := client.GetDefaultConfig()
	if err != nil {
		return err
	}
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	handle(data)
	return nil
}
