package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
)

type httpClient struct {
	name  string
	raddr string
}

func NewHttpClient(raddr string) Client {
	c := &httpClient{
		name:  "httpClient",
		raddr: raddr,
	}
	return c
}

// Dial implements Client.
func (c *httpClient) Dial(network string, dstAddr string) (net.Conn, error) {
	var proxyConn net.Conn
	var err error
	for i := 0; i < 3; i++ { // Retry up to 3 times
		proxyConn, err = net.Dial(network, c.raddr)
		if err != nil {
			if err == io.EOF {
				continue // If "unexpected EOF", retry the connection
			}
			return nil, err
		}
		break
	}
	if dstAddr == "" {
		return proxyConn, nil
	}
	// Connect to dst addr
	req := &http.Request{
		Method: http.MethodConnect,
		URL:    &url.URL{Host: dstAddr},
	}
	err = req.Write(proxyConn)
	if err != nil {
		return nil, err
	}
	// Read response from proxy server
	resp, err := http.ReadResponse(bufio.NewReader(proxyConn), nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status code from proxy server: %d", resp.StatusCode)
	}
	return proxyConn, nil
}
