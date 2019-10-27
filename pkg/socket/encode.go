package socket

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"regexp"
)

var _magicKey string = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

func randomByte() byte {
	return byte(rand.Intn(255))
}

func createAcceptKey(key string) string {
	hasher := sha1.New()
	hasher.Write([]byte(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func getSocketKey(b *[]byte) string {
	reg, _ := regexp.Compile(`Sec-WebSocket-Accept: [A-z0-9+\/]+`)
	return reg.FindString(string(*b))
}

func maskMessage(msg string, mask []byte) []byte {
	encPayload := []byte{}

	for i, c := range msg {
		fmt.Println(mask)
		encPayload = append(encPayload, byte(c)^mask[i%4])
	}

	return encPayload
}

func CreatePayload(msg string) []byte {
	payload := []byte{
		0x81,
		byte(0x80 | len(msg)),
		byte(rand.Intn(255)),
		byte(rand.Intn(255)),
		byte(rand.Intn(255)),
		byte(rand.Intn(255)),
	}

	payload = append(payload, maskMessage(msg, payload[2:6])...)
	return payload
}
