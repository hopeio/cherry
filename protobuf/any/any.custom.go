package any

import (
	"github.com/hopeio/cherry/utils/encoding/json"
)

func NewAny(v interface{}) (*RawJson, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return &RawJson{Data: data}, nil
}

func BytesToJsonAny(b []byte) *RawJson {
	b = append([]byte{'"'}, b...)
	return &RawJson{Data: append(b, '"')}
}

func StringToJsonAny(s string) *RawJson {
	return &RawJson{Data: []byte("\"" + s + "\"")}
}

func (a *RawJson) MarshalJSON() ([]byte, error) {
	if len(a.Data) == 0 {
		return []byte("null"), nil
	}
	return a.Data, nil
}

func (a *RawJson) Size() int {
	return len(a.Data)
}

func (a *RawJson) MarshalTo(b []byte) (int, error) {
	return copy(b, a.Data), nil
}

func (a *RawJson) Unmarshal(b []byte) error {
	a.Data = b
	return nil
}

func (a *RawJson) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	i -= len(a.Data)
	copy(dAtA[i:], a.Data)
	return len(a.Data), nil
}
