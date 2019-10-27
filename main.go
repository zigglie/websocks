package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"time"

	"github.com/Zigglie/websocks/pkg/socket"
)

func exitError(e error) {
	fmt.Println(e.Error())
	os.Exit(1)
}

func hasRNRN(b *[]byte) bool {
	for i, by := range *b {
		if by == 13 {
			if i < len(*b)-5 {
				if by == 13 && (*b)[i+1] == 10 && (*b)[i+2] == 13 && (*b)[i+3] == 10 {
					return true
				}
			}
		}
	}
	return false
}

func goReadBytes(c chan []byte, s *socket.Socket) {
	for {
		tmp := make([]byte, 1024)
		n, err := s.GetConn().Read(tmp)

		if err != nil {
			exitError(err)
		}

		if n != 0 {
			c <- tmp
		}
	}
}

func readBytes(r *tls.Conn, b *[]byte) {
	tmp := make([]byte, 1024)
	for {
		n, err := r.Read(tmp)
		fmt.Println(n)

		if n == 0 {
			break
		}

		if err != nil {
			exitError(err)
		}
		*b = append(*b, tmp...)

		if hasRNRN(b) {
			break
		}
	}
}

func goSendBytes(t string, s *socket.Socket) {
	i := 0
	for {
		msg := fmt.Sprintf("%s %d", t, i)
		time.Sleep(1 * time.Second)
		fmt.Println(msg)
		s.SendMessage(msg)
		i++
	}
}

func main() {
	s, err := socket.NewSocket("wss://echo.websocket.org")

	if err != nil {
		exitError(err)
	}

	c := make(chan []byte)
	go goReadBytes(c, s)

	go goSendBytes("Bytes", s)

	for {
		tmp := <-c
		pck := socket.DecodePacket(&tmp)
		fmt.Println(pck.Msg)
	}
}
