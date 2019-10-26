package socket

type Packet struct {
	Length int
	Bytes  []byte
	Msg    string
}

func DecodePacket(b *[]byte) (p Packet) {
	p.Length = int((*b)[1])

	i := 0

	for ; i < len(*b); i++ {
		if (*b)[i] == 0 {
			break
		}
	}

	p.Bytes = (*b)[:i]
	p.Msg = string(p.Bytes[2:])

	return p
}
