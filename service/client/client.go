package client

import (
	"net"
)

type Client interface {
	Dial(network, dstAddr string) (net.Conn, error)
}
