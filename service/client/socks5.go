package client

import (
	"fmt"
	"io"
	"log"
	"net"

	"golang.org/x/net/proxy"
)

type socks5Client struct {
	name   string
	dialer *proxy.Dialer
}

func NewSocks5Client(raddr string) Client {
	c := &socks5Client{
		name: "socks5Client",
	}
	dialer, err := proxy.SOCKS5("tcp", raddr, nil, proxy.Direct)
	if err != nil {
		log.Printf("Error NewSocks5Client in %s: %v\n", c.name, err)
		return nil
	}
	c.dialer = &dialer
	return c
}

// Dial implements Client.
func (c *socks5Client) Dial(protocol, dstAddr string) (net.Conn, error) {
	if c.dialer == nil {
		return nil, fmt.Errorf("dialer is empty")
	}
	dialer := *c.dialer
	var proxyConn net.Conn
	var err error
	for i := 0; i < 3; i++ { // Retry up to 3 times
		proxyConn, err = dialer.Dial("tcp", dstAddr)
		if err != nil {
			if err == io.EOF {
				continue // If "unexpected EOF", retry the connection
			}
			return nil, err
		}
		break
	}
	return proxyConn, nil
}
