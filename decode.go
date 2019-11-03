package socket

import "fmt"

type Packet struct {
	Length uint64
	Bytes  []byte
	Msg    string
}

// *  %x0 denotes a continuation frame

// *  %x1 denotes a text frame

// *  %x2 denotes a binary frame

// *  %x3-7 are reserved for further non-control frames

// *  %x8 denotes a connection close

// *  %x9 denotes a ping

// *  %xA denotes a pong

// *  %xB-F are reserved for further control frames

const (
	CONTCODE  int = 0x0
	TEXTCODE  int = 0x1
	BINCODE   int = 0x2
	CLOSECODE int = 0x8
	PINGCODE  int = 0x9
	PONGCODE  int = 0xA
)

func getOpCode(b byte) int {
	return int(b & 0xF)
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
func getLength(b *[]byte) uint64 {
	length := uint64((*b)[1])
	if length == 127 {
		length = 0

		length = (uint64((*b)[2]) | length) << 8
		length = (uint64((*b)[3]) | length) << 8
		length = (uint64((*b)[4]) | length) << 8
		length = (uint64((*b)[5]) | length) << 8
		length = (uint64((*b)[6]) | length) << 8
		length = (uint64((*b)[7]) | length) << 8
		length = (uint64((*b)[8]) | length) << 8
		length = (uint64((*b)[9]) | length)

		return length
	} else if length == 126 {
		// Get next two bytes
		length = uint64((*b)[2]) << 8
		length = length | uint64((*b)[3])
		return length
	}
	return length
}

func DecodePacket(b *[]byte) (p Packet) {
	// opcode := getOpCode((*b)[0])

	if (*b)[0] == 0x72 {
		fmt.Println("Probably not websocket data")

		i := 0

		for ; i < len(*b); i++ {
			if (*b)[i] == 0 {
				break
			}
		}

		p.Length = uint64(i)
		p.Bytes = (*b)[:i]
		p.Msg = string(p.Bytes)

	} else {
		p.Length = getLength(b)
		push := 2

		// Handle differing lengths
		// need to push back payload read point
		if p.Length == 126 {
			// next 2 bytes as unsigned 16bit int as payload length
			push = 4
		} else if p.Length == 127 {
			push = 10
		}

		i := 0

		for ; i < len(*b); i++ {
			if (*b)[i] == 0 {
				break
			}
		}

		p.Bytes = (*b)[:i]
		p.Msg = string(p.Bytes[push:])
	}

	return p
}
