/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package runtime

import (
	"context"
	"github.com/hopeio/protobuf/response"
	"google.golang.org/protobuf/proto"
	"net/http"
)

func HttpResponseHook(ctx context.Context, writer http.ResponseWriter, message proto.Message) error {
	if res, ok := message.(*response.HttpResponse); ok {
		hlen := len(res.Header)
		for i := 0; i < hlen && i+1 < hlen; i += 2 {
			writer.Header().Set(res.Header[i], res.Header[i+1])
		}
		writer.WriteHeader(int(res.StatusCode))
		writer.Write(res.Body)
	}
	/*	if message == nil{
		*(&message) = &response.TinyRep{}
	}*/
	return nil
}
