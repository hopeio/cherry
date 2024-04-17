package main

import (
	execi "github.com/hopeio/cherry/utils/os/exec"
	"strings"
)

// 直接执行protoc linux下会报 Could not make proto path relative: /xxx/*.proto: No such file or directory,找不到原因，无解
func protoc(plugins []string, file string) {
	for _, plugin := range plugins {
		arg := "bash -c \"protoc " + config.include + " " + file + " --" + plugin + ":" + config.genpath
		if strings.HasPrefix(plugin, "openapiv2_out") || strings.HasPrefix(plugin, "gql_out") {
			arg = arg + "/api"
		}
		arg += "\""
		execi.Run(arg)
	}
}
