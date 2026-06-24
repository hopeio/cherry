/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package main

import (
	"github.com/hopeio/cherry"
	"github.com/hopeio/cherry/_example/api"
)

func main() {
	cherry.NewServer(cherry.WithGrpcHandler(api.GrpcRegister), cherry.WithHttpHandler(api.HttpHandler)).Run()
}
