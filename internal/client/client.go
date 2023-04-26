package client

import (
	"fmt"
	"net"
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

		fmt.Println(string(buff[:n]))

	}

}
