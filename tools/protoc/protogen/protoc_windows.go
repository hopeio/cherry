package main

import (
	execi "github.com/hopeio/cherry/utils/os/exec"
	"strings"
)

func protoc(plugins []string, file, mod, modDir string) {
	cmd := "protoc " + config.include + " " + file
	var args string

	for _, plugin := range plugins {
		genpath := config.genpath
		if strings.HasPrefix(plugin, "openapiv2_out") {
			plugin += mod
			genpath += "/api/" + modDir
		}

		if strings.HasPrefix(plugin, "gql_out") {
			genpath += "/api/"
		}
		args += " --" + plugin + ":" + genpath

	}
	execi.Run(cmd + args)
}
