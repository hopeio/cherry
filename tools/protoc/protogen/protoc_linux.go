package main

import (
	execi "github.com/hopeio/cherry/utils/os/exec"
)

// 直接执行protoc linux下会报 Could not make proto path relative: /xxx/*.proto: No such file or directory,找不到原因，无解
func protoc(cmd string) {
	cmd := "bash -c \"" + cmd + "\""
	execi.Run(cmd)
}
