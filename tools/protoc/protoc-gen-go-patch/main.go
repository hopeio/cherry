package main

import (
	"fmt"
	patch2 "github.com/hopeio/cherry/tools/protoc/protoc-gen-go-patch/patch"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	log.SetFlags(0)

	if os.Getenv("PROTO_PATCH_DEBUG_LOGGING") == "" {
		log.SetOutput(io.Discard)
	}

	var plugin string

	protogen.Options{
		ParamFunc: func(name, value string) error {
			switch name {
			case "plugin":
				plugin = value
			}
			return nil // Ignore unknown params.
		},
	}.Run(func(gen *protogen.Plugin) error {
		if plugin == "" {
			s := strings.TrimPrefix(filepath.Base(os.Args[0]), "protoc-gen-")
			return fmt.Errorf("no protoc plugin specified; use 'protoc --%s_out=plugin=$PLUGIN:...'", s)
		}

		// Strip our custom param(s).
		patch2.StripParam(gen.Request, "plugin")

		// Run the specified plugin and unmarshal the CodeGeneratorResponse.
		res, err := patch2.RunPlugin(plugin, gen.Request, nil)
		if err != nil {
			return err
		}

		// Initialize a Patcher and scan source proto files.
		patcher, err := patch2.NewPatcher(gen)
		if err != nil {
			return err
		}

		// Patch the CodeGeneratorResponse.
		err = patcher.Patch(res)
		if err != nil {
			return err
		}

		// Write the patched CodeGeneratorResponse to stdout.
		return patch2.Write(res, os.Stdout)
	})
}
