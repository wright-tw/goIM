package main

import (
	"fmt"
	"net/http"
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

type UserConnsStruct struct {
	Lock    sync.Mutex
	ConnMap map[int]*websocket.Conn
}

var UserConns = UserConnsStruct{}
var UserIncr int = 0

func main() {

	go showPeople()
	server := gin.Default()

	// 靜態資料
	server.LoadHTMLGlob("src/views/*.html")
	server.Static("static", "src/views/assets")

	server.GET("/", func(c *gin.Context) {
		// 渲染HTML模板並回應
		c.HTML(http.StatusOK, "index", gin.H{
			"title": "Gin ChatRoom",
		})
	})

	server.GET("/ws", func(c *gin.Context) {

		username := c.Query("name")
		if username == "" {
			username = "No Name"
		}
		fmt.Println(username)

		// 升級連線
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer ws.Close()

		// 歡迎語
		sendText(ws, "Server: hello!")

		// 連線處理
		userId := registerConn(ws, username)
		defer unregisterConn(userId, username)

		// 增加人數
		addPeople()
		defer subPeople()

		for {
			// 處理訊息
			_, textByte, err := ws.ReadMessage()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			text := string(textByte)

			// 處理心跳
			if text == "ping" {
				sendText(ws, "pong")
				continue
			}

			// 處理訊息
			sendMsgToAllPeople(username + ": " + text)
		}
	})

	server.Run(":8081")
}

func sendMsgToAllPeople(msg string) {
	// 發給其他人
	for _, conn := range UserConns.ConnMap {
		sendText(conn, msg)
	}
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
		fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " 目前有 " + strconv.Itoa(Stats.People) + " 人")
		time.Sleep(time.Second * 3)
	}
}

func registerConn(ws *websocket.Conn, username string) int {
	UserConns.Lock.Lock()

	if UserConns.ConnMap == nil {
		UserConns.ConnMap = make(map[int]*websocket.Conn)
	}
	UserIncr++
	UserConns.ConnMap[UserIncr] = ws
	UserConns.Lock.Unlock()

	sendMsgToAllPeople("Server: " + username + " 已進入房間")
	return UserIncr
}

func unregisterConn(userId int, username string) {
	if UserConns.ConnMap != nil {
		delete(UserConns.ConnMap, userId)
	}
	sendMsgToAllPeople("Server: " + username + " 已離開房間")
}
