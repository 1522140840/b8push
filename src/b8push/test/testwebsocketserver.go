package main

import (
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
)

var upgrader = websocket.Upgrader{}
func main() {
	http.HandleFunc("/echo", echoHandler)

	http.ListenAndServe(":8081", nil)
}

func  echoHandler(w http.ResponseWriter, r *http.Request)  {
	c, _ := upgrader.Upgrade(w, r, nil)


	for {

		_, message, err := c.ReadMessage()
		if err != nil {

			break
		}
		fmt.Println(message)
	}


}
