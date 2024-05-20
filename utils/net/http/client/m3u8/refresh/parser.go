package m3u8

import (
	bytesp "bytes"
	"errors"
	"fmt"
	"github.com/hopeio/cherry/utils/crypto/aes"
	em3u8 "github.com/hopeio/cherry/utils/encoding/m3u8"
	"github.com/hopeio/cherry/utils/net/http/client"
	"io"
	"net/http"
	"net/url"
)

type Result struct {
	URL  *url.URL
	M3u8 *em3u8.M3u8
	Keys map[int]string
}

func FromURL(link string) (*Result, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, err
	}
	link = u.String()
	var body client.RawBytes
	err = client.DefaultHeaderRequest().RetryTimes(20).DisableLog().Get(link, &body)
	if err != nil {
		return nil, fmt.Errorf("request m3u8 URL failed: %s", err.Error())
	}
	m3u8, err := em3u8.Parse(body)
	if err != nil {
		return nil, err
	}
	if len(m3u8.MasterPlaylist) != 0 {
		sf := m3u8.MasterPlaylist[0]
		return FromURL(client.ResolveURL(u, sf.URI))
	}
	if len(m3u8.Segments) == 0 {
		return nil, errors.New("can not found any TS file description")
	}
	result := &Result{
		URL:  u,
		M3u8: m3u8,
		Keys: make(map[int]string),
	}

	for idx, key := range m3u8.Keys {
		switch {
		case key.Method == "" || key.Method == em3u8.CryptMethodNONE:
			continue
		case key.Method == em3u8.CryptMethodAES:
			// Request URL to extract decryption key
			keyURL := key.URI
			keyURL = client.ResolveURL(u, keyURL)
			var keyByte client.RawBytes
			err = client.DefaultHeaderRequest().RetryTimes(20).DisableLog().ResponseHandler(func(response *http.Response) (retry bool, data []byte, err error) {
				data, err = io.ReadAll(response.Body)
				if err != nil {
					return false, nil, err
				}
				if bytesp.HasPrefix(data, []byte("<html>")) {
					return true, nil, nil
				}
				if len(data) == 0 {
					return false, nil, fmt.Errorf("no key")
				}
				return false, data, err
			}).Get(keyURL, &keyByte)
			if err != nil {
				return nil, fmt.Errorf("request m3u8 URL failed: %s", err.Error())
			}
			fmt.Printf("decryption key: %s\r", string(keyByte))
			result.Keys[idx] = string(keyByte)
		default:
			return nil, fmt.Errorf("unknown or unsupported cryption method: %s", key.Method)
		}
	}
	return result, nil
}

func (r *Result) Download(segIndex int) ([]byte, error) {
	sf := r.M3u8.Segments[segIndex]

	if sf == nil {
		return nil, fmt.Errorf("invalid segment index: %d", segIndex)
	}

	tsUrl := client.ResolveURL(r.URL, sf.URI)

	var bytes client.RawBytes
	err := client.DefaultHeaderRequest().DisableLog().ResponseHandler(func(response *http.Response) (retry bool, data []byte, err error) {
		data, err = io.ReadAll(response.Body)
		if err != nil {
			return false, nil, err
		}
		if len(data) == 0 {
			return false, data, errors.New("empty response body")
		}
		if bytesp.HasPrefix(data, []byte("<html>")) {
			return true, nil, nil
		}
		return false, data, err
	}).Get(tsUrl, &bytes)
	if err != nil {
		return nil, fmt.Errorf("request %s, %s", tsUrl, err.Error())
	}

	key, ok := r.Keys[sf.KeyIndex]
	if ok && key != "" {
		bytes, err = aes.CBCDecrypt(bytes, []byte(key),
			[]byte(r.M3u8.Keys[sf.KeyIndex].IV))
		if err != nil {
			return nil, fmt.Errorf("decryt: %s, %s", tsUrl, err.Error())
		}
	}
	// https://en.wikipedia.org/wiki/MPEG_transport_stream
	// Some TS files do not start with SyncByte 0x47, they can not be played after merging,
	// Need to remove the bytes before the SyncByte 0x47(71).
	syncByte := uint8(71) //0x47
	bLen := len(bytes)
	for j := 0; j < bLen; j++ {
		if bytes[j] == syncByte {
			bytes = bytes[j:]
			break
		}
	}
	return bytes, err
}
