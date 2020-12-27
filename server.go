package main

import (
	"net"
	"log"
	"sync"
)

type Server struct {
	Clients    []*Client
	Messages []ChatMessage
	Mutex *sync.RWMutex
}

func (s *Server) Join(c *Client) {
	s.Mutex.Lock()
	s.Clients = append(s.Clients, c)
	s.Mutex.Unlock()
}

func (s *Server) Remove(c *Client) {
	s.Mutex.Lock()
	for i, client := range(s.Clients) {
		if client == c {
			client.Connection.Close()
			s.Clients = append(s.Clients[:i], s.Clients[i+1:]...)
			break
		}
	}
	s.Mutex.Unlock()
}

// Broadcasts the command to the server.
func (s *Server) Broadcast(id commandID, data []byte) {
	s.Mutex.RLock()
	for _, c := range(s.Clients) {
		err := SendCommand(c.Connection, id, data)
		if err != nil {
			log.Println("WARNING: Command not send to " + c.Username)
		}
	}
	s.Mutex.RUnlock()
}

// Disconnects client from server and announces it. If the client timed out, timeout must be set to true.
func (s *Server) Disconnect(c *Client, timeout bool) {
	s.Broadcast(CMD_EXIT_NOTIF, []byte(c.Username))
	if !timeout {c.Disconnect <- true}
	s.Remove(c)
}

// Gets a client from the server client list based on the connection given.
func (s *Server) ClientFromConnection(conn net.Conn) *Client {
	s.Mutex.RLock()
	for _, c := range(s.Clients) {
		if c.Connection == conn {
			s.Mutex.RUnlock()
			return c
		}
	}
	s.Mutex.RUnlock()
	return nil
}

// Gets a client from the server client list based on the name given.
func (s *Server) ClientFromName(name string) *Client {
	s.Mutex.RLock()
	for _, c := range(s.Clients) {
		if c.Username == name {
			s.Mutex.RUnlock()
			return c
		}
	}
	s.Mutex.RUnlock()
	return nil
}
