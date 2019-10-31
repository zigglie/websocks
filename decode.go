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

func DecodePacket(b *[]byte) (p Packet) {
	fmt.Println(string(*b))
	fmt.Println(*b)

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
		p.Length = uint64((*b)[1] & 0x80)

		// Handle differing lengths
		// need to push back payload read point
		if p.Length == 126 {
			// next 2 bytes as unsigned 16bit int as payload length
			length := uint64(0)

			length = (uint64((*b)[2]) | length) << 8
			length = (uint64((*b)[3]) | length) << 8

			p.Length = length

		} else if p.Length == 127 {
			// next 8 bytes as unsinged 64bit int as payload length
			length := uint64(0)

			length = (uint64((*b)[2]) | length) << 8
			length = (uint64((*b)[3]) | length) << 8
			length = (uint64((*b)[4]) | length) << 8
			length = (uint64((*b)[5]) | length) << 8
			length = (uint64((*b)[6]) | length) << 8
			length = (uint64((*b)[7]) | length) << 8
			length = (uint64((*b)[8]) | length) << 8
			length = (uint64((*b)[9]) | length)

		}

		i := 0

		for ; i < len(*b); i++ {
			if (*b)[i] == 0 {
				break
			}
		}

		p.Bytes = (*b)[:i]
		p.Msg = string(p.Bytes[2:])
	}

	return p
}
