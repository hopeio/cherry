package main

import (
	execi "github.com/hopeio/cherry/utils/os/exec"
	"strings"
)

func protoc(plugins []string, file string) {
	for _, plugin := range plugins {
		arg := "protoc " + config.include + " " + file + " --" + plugin + ":" + config.genpath
		if strings.HasPrefix(plugin, "openapiv2_out") || strings.HasPrefix(plugin, "gql_out") {
			arg = arg + "/api"
		}
		execi.Run(arg)
	}
}
