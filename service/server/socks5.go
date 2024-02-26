package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sysproxy/service/server/packet"
	"time"
)

type socks5erver struct {
	name     string
	addr     string
	username string
	password string
	listener *net.Listener
	connPool chan *ConnEntry
}

func NewSocks5Server(laddr, username, password string) Server {
	s := &socks5erver{
		name:     "socks5Server",
		addr:     laddr,
		username: username,
		password: password,
		connPool: make(chan *ConnEntry, 100),
	}
	return s
}

// StartListen implements Server.
func (s *socks5erver) StartListen() {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Printf("Error created listener in %s: %v\n", s.name, err)
		return
	}
	s.listener = &listener
	log.Printf("%s listener is starting at %s", s.name, listener.Addr())
	go func() {
		for {
			if s.listener == nil {
				log.Printf("listener is empty in %s\n", s.name)
				break
			}
			conn, err := (*s.listener).Accept()
			if err != nil {
				log.Printf("Error accepted connection in %s: %v\n", s.name, err)
				time.Sleep(time.Second) // avoid repeat err
				continue
			}
			go s.HandleConnection(conn)
		}
	}()
}

func (s *socks5erver) HandleConnection(conn net.Conn) {
	err := s.GreetMethod(conn)
	if err != nil {
		log.Printf("Error greet method in %s: %v\n", s.name, err)
		conn.Close()
		return
	}
	err = s.Negotiate(conn)
	if err != nil {
		log.Printf("Error negotiate in %s: %v\n", s.name, err)
		conn.Close()
		return
	}
	connEntry, err := s.Deliver(conn)
	if err != nil {
		log.Printf("Error deliver in %s: %v\n", s.name, err)
		conn.Close()
		return
	}
	err = s.PutWithTimeout(connEntry, time.Duration(10*time.Second))
	if err != nil {
		log.Printf("Error put conn to pool in %s: %v\n", s.name, err)
		conn.Close()
		connEntry = nil
		return
	}
}

func (s *socks5erver) PutWithTimeout(connEntry *ConnEntry, timeout time.Duration) error {
	select {
	case s.connPool <- connEntry:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("put %s connection timeout", connEntry.DstAddr)
	}
}

func (s *socks5erver) GreetMethod(conn net.Conn) error {
	buf := make([]byte, 2)
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		return err
	}
	req := &packet.MethodRequest{
		Ver:      buf[0],
		NMethods: buf[1],
	}
	buf = make([]byte, int(req.NMethods))
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return err
	}
	req.Methods = buf
	if req.Ver != packet.Ver {
		return fmt.Errorf("not a socks5 client")
	}
	var serverMethod byte
	if s.username == "" && s.password == "" {
		serverMethod = packet.MethodNoAuth
	} else {
		serverMethod = packet.MethodUPAuth
	}
	var isClientSupported bool
	for _, method := range req.Methods {
		if method == serverMethod {
			isClientSupported = true
			break
		}
	}
	if !isClientSupported {
		return fmt.Errorf("client methods are not supported")
	}
	resp := &packet.MethodResponse{
		Ver:    req.Ver,
		Method: serverMethod,
	}
	// Write to client
	_, err = conn.Write(resp.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (s *socks5erver) Negotiate(conn net.Conn) error {
	if s.username == "" && s.password == "" {
		return nil
	}
	buf := make([]byte, 2)
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		return err
	}
	// Read socks ver and username length
	req := &packet.NegotiationRequest{
		Ver:  buf[0],
		ULen: buf[1],
	}
	if req.Ver != packet.Ver {
		return fmt.Errorf("not a socks5 client")
	}
	// Read username
	buf = make([]byte, int(req.ULen))
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return err
	}
	req.UName = buf
	// Read password
	buf = make([]byte, 1)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return err
	}
	req.PLen = buf[0]
	buf = make([]byte, int(req.PLen))
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return err
	}
	req.Passwd = buf
	// Auth
	resp := &packet.NegotiationResponse{Ver: req.Ver}
	if string(req.UName) != s.username || string(req.Passwd) != s.password {
		resp.Status = packet.StatuFailure
	} else {
		resp.Status = packet.StatuFailure
	}
	// Write to client
	_, err = conn.Write(resp.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (s *socks5erver) Deliver(conn net.Conn) (*ConnEntry, error) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		if err == io.EOF {
			return nil, fmt.Errorf("current connection don't have data")
		}
		return nil, err
	}
	req := &packet.DeliverRequest{
		Ver:  buf[0],
		Cmd:  buf[1],
		Rsv:  buf[2],
		ATyp: buf[3],
	}
	if req.Ver != packet.Ver {
		return nil, fmt.Errorf("not a socks5 client")
	}
	// Switch address type
	var addrLen int = 0
	if req.ATyp == packet.ATypIPV4 {
		addrLen = 4
	} else if req.ATyp == packet.ATypIPV6 {
		addrLen = 16
	} else if req.ATyp == packet.ATypDomain {
		domainLenBuf := make([]byte, 1)
		_, err := io.ReadFull(conn, domainLenBuf)
		if err != nil {
			return nil, err
		}
		req.DstHostLen = domainLenBuf[0]
		addrLen = int(req.DstHostLen)
	}
	if addrLen == 0 {
		return nil, fmt.Errorf("not supported atyp")
	}
	// Read address
	buf = make([]byte, addrLen)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return nil, err
	}
	req.DstHost = buf
	// Read Port
	buf = make([]byte, 2)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return nil, err
	}
	req.DstPort = [2]byte(buf)
	// New response
	resp := &packet.DeliverResponse{
		Ver:        req.Ver,
		Rep:        packet.UnknowMessage,
		Rsv:        req.Rsv,
		ATyp:       req.ATyp,
		BndHostLen: req.DstHostLen,
		BndHost:    req.DstHost,
		BndPort:    req.DstPort,
	}
	if req.Cmd != packet.CmdConnect && req.Cmd != packet.CmdUDP {
		resp.Rep = packet.RepCommandNotSupported
		_, err = conn.Write(resp.Bytes())
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(packet.ReqMessageMap[resp.Rep])
	}
	// Convert bytes to dst addr
	var host string = ""
	if req.ATyp == packet.ATypDomain {
		host = string(req.DstHost)
	} else {
		ip := net.IP(req.DstHost[:])
		host = string(ip)
	}
	port := int(req.DstPort[0])<<8 + int(req.DstPort[1])
	dstAddr := net.JoinHostPort(host, strconv.Itoa(port))
	connEntry := new(ConnEntry)
	if req.Cmd == packet.CmdConnect {
		connEntry.Network = "tcp"
	} else {
		connEntry.Network = "udp"
	}
	connEntry.ClientConn = conn
	connEntry.DstAddr = dstAddr
	connEntry.ClientResp = resp.Bytes()
	return connEntry, nil
}

// GetConnEntry implements Server.
func (s *socks5erver) GetConnEntry() *ConnEntry {
	return <-s.connPool
}

// StopListen implements Server.
func (s *socks5erver) StopListen() {
	if s.listener == nil {
		log.Printf("listener is empty in %s\n", s.name)
		return
	}
	log.Printf("%s listener is closing at %s\n", s.name, (*s.listener).Addr())
	err := (*s.listener).Close()
	if err != nil {
		log.Printf("Error close listener in %s: %v\n", s.name, err)
		return
	}
	s.listener = nil
	log.Printf("%s listener was closed\n", s.name)
}
