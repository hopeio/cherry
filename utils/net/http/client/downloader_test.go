package client

import "testing"

func TestFetch(t *testing.T) {
	_, err := GetFile("")
	if err != nil {
		t.Log(err)
	}
}
