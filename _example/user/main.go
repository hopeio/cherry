package main

import (
	"github.com/hopeio/cherry"
	"github.com/hopeio/cherry/_example/user/api"
)

func main() {
	cherry.NewServer(cherry.WithGrpcHandler(api.GrpcRegister), cherry.WithGinHandler(api.GinRegister)).Run()

}
