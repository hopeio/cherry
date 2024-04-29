// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package main

import (
	"flag"

	"github.com/hopeio/cherry/tools/protoc/protoc-gen-validator/plugin"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(p *protogen.Plugin) error {

		return plugin.New(p).Generate()
	})
}
