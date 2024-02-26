package server

import "net"

type ConnEntry struct {
	Network    string
	ClientConn net.Conn
	DstAddr    string
	ClientResp []byte
}

type Server interface {
	StartListen()
	Deliver(conn net.Conn) (*ConnEntry, error)
	GetConnEntry() *ConnEntry
	StopListen()
}
