package server

import (
	"fmt"
	"net"
)

type Client struct {
	name     string
	conn     net.Conn
	isActual bool
}

func (client *Client) newClient(list *ClientList) {
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
		if !client.isActual {
			break
		}
	}

}
