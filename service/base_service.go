package service

import (
	"io"
	"log"
	"net"
	"sync"
	"sysproxy/service/server"
)

type baseService struct {
	name   string
	server server.Server
	conns  sync.Map
	IChild
}

type IChild interface {
	GetProxyConn(connEntry *server.ConnEntry) (net.Conn, error)
}

func newBaseService(name string, server server.Server, child IChild) BaseService {
	s := &baseService{
		name:   name,
		server: server,
		IChild: child,
	}
	return s
}

// StartService implements BaseService.
func (s *baseService) StartService() {
	s.server.StartListen()
	s.Connect()
}

func (s *baseService) Connect() {
	for {
		connEntry := s.server.GetConnEntry()
		proxyConn, err := s.IChild.GetProxyConn(connEntry)
		if err != nil {
			log.Printf("Error get proxyConn: %v", err)
			connEntry.ClientConn.Close()
			continue
		}
		s.conns.Store(connEntry.ClientConn, nil)
		s.conns.Store(connEntry.ClientConn, nil)
		go s.Forward(proxyConn, connEntry.ClientConn)
	}
}

// Forward implements BaseService.
func (s *baseService) Forward(proxyConn, clientConn net.Conn) {
	defer s.disconnect(proxyConn)
	defer s.disconnect(clientConn)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := io.Copy(proxyConn, clientConn)
		if err != nil {
			log.Printf("Error copying from clientConn to proxyConn: %v", err)
		}
	}()
	go func() {
		defer wg.Done()
		_, err := io.Copy(clientConn, proxyConn)
		if err != nil {
			log.Printf("Error copying from proxyConn to clientConn: %v", err)
		}
	}()
	wg.Wait()
}

func (s *baseService) disconnect(conn net.Conn) {
	_, ok := s.conns.Load(conn)
	if !ok {
		return
	}
	s.conns.Delete(conn)
	err := conn.Close()
	if err != nil {
		log.Printf("Error disconnect: %v", err)
	}
}

func (s *baseService) disconnectAll() {
	s.conns.Range(func(key, _ any) bool {
		s.conns.Delete(key)
		conn, ok := key.(net.Conn)
		if !ok {
			log.Printf("Error disconnect all because dessert net.Conn failed")
		} else {
			err := conn.Close()
			if err != nil {
				log.Printf("Error disconnect: %v", err)
			}
		}
		return true
	})
}

// StopService implements BaseService.
func (s *baseService) StopService() {
	s.disconnectAll()
	s.server.StopListen()
}
