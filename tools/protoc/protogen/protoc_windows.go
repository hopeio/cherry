package main

import (
	execi "github.com/hopeio/cherry/utils/os/exec"
	"strings"
)

func protoc(plugins []string, file string) {
	cmd := "protoc " + config.include + " " + file
	var args string
	for _, plugin := range plugins {
		args += " --" + plugin + ":" + config.genpath
		if strings.HasPrefix(plugin, "openapiv2_out") || strings.HasPrefix(plugin, "gql_out") {
			args += "/api"
		}

	}
	execi.Run(cmd + args)
}
