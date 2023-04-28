package server

import (
	"fmt"
	"github.com/pterm/pterm"
	"net"
	"strconv"
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

					_, err := client.conn.Write([]byte(result))
					if err != nil {
						return
					}

					buff := make([]byte, 1024)

					n, err := client.conn.Read(buff)

					if err != nil {
						fmt.Println(err)
						continue
					}
					fileSize, _ := strconv.ParseInt(strings.Trim(string(buff[:n]), ":"), 10, 64)
					fmt.Println("Size ", fileSize)

					n, err = client.conn.Read(buff)

					if err != nil {
						fmt.Println(err)
						continue
					}
					fileName := strings.Trim(string(buff[:n]), ":")
					fmt.Println("Filename  ", fileName)

					//bufferFileName := make([]byte, 1024)
					//bufferFileSize := make([]byte, 1024)
					////
					//fmt.Println(string(bufferFileSize))
					//fmt.Println(string(bufferFileName))

					//client.conn.Read(bufferFileSize)
					//fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
					//
					//client.conn.Read(bufferFileName)
					//fileName := strings.Trim(string(bufferFileName), ":")
					//
					//fmt.Println(fileSize)
					//fmt.Println(fileName)

					//newFile, _ := os.Create(fileName + "_2")
					//
					//var receivedBytes int64
					//
					//for {
					//	if (fileSize - receivedBytes) < 1024 {
					//		io.CopyN(newFile, client.conn, (fileSize - receivedBytes))
					//		client.conn.Read(make([]byte, (receivedBytes+1024)-fileSize))
					//		break
					//	}
					//	io.CopyN(newFile, client.conn, 1024)
					//	receivedBytes += 1024
					//}
					//fmt.Println("Received file completely!")
					//
					//newFile.Close()

				}
			}
		}
	} else {
		pterm.DefaultBasicText.WithStyle(pterm.FgDarkGray.ToStyle()).Println("Wrong command. Please input 'ls' or 'list' to show all connections")
	}
}
