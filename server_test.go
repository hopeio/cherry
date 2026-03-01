/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package cherry

import (
	"log"
	"net/http/httptest"
	"reflect"
	"testing"
	"unsafe"

	"go.uber.org/zap"
)

type Foo struct {
	A func()
	B func()
	C func()
	D func()
}

func TestServer(t *testing.T) {
	s := Server{}
	typ := reflect.TypeOf(&s).Elem()
	log.Println(typ.Size())
	f := Foo{}
	typ = reflect.TypeOf(&f).Elem()
	log.Println(typ.Size())
}

func TestPtr(t *testing.T) {
	recorder := httptest.NewRecorder()
	recorder.WriteHeader(200)
	log.Println(recorder.Code)
	//字节对齐了,recorder size 56 bool size 1
	*(*bool)(unsafe.Pointer(uintptr(unsafe.Pointer(recorder)) + 48)) = false
	recorder.WriteHeader(300)
	log.Println(recorder.Code)
}

func TestSize(t *testing.T) {
	t.Log(reflect.TypeOf(zap.Logger{}).Size())
	t.Log(reflect.TypeOf(zap.Field{}).Size())
}
