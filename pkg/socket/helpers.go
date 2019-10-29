package socket

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func getHeaderKey() string {
	key := []byte{}
	for i := 0; i < 16; i++ {
		key = append(key, randomByte())
	}
	return base64.URLEncoding.EncodeToString(key)
}

func getPort(addr string) (string, bool) {
	splits := strings.Split(addr, ":")
	if splits[0] == "ws" {
		return "80", false
	}
	return "443", true
}

func getEndpoint(addr string) string {
	splits := strings.Split(addr, "/")

	if strings.Compare(addr[0:3], "ws:") == 0 ||
		strings.Compare(addr[0:4], "wss:") == 0 {
		return strings.Join(splits[3:], "")
	}

	if len(splits) > 1 {
		return splits[1]
	}
	return ""
}

func getHost(addr string) string {
	splits := strings.Split(addr, "/")
	return splits[2]
}

func createWSURI(addr string) (string, bool) {
	port, secure := getPort(addr)
	splits := strings.Split(addr, "/")
	fmt.Println(getEndpoint(addr))

	return splits[2] + ":" + port, secure
}
