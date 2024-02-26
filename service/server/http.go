package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type httpServer struct {
	name     string
	addr     string
	filename string
	listener *net.Listener
	connPool chan *ConnEntry
}

func NewHttpServer(laddr string, filename string) Server {
	s := &httpServer{
		name:     "httpServer",
		addr:     laddr,
		filename: filename,
		connPool: make(chan *ConnEntry, 100),
	}
	return s
}

// StartListen implements Server.
func (s *httpServer) StartListen() {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Printf("Error created listener in %s: %v\n", s.name, err)
		return
	}
	s.listener = &listener
	log.Printf("%s listener is starting at %s\n", s.name, listener.Addr())
	if s.filename != "" {
		go http.Serve(listener, http.HandlerFunc(s.handlePACRequest))
	} else {
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
				go s.handleConnection(conn)
			}
		}()
	}
}

func (s *httpServer) handlePACRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("request pac file from %s", r.RemoteAddr)
	_, err := os.Stat(s.filename)
	if err != nil {
		http.Error(w, "404 page not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/x-ns-proxy-autoconfig")
	http.ServeFile(w, r, s.filename)
}

func (s *httpServer) handleConnection(conn net.Conn) {
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

// Deliver implements Server.
func (s *httpServer) Deliver(conn net.Conn) (*ConnEntry, error) {
	connEntry := new(ConnEntry)
	connEntry.Network = "tcp"
	connEntry.ClientConn = conn
	connEntry.DstAddr = ""
	return connEntry, nil
}

func (s *httpServer) PutWithTimeout(connEntry *ConnEntry, timeout time.Duration) error {
	select {
	case s.connPool <- connEntry:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("put %s connection timeout", connEntry.DstAddr)
	}
}

// GetConnEntry implements Server.
func (s *httpServer) GetConnEntry() *ConnEntry {
	return <-s.connPool
}

// StopListen implements Server.
func (s *httpServer) StopListen() {
	if s.listener == nil {
		log.Printf("listener is empty in %s\n" + s.name)
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

// Deprecated
// func (s *httpServer) getRawRequest(r *http.Request) string {
// 	// Create the request line with the absolute URL
// 	requestLine := fmt.Sprintf("%s %s HTTP/1.1\r\nHost: %s\r\n", r.Method, r.URL.String(), r.URL.Hostname())
// 	// Build the headers of the request
// 	var builder strings.Builder
// 	for name, values := range r.Header {
// 		for _, value := range values {
// 			builder.WriteString(name)
// 			builder.WriteString(": ")
// 			builder.WriteString(value)
// 			builder.WriteString("\r\n")
// 		}
// 	}
// 	headers := builder.String()
// 	// Read the body into a string
// 	bodyBytes, _ := io.ReadAll(r.Body)
// 	// Combine the request line, headers, and body into a raw HTTP request
// 	rawRequest := requestLine + headers + "\r\n" + string(bodyBytes)
// 	return rawRequest
// }
