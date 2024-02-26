package service

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"sysproxy/config"
	"sysproxy/service/client"
	"sysproxy/service/server"
)

const (
	errServiceUnavailable    string = "HTTP/1.1 503 Service Unavailable\r\n\r\n"
	msgConnectionEstablished string = "HTTP/1.1 200 Connection Established\r\n\r\n"
)

type httpService struct {
	BaseService
	client   client.Client
	dstProto string
}

func NewHttpService(outbound config.Outbound) BaseService {
	s := new(httpService)
	laddr := net.JoinHostPort(outbound.SrcIP, outbound.SrcPort)
	raddr := net.JoinHostPort(outbound.DstIP, outbound.DstPort)
	var httpServer server.Server
	if outbound.SrcProto == PAC {
		httpServer = server.NewHttpServer(laddr, PACTempFilename)
		s.client = nil
	} else if outbound.DstProto == HTTP {
		httpServer = server.NewHttpServer(laddr, "")
		s.client = client.NewHttpClient(raddr)
	} else if outbound.DstProto == SOCKS5 {
		httpServer = server.NewHttpServer(laddr, "")
		s.client = client.NewSocks5Client(raddr)
	} else {
		return nil
	}
	s.dstProto = outbound.DstProto
	s.BaseService = newBaseService("httpService", httpServer, s)
	return s
}

// GetProxyConn implements IChild.
func (s *httpService) GetProxyConn(connEntry *server.ConnEntry) (net.Conn, error) {
	switch s.dstProto {
	case HTTP:
		return s.Http2Http(connEntry.Network)
	case SOCKS5:
		return s.Http2Socks5(connEntry.Network, connEntry.ClientConn)
	}
	return nil, fmt.Errorf("unknown dstProto: %s", s.dstProto)
}

func (s *httpService) Http2Http(network string) (net.Conn, error) {
	return s.client.Dial(network, "")
}

func (s *httpService) Http2Socks5(network string, clientConn net.Conn) (net.Conn, error) {
	req, err := http.ReadRequest(bufio.NewReader(clientConn))
	if err != nil {
		return nil, fmt.Errorf("read request from clientConn failed: %v", err)
	}
	host := req.URL.Hostname()
	port := req.URL.Port()
	if port == "" {
		port = "443"
		if req.TLS == nil {
			port = "80"
		}
	}
	dstAddr := net.JoinHostPort(host, port)
	proxyConn, err := s.client.Dial(network, dstAddr)
	if err != nil {
		clientConn.Write([]byte(errServiceUnavailable))
		return nil, fmt.Errorf("dial %s timeout: %v", req.Host, err)
	}
	// Establish https tunnel
	if req.Method == http.MethodConnect {
		clientConn.Write([]byte(msgConnectionEstablished))
		return proxyConn, nil
	}
	// Forward http first request
	err = req.Write(proxyConn)
	if err != nil {
		clientConn.Write([]byte(errServiceUnavailable))
		return nil, fmt.Errorf("forward http reqeust failed: %v", err)
	}
	resp, err := http.ReadResponse(bufio.NewReader(proxyConn), nil)
	if err != nil {
		clientConn.Write([]byte(errServiceUnavailable))
		return nil, fmt.Errorf("read response from proxyConn failed: %v", err)
	}
	err = resp.Write(clientConn)
	if err != nil {
		clientConn.Write([]byte(errServiceUnavailable))
		return nil, fmt.Errorf("write response to clienConn failed: %v", err)
	}
	if req.Close {
		// header don't set keep-alive
		return nil, fmt.Errorf("http request end of life")
	}
	return proxyConn, nil
}
