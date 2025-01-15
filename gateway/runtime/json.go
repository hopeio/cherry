/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package runtime

import (
	"github.com/hopeio/utils/encoding/json"
	"github.com/hopeio/utils/encoding/protobuf/jsonpb"
	"io"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

var JsonPb = &JSONPb{}

type JSONPb struct {
	jsonpb.JSONPb
}

// NewDecoder returns a runtime.Decoder which reads JSON stream from "r".
func (j *JSONPb) NewDecoder(r io.Reader) runtime.Decoder {
	return json.NewDecoder(r)
}

// NewEncoder returns an Encoder which writes JSON stream into "w".
func (j *JSONPb) NewEncoder(w io.Writer) runtime.Encoder {
	return json.NewEncoder(w)
}
