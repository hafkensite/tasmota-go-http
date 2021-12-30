package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var activeConns = make(map[*websocket.Conn]bool, 0)

var websocketStatus http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	conn.WriteJSON(statuses)

	activeConns[conn] = true
	defer func() {
		delete(activeConns, conn)
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		stateChange := make(map[string]string)
		err = json.Unmarshal(message, &stateChange)
		if err != nil {
			fmt.Printf("%v\n", err)
			break
		}
		for topic, state := range stateChange {
			t := "cmnd/" + topic
			fmt.Printf("Sending command %s = %v", t, state)
			c.Publish(t, 0, false, state)
		}
	}

}
