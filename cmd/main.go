package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ServerStats struct {
	People int
}

var stats = ServerStats{}

func main() {

	go showPeople()

	router := gin.Default()
	router.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		// 增加人數
		addPeople()

		defer conn.Close()
		defer subPeople()

		sendText(conn, "hello!")
		for {
		}
	})

	router.Run("127.0.0.1:8081")
}

func sendText(conn *websocket.Conn, msg string) {
	conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func subPeople() {
	stats.People--
}

func addPeople() {
	stats.People++
}

func showPeople() {
	for {
		fmt.Println("目前有 " + strconv.Itoa(stats.People) + " 人")
		time.Sleep(time.Second)
	}
}
