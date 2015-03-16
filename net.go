package main

import (
	"net"

	"code.google.com/p/go.net/proxy"
)

type ProxyConn struct {
	Address  string
	User     string
	Password string
}

func NewConn(p *Proxy, address string) (conn net.Conn, err error) {

	forwardDialer := &net.Dialer{

		DualStack: true,
	}

	if p != nil {

		auth := &proxy.Auth{

			User:     p.User,
			Password: p.Password,
		}

		// setup the socks proxy
		dialer, err := proxy.SOCKS5("tcp", p.Address, auth, forwardDialer)
		if err != nil {

			return nil, err
		}
		return dialer.Dial("tcp", address)
	}

	return forwardDialer.Dial("tcp", address)
}
