package main

import (
	"log"
	"net"
	"sync"
	"time"
)

const PORT string = "9000"

func main() {
	ln, err := net.Listen("tcp", ":" + PORT)

	s := &Server{Mutex: &sync.RWMutex{}}

	if err != nil {
		log.Fatalf("Error listening to connection: " + err.Error())
	}

	log.Println("Server has started at port 9000 of this server")

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Println("Error recieving connection: " + err.Error())
		}

		c := &Client{Connection: conn, Ping: make(chan bool, 1), Disconnect: make(chan bool, 1)}

		go Listen(c, s)
		go Timeout(c, s)

	}
}

func Timeout(client *Client, server *Server) {
	for {
		select {
		case <-client.Ping:
			time.Sleep(time.Second * 10)
			SendCommand(client.Connection, CMD_PING, []byte{})
		case <-time.After(time.Second * 10):
			log.Println("Timed out client with name " + client.Username)
			server.Disconnect(client, true)
			return
		case <-client.Disconnect:
			log.Println("Client disconnected with name " + client.Username)
			return
		}
	}
}

/*
The listen function will be ran for each connection made to the server, so each user
is being listened by it's own routine.
*/
func Listen(client *Client, server *Server) {

	defer server.Disconnect(client, false)

	for {
		commands, ok := RetrieveCommands(client.Connection)

		if !ok {
			break
		}

		for _, c := range commands {
			switch c.id {
			case CMD_JOIN_REQUEST:
				client.Username = string(c.data)
				server.Join(client)
				server.Mutex.RLock()

				for _, cl := range server.Clients {
					if client != cl {
						SendCommand(client.Connection, CMD_JOIN_NOTIF, []byte(cl.Username))
					}
				}

				server.Mutex.RUnlock()

				server.Broadcast(CMD_JOIN_NOTIF, c.data)
			case CMD_PING_RESPONSE:
				go func() { client.Ping <- true }()
			case CMD_MESSAGE:
				server.Broadcast(CMD_MESSAGE, append(append([]byte(client.Username), MESSAGE_DELIMITER), c.data...))
			}
		}

	}
}
