package server

import (
	"fmt"
	"net"
)

type Client struct {
	name    string
	command string
	conn    net.Conn
}

func (client *Client) clientHandler(list *ClientList) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(client.conn)
	defer list.removeClient(client)
	list.addClient(client)

	for {
		continue
	}

}
