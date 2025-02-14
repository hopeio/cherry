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
		header := writer.Header()
		for k, v := range res.Headers {
			header.Set(k, v)
		}
		writer.WriteHeader(int(res.Status))
		writer.Write(res.Body)
	}
	/*	if message == nil{
		*(&message) = &response.TinyRep{}
	}*/
	return nil
}
