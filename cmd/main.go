package main

import (
	"fmt"
	"strconv"
	"sync"
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

var Stats = ServerStats{}

type Room struct {
	UserIds []int
	Msgs    []string
}

var room1 = Room{}

var NowUserIncr int64 = 0

type UserConnsStruct struct {
	Lock    sync.Mutex
	ConnMap map[int]*websocket.Conn
}

var UserConns = UserConnsStruct{}
var UserIncr = 0

func main() {

	go showPeople()

	router := gin.Default()
	router.GET("/ws", func(c *gin.Context) {

		// 升級連線
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer ws.Close()

		// 連線處理
		userId := registerConn(ws)
		defer unregisterConn(userId)

		// 增加人數
		addPeople()
		defer subPeople()

		// 房間人數增加
		room1.UserIds = append(room1.UserIds, userId)

		sendText(ws, "hello!")
		for {
			_, textByte, err := ws.ReadMessage()
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			text := string(textByte)
			fmt.Println(text)
		}
	})

	router.Run("127.0.0.1:8081")
}

func sendText(conn *websocket.Conn, msg string) {
	conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func subPeople() {
	Stats.People--
}

func addPeople() {
	Stats.People++
}

func showPeople() {
	for {
		fmt.Println("目前有 " + strconv.Itoa(Stats.People) + " 人")
		fmt.Println(UserConns.ConnMap)
		time.Sleep(time.Second)
	}
}

func registerConn(c *websocket.Conn) int {
	UserConns.Lock.Lock()

	if UserConns.ConnMap == nil {
		UserConns.ConnMap = make(map[int]*websocket.Conn)
	}
	UserIncr++
	UserConns.ConnMap[UserIncr] = c
	UserConns.Lock.Unlock()
	return UserIncr
}

func unregisterConn(userId int) {
	if UserConns.ConnMap != nil {
		delete(UserConns.ConnMap, userId)
	}
}
