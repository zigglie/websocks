package socket

import (
	"crypto/tls"
	"net"
	"strings"
)

type Socket struct {
	Address string
	config  *tls.Config
	conn    net.Conn
	secure  bool
}

func getPort(addr string) string {
	splits := strings.Split(addr, ":")
	if splits[0] == "ws" {
		return "80"
	}
	return "443"
}

func dial(s *Socket) {

}

// NewSocket creates a connection to addr
func NewSocket(addr string) (*Socket, error) {
	port := getPort(addr)

	conf := tls.Config{}
	addr += ":" + port
	splits := strings.Split(addr, "//")
	conn, err := tls.Dial("tcp", splits[1], &conf)

	if err != nil {
		return nil, err
	}

	s := &Socket{}

	s.Address = addr
	s.config = &conf
	s.conn = conn

	s.upgrade()
	//s.conn.Handshake()
	return s, nil
}

func (s *Socket) GetConn() net.Conn {
	return s.conn
}

func (s *Socket) upgrade() {
	payload := "GET /?encoding=string HTTP/1.1\r\n" +
		"Host: echo.websocket.org\r\n" +
		"Accept: */*\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: keep-alive, Upgrade\r\n" +
		"Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==\r\n" +
		"Sec-WebSocket-Version: 13\r\n\r\n"
	s.conn.Write([]byte(payload))
}

func (s *Socket) SendMessage(msg string) {
	payload := CreatePayload(msg)
	s.conn.Write(payload)
}
