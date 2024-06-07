package binding

import (
	"errors"
	"mime/multipart"
	"net/http"
	"reflect"
)

type MultipartSource http.Request

var _ Setter = (*MultipartSource)(nil)

// TrySet tries to set a value by the multipart request with the binding a form file
func (r *MultipartSource) TrySet(value reflect.Value, field reflect.StructField, key string, opt SetOptions) (isSet bool, err error) {
	if files := r.MultipartForm.File[key]; len(files) != 0 {
		return SetByMultipartFormFile(value, field, files)
	}

	return SetByKV(value, field, FormSource(r.MultipartForm.Value), key, opt)
}

func SetByMultipartFormFile(value reflect.Value, field reflect.StructField, files []*multipart.FileHeader) (isSet bool, err error) {
	switch value.Kind() {
	case reflect.Ptr:
		switch value.Interface().(type) {
		case *multipart.FileHeader:
			value.Set(reflect.ValueOf(files[0]))
			return true, nil
		}
	case reflect.Struct:
		switch value.Interface().(type) {
		case multipart.FileHeader:
			value.Set(reflect.ValueOf(*files[0]))
			return true, nil
		}
	case reflect.Slice:
		slice := reflect.MakeSlice(value.Type(), len(files), len(files))
		isSet, err = setArrayOfMultipartFormFiles(slice, field, files)
		if err != nil || !isSet {
			return isSet, err
		}
		value.Set(slice)
		return true, nil
	case reflect.Array:
		return setArrayOfMultipartFormFiles(value, field, files)
	}
	return false, errors.New("unsupported field type for multipart.FileHeader")
}

func setArrayOfMultipartFormFiles(value reflect.Value, field reflect.StructField, files []*multipart.FileHeader) (isSet bool, err error) {
	if value.Len() != len(files) {
		return false, errors.New("unsupported len of array for []*multipart.FileHeader")
	}
	for i := range files {
		setted, err := SetByMultipartFormFile(value.Index(i), field, files[i:i+1])
		if err != nil || !setted {
			return setted, err
		}
	}
	return true, nil
}
