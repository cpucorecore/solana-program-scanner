package http_client

import (
	"net"
	"net/http"
	"time"
)

var dialer = net.Dialer{
	Timeout:   30 * time.Second,
	KeepAlive: 20 * time.Second,
}

var getBlockTransport = &http.Transport{
	DialContext:           dialer.DialContext,
	MaxIdleConns:          30,
	MaxIdleConnsPerHost:   30,
	MaxConnsPerHost:       60,
	IdleConnTimeout:       90 * time.Second,
	ReadBufferSize:        1024 * 1024 * 1,
	ForceAttemptHTTP2:     true,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

var commonTransport = &http.Transport{
	DialContext:           dialer.DialContext,
	MaxIdleConns:          10,
	MaxIdleConnsPerHost:   10,
	MaxConnsPerHost:       10,
	IdleConnTimeout:       90 * time.Second,
	ReadBufferSize:        1024 * 16,
	ForceAttemptHTTP2:     true,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}
