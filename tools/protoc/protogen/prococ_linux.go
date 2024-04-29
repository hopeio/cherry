package main

import (
	execi "github.com/hopeio/cherry/utils/os/exec"
	"strings"
)

// 直接执行protoc linux下会报 Could not make proto path relative: /xxx/*.proto: No such file or directory,找不到原因，无解
func protoc(plugins []string, file string) {
	cmd := "bash -c \"protoc " + config.include + " " + file
	var args string
	for _, plugin := range plugins {
		args += " --" + plugin + ":" + config.genpath
		if strings.HasPrefix(plugin, "openapiv2_out") || strings.HasPrefix(plugin, "gql_out") {
			args += "/api"
		}
		arg += "\""

	}
	execi.Run(cmd + args)
}
