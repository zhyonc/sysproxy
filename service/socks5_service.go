package service

import (
	"fmt"
	"net"
	"sysproxy/config"
	"sysproxy/service/client"
	"sysproxy/service/server"
)

type socks5Service struct {
	BaseService
	client client.Client
}

func NewSocks5Service(outbound config.Outbound) BaseService {
	s := new(socks5Service)
	laddr := net.JoinHostPort(outbound.SrcIP, outbound.SrcPort)
	raddr := net.JoinHostPort(outbound.DstIP, outbound.DstPort)
	if outbound.DstProto == HTTP {
		s.client = client.NewHttpClient(raddr)
	} else if outbound.DstProto == SOCKS5 {
		s.client = client.NewSocks5Client(raddr)
	} else {
		return nil
	}
	s.BaseService = newBaseService("httpService", server.NewSocks5Server(laddr, "", ""), s)
	return s
}

// GetProxyConn implements IChild.
func (s *socks5Service) GetProxyConn(connEntry *server.ConnEntry) (net.Conn, error) {
	proxyConn, err := s.client.Dial(connEntry.Network, connEntry.DstAddr)
	resp := connEntry.ClientResp
	if err != nil {
		resp[1] = 0x03 // network unreachable
		_, _ = connEntry.ClientConn.Write(resp)
		return nil, fmt.Errorf("dial %s failed: %v", connEntry.DstAddr, err)
	}
	resp[1] = 0x00 // connection established successfully
	_, err = connEntry.ClientConn.Write(resp)
	if err != nil {
		proxyConn.Close()
		return nil, fmt.Errorf("clientConn write socks5 response failed")
	}
	return proxyConn, nil
}
