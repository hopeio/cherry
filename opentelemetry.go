/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	_ "github.com/hopeio/gox/net/http"
	"go.opentelemetry.io/otel"
)

var(
	tracer  = otel.Tracer("server")
	meter   = otel.Meter("server")
)