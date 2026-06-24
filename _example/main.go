/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

import (
	"github.com/hopeio/mix"
	"github.com/hopeio/mix/_example/api"
)

func main() {
	mix.NewServer(mix.WithGrpcHandler(api.GrpcRegister), mix.WithHttpHandler(api.HttpHandler)).Run()
}
