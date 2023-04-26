package server

import (
	"sync"
)

type ClientList struct {
	mu   sync.Mutex
	list []*Client
}

func (list *ClientList) addClient(client *Client) {
	list.mu.Lock()
	defer list.mu.Unlock()
	list.list = append(list.list, client)
}

func (list *ClientList) removeClient(client *Client) {
	list.mu.Lock()
	defer client.conn.Close()
	defer list.mu.Unlock()

	for i, cl := range list.list {
		if cl == client {
			list.list = append(list.list[:i], list.list[i+1:]...)
		}
	}

}
