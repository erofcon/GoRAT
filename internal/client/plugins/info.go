package plugins

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func Info() (string, error) {

	response := ""

	client := http.Client{Timeout: 6 * time.Second}

	res, err := client.Get("https://geolocation-db.com/json")
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	if err != nil {
		fmt.Println(err)
		return response, err
	}

	for {
		buff := make([]byte, 1024)

		n, err := res.Body.Read(buff)
		response += string(buff[:n])

		if n == 0 || err != nil {
			break
		}

	}

	return response, nil
}
