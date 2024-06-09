package main

import (
	execi "github.com/hopeio/cherry/utils/os/exec"
)

func protoc(cmd string) {
	execi.Run(cmd)
}
