package socket

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"
)

type Socket struct {
	Address  string
	Endpoint string
	Headers  map[string]string
	config   *tls.Config
	conn     net.Conn
	secure   bool
}

func (s *Socket) dial() error {
	if s.secure {
		conf := tls.Config{}
		conn, err := tls.Dial("tcp", s.Address, &conf)

		if err != nil {
			return err
		}

		s.config = &conf
		s.conn = conn
	} else {
		conn, err := net.Dial("tcp", s.Address)

		if err != nil {
			return err
		}

		s.conn = conn
	}
	s.upgrade()
	return nil
}

// NewSocket creates a connection to addr
func (s *Socket) Init(addr string) error {
	s.Address, s.secure = createWSURI(addr)
	s.Endpoint = getEndpoint(addr)

	err := s.dial()
	return err
}

func (s *Socket) getHost() string {
	split := strings.Split(s.Address, ":")
	return split[0]
}

func (s *Socket) GetConn() net.Conn {
	return s.conn
}

func (s *Socket) upgrade() {
	payload := "GET /" + s.Endpoint + " HTTP/1.1\r\n" +
		"Host: " + s.getHost() + "\r\n" +
		"Connection: Upgrade\r\n" +
		"Accept: */*\r\n" +
		"upgrade: websocket\r\n" +
		"Sec-WebSocket-Key: " + getHeaderKey() + "\r\n" +
		"Sec-WebSocket-Version: 13\r\n"

	for k, v := range s.Headers {
		payload += k + ": " + v + "\r\n"
	}

	payload += "\r\n"
	fmt.Println(payload)
	s.conn.Write([]byte(payload))
}

func (s *Socket) SendMessage(msg string) {
	payload := CreatePayload(msg)
	s.conn.Write(payload)
}
