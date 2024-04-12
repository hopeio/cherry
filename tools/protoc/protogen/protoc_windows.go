package main

import (
	execi "github.com/hopeio/cherry/utils/os/exec"
	"strings"
)

func protoc(plugins []string, file string) {
	for _, plugin := range plugins {
		arg := "protoc " + include + " " + file + " --" + plugin + ":" + genpath
		if strings.HasPrefix(plugin, "openapiv2_out") || strings.HasPrefix(plugin, "gql_out") {
			arg = arg + "/api"
		}
		execi.Run(arg)
	}
}
