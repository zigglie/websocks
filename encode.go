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
	hasher.Write([]byte(key + _magicKey))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func getSocketKey(b *[]byte) string {
	reg, _ := regexp.Compile(`Sec-WebSocket-Accept: [A-z0-9+\/]+`)
	return reg.FindString(string(*b))
}

// maskMessage performs necessary client payload masking
// as required by RFC6455. https://tools.ietf.org/html/rfc6455#section-4.1
func maskMessage(msg string, mask []byte) []byte {
	encPayload := []byte{}

	for i, c := range msg {
		encPayload = append(encPayload, byte(c)^mask[i%4])
	}

	return encPayload
}

func createMask() []byte {
	return []byte{
		randomByte(),
		randomByte(),
		randomByte(),
		randomByte(),
	}
}

// https://tools.ietf.org/html/rfc6455#section-5.2
// The length of the "Payload data", in bytes: if 0-125, that is the
// payload length.  If 126, the following 2 bytes interpreted as a
// 16-bit unsigned integer are the payload length.  If 127, the
// following 8 bytes interpreted as a 64-bit unsigned integer (the
// most significant bit MUST be 0) are the payload length.  Multibyte
// length quantities are expressed in network byte order.  Note that
// in all cases, the minimal number of bytes MUST be used to encode
// the length, for example, the length of a 124-byte-long string
// can't be encoded as the sequence 126, 0, 124.  The payload length
// is the length of the "Extension data" + the length of the
// "Application data".  The length of the "Extension data" may be
// zero, in which case the payload length is the length of the
// "Application data".
func createLengthBytes(msg *string, b *[]byte) {
	length := len(*msg)
	// fmt.Println("lennnnnnnnnnnnngth", length)
	if length > 125 && length < 0xFFFF {
		*b = append(*b, 126|0x80)
		*b = append(*b, byte((0xFF00&length)>>8))
		*b = append(*b, byte(0xFF&length))
	} else if length > 0xFFFF {
		*b = append(*b, 127)
	} else {
		*b = append(*b, byte(length|0x80))
	}
}

func CreatePayload(msg string) []byte {
	payload := []byte{0x81}

	createLengthBytes(&msg, &payload)

	fmt.Println(payload)

	mask := createMask()

	payload = append(payload, mask...)
	payload = append(payload, maskMessage(msg, mask)...)

	return payload
}
