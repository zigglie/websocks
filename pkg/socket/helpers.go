package socket

import "encoding/base64"

func getHeaderKey() string {
	key := []byte{}
	for i := 0; i < 16; i++ {
		key = append(key, randomByte())
	}
	return base64.URLEncoding.EncodeToString(key)
}
