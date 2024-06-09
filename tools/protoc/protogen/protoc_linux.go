package main

import (
	execi "github.com/hopeio/cherry/utils/os/exec"
	"strings"
)

// 直接执行protoc linux下会报 Could not make proto path relative: /xxx/*.proto: No such file or directory,找不到原因，无解
func protoc(plugins []string, file, mod, modDir string) {
	cmd := "bash -c \"protoc " + config.include + " " + file
	var args string
	for _, plugin := range plugins {
		genpath := config.genpath

		if strings.HasPrefix(plugin, "gql_out") {
			genpath += "/apidoc/" + mod
		}
		if strings.HasPrefix(plugin, "openapiv2_out") {
			plugin += mod
			genpath += "/apidoc/" + modDir
		}
		args += " --" + plugin + ":" + genpath + "\""

	}
	execi.Run(cmd + args)
}
