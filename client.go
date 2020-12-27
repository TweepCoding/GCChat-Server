package main

import "net"

type Client struct {
	Username string
	Connection net.Conn
	Ping chan bool
	Disconnect chan bool
}
