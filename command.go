package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net"
	"strconv"
)

/*
This package is present on both the client and server recieving it. This will effectively
allow communication in a bidirectional way by only reading or doing an action when the
sender has said to do so.

Both parties must have all od the constants, such that they both can understand and send
messages in the same format.

The structure of the information is:


First 6 bytes: Identifier
7th byte: Command ID
8th and onward: Data packaged with it

Data packaged by commands is:

JOIN_REQUEST: Owner in string
CLIENT_EXIT_REQUEST, PING_REQUEST, PING_RESPONSE: None
SERVER_EXIT_REQUEST: Who left the server in string
JOIN_NOTIF: Who joined the server in string
MESSAGE: Message in string

Useage of the commands:

JOIN_REQUEST: Sent by client to server to notify a user join
JOIN_NOTIF: Sent by server to notify clients that someone has joined the server
EXIT_NOTIF: Sent by server to notify clients that someone has left the server
PING: Sent by server to check if clients are still online
PING_RESPONSE: Sent by client to respond to a server ping

MESSAGE: Sent by clients to send a message to the server, which then
the server will broadcast back to every client and when the clients recieve
said request, they will display the message sent
*/

type commandID byte

var verify []byte = []byte("GCCHAT")

const (
	MESSAGE_DELIMITER byte = 0xF9
	EOC               byte = 0xFF
)

const (
	CMD_JOIN_REQUEST commandID = iota
	CMD_JOIN_NOTIF
	CMD_EXIT_NOTIF
	CMD_PING
	CMD_PING_RESPONSE
	CMD_MESSAGE
)

type command struct {
	id   commandID
	data []byte
}

func SendCommand(conn net.Conn, id commandID, data []byte) error {

	_, err := conn.Write(append(append(append(verify, byte(id)), data...), EOC))

	if id == CMD_PING {
		log.Println("Sent a ping")
	} else if id == CMD_PING_RESPONSE {
		log.Println("Sent a ping response")
	} else {
		log.Println("Sent command with id: " + id.ToString() + " with data: " + string(data))
	}

	return err
}

func RetrieveCommands(conn net.Conn) ([]command, bool) {

	commands := []command{}

	reader := bufio.NewReader(conn)

	log.Println("Reading...")

	for left := true; left; left = reader.Buffered() != 0 {
		comm, err := reader.ReadBytes(EOC)

		if err == io.EOF {
			return []command{}, false
		}

		if len(comm) != len(verify) {
			log.Println("Read " + strconv.Itoa(len(comm)) + " bytes")
		}

		comm = comm[:len(comm)-1]
		id := commandID(comm[len(verify)])

		if id == CMD_PING {
			log.Println("Recieved a ping")
		} else if id == CMD_PING_RESPONSE {
			log.Println("Recieved a ping response")
		} else {
			log.Println("Command was read, with id " + id.ToString() + " and data " + string(comm[len(verify)+1:]))
		}

		if !bytes.Equal(comm[:len(verify)], verify) {
			panic("Error at retrieving command: Command is not valid")
		}

		commands = append(commands, command{id, comm[len(verify)+1:]})
	}

	return commands, true
}

func (c commandID) ToString() string {
	switch c {
	case CMD_JOIN_REQUEST:
		return "CMD_JOIN_REQUEST"
	case CMD_JOIN_NOTIF:
		return "CMD_JOIN_NOTIF"
	case CMD_EXIT_NOTIF:
		return "CMD_EXIT_NOTIF"
	case CMD_PING:
		return "CMD_PING"
	case CMD_PING_RESPONSE:
		return "CMD_PING_RESPONSE"
	case CMD_MESSAGE:
		return "CMD_MESSAGE"
	}
	return "INVALID ID"
}
