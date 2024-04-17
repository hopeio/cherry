package main

import (
	"github.com/hopeio/cherry/utils/io/fs"
	osi "github.com/hopeio/cherry/utils/os"
	execi "github.com/hopeio/cherry/utils/os/exec"
	_go "github.com/hopeio/cherry/utils/sdk/go"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

/*
*文件名正则不支持以及enum生成和model生成用的都是gogo的，所以顺序必须是gogo_out在前，enum_out在后
 */

//go:generate mockgen -destination ../protobuf/user/user.mock.go -package user -source ../protobuf/user/user.service_grpc.pb.go UserServiceServer

const (
	goOut           = "go-patch_out=plugin=go,paths=source_relative"
	grpcOut         = "go-patch_out=plugin=go-grpc,paths=source_relative"
	enumOut         = "enum_out=paths=source_relative"
	gatewayOut      = "grpc-gin_out=paths=source_relative"
	openapiv2Out    = "openapiv2_out=logtostderr=true"
	govalidatorsOut = "govalidators_out=paths=source_relative"
	gqlOut          = "gql_out=svc=true,merge=true,paths=source_relative"
	gogqlOut        = "gogql_out=svc=true,merge=true,paths=source_relative"
	dartOut         = "dart_out=grpc"
)

const (
	goListDir     = `go list -m -f {{.Dir}} `
	goListDep     = `go list -m -f {{.Path}}@{{.Version}} `
	DepGoogleapis = "github.com/googleapis/googleapis@v0.0.0-20220520010701-4c6f5836a32f"
	Depcherry     = "github.com/hopeio/cherry"
)

var (
	DepGrpcGateway = "github.com/grpc-ecosystem/grpc-gateway/v2"
	DepProtopatch  = "github.com/alta/protopatch"
)

var plugin = []string{goOut, grpcOut}

//"gqlgencfg_out=paths=source_relative",

var enumPlugin = enumOut
var gatewayPlugin = []string{gatewayOut, openapiv2Out}
var validatorsOutPlugin = govalidatorsOut
var gqlPlugin = []string{gqlOut, gogqlOut}

var (
	proto, genpath, dproto                                                          string
	include                                                                         string
	useEnumPlugin, useGateWayPlugin, useValidatorsOutPlugin, useGqlPlugin, stdPatch bool
)

func init() {
	protodef, _ := filepath.Abs("/proto")
	pwd, _ := os.Getwd()
	pflag := rootCmd.PersistentFlags()
	pflag.StringVarP(&proto, "proto", "p", protodef, "proto dir")
	pflag.StringVarP(&genpath, "genpath", "g", pwd+"/protobuf", "generate dir")
	pflag.StringVarP(&dproto, "cherry", "d", "/proto", "cherry proto dir")
	pflag.BoolVarP(&useEnumPlugin, "enum", "e", false, "是否使用enum扩展插件")
	pflag.BoolVarP(&useGateWayPlugin, "gw", "w", false, "是否使用grpc-gateway插件")
	pflag.BoolVarP(&useValidatorsOutPlugin, "validator", "v", false, "是否使用validators插件")
	pflag.BoolVarP(&useGqlPlugin, "graphql", "q", false, "是否使用graphql插件")
	pflag.BoolVar(&stdPatch, "patch", false, "是否使用原生protopatch")
	rootCmd.AddCommand(&cobra.Command{
		Use: "test",
		Run: func(cmd *cobra.Command, args []string) {
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use: "go",
		Run: func(cmd *cobra.Command, args []string) {
			run(proto)
			if useGqlPlugin {
				gengql()
			}
		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use: "dart",
		Run: func(cmd *cobra.Command, args []string) {

		},
	})
	rootCmd.AddCommand(&cobra.Command{
		Use: "ts",
		Run: func(cmd *cobra.Command, args []string) {

		},
	})

}

var rootCmd = &cobra.Command{
	Use: "protogen",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if fs.CheckNotExist(genpath) {
			os.MkdirAll(genpath, os.ModePerm)
		}
		if useEnumPlugin {
			plugin = append(plugin, enumPlugin)
		}
		if useGateWayPlugin {
			plugin = append(plugin, gatewayPlugin...)
		}
		if useValidatorsOutPlugin {
			plugin = append(plugin, validatorsOutPlugin)
		}
		if useGqlPlugin {
			plugin = append(plugin, gqlPlugin...)
		}
		getInclude()
	},
}

func main() {
	//single("/content/moment.model.proto")
	rootCmd.Execute()
	//gengql()

}

func run(dir string) {
	protoc(plugin, dir+"/*.proto")
	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalln(err)
	}
	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			run(dir + "/" + fileInfos[i].Name())
		}
	}
}
func getInclude() {
	pwd, _ := os.Getwd()
	proto, _ = filepath.Abs(proto)
	genpath, _ = filepath.Abs(genpath)
	log.Println("proto:", proto)
	log.Println("genpath:", genpath)
	if useGateWayPlugin || useGqlPlugin {
		_, err := os.Stat(genpath + "/api")
		if os.IsNotExist(err) {
			err = os.Mkdir(genpath+"/api", os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	/*	generatePath := "generate" + time.Now().Format("150405")
		err = os.Mkdir(generatePath, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		generatePath = pwd + "/" + generatePath
		defer os.RemoveAll(generatePath)
		err = os.Chdir(generatePath)
		if err != nil {
			log.Fatal(err)
		}
		osi.Cmd("go mod init generate")*/

	libcherryDir, err := osi.Cmd(_go.GoListDir + Depcherry)
	if err == nil {
		dproto = libcherryDir + "/protobuf/_proto"
	}
	include = "-I" + dproto + " -I" + proto
	/*	os.Chdir(libcherryDir)
		DepGrpcGateway, _ = osi.Cmd(goListDep + DepGrpcGateway)
		DepProtopatch, _ = osi.Cmd(goListDep + DepProtopatch)
		os.Chdir(generatePath)
		libGoogleDir := _go.GetDepDir(DepGoogleapis)
		libGatewayDir := _go.GetDepDir(DepGrpcGateway)*/

	os.Chdir(pwd)

	log.Println("include:", include)

}

// 找出所以包含go文件的包

func getPackages(dir string) []string {
	p := make(map[string]struct{})
	getPackagesHelper(dir, "", p)
	var packages []string
	for pkg, _ := range p {
		packages = append(packages, pkg)
	}
	return packages
}

func getPackagesHelper(dir, pre string, p map[string]struct{}) {
	fileInfos, err := os.ReadDir(dir)
	if err != nil {
		log.Panicln(err)
	}
	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			files, err := os.ReadDir(dir + "/" + fileInfos[i].Name())
			if err != nil {
				log.Panicln(err)
			}
			for j := range files {
				if strings.HasSuffix(files[j].Name(), ".go") {
					if pre != "" {
						p[pre+"/"+fileInfos[i].Name()] = struct{}{}
					} else {
						p[fileInfos[i].Name()] = struct{}{}
					}

					break
				}
			}
			getPackagesHelper(dir+"/"+fileInfos[i].Name(), fileInfos[i].Name(), p)
		}
	}
}
func gengql() {
	// 完整路径
	compath, err := filepath.Abs(genpath)
	if err != nil {
		log.Panicln(err)
	}
	// mod名
	err = os.Chdir(genpath)
	if err != nil {
		log.Panicln(err)
	}
	out, err := osi.Cmd("go list -m")
	if err != nil {
		log.Panicln(err)
	}
	mods := strings.Split(out, "\n")
	mod := mods[len(mods)-1]
	// 调用方mod路径
	out, err = osi.Cmd("go list -m -f {{.Dir}}")
	// 如果生成路径包含模块名
	_, after, _ := strings.Cut(compath, out)
	gomod := strings.ReplaceAll(mod+after, "\\", "/")
	packages := getPackages(compath)
	gqldir := genpath + "/api"
	fileInfos, err := os.ReadDir(gqldir)
	if err != nil {
		log.Panicln(err)
	}
	for i := range fileInfos {
		if fileInfos[i].IsDir() {
			files, err := os.ReadDir(gqldir + "/" + fileInfos[i].Name())
			if err != nil {
				log.Panicln(err)
			}
			for j := range files {
				if strings.HasSuffix(files[j].Name(), ".graphql") {
					os.Chdir(gqldir + "/" + fileInfos[i].Name())
					/*			data, err := os.ReadFile(fileInfos[i].Name() + ".graphql")
								if err != nil {
									return
								}
								dataStr := stringsi.ToString(data)
								dataStr = strings.ReplaceAll(dataStr, ": Int", ": Int!")
								dataStr = strings.ReplaceAll(dataStr, ": String", ": String!")
								dataStr = strings.ReplaceAll(dataStr, ": Boolean", ": Boolean!")
								dataStr = strings.ReplaceAll(dataStr, ": Float", ": Float!")*/

					//这里用模板生成yml
					t := template.Must(template.New("yml").Parse(ymlTpl))
					config := fileInfos[i].Name() + `.gqlgen.yml`
					_, err := os.Stat(config)
					var file *os.File
					file, err = os.Create(config)
					if err != nil {
						log.Panicln(err)
					}
					renderValue := map[string]any{"GoMod": gomod, "SubDir": fileInfos[i].Name(), "Packages": packages}
					err = t.Execute(file, renderValue)
					if err != nil {
						log.Panicln(err)
					}
					file.Close()
					execi.Run(`gqlgen --verbose --config ` + config)
					break
				}
			}
		}
	}
}
