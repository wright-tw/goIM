package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	go logPeople()
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
		sendMsgToPeople(ws, ACTION_SYSTEM_MSG, ServerName, "hello")

		// 連線處理
		userId := registerConn(ws, username)
		defer unregisterConn(userId, username)

		// 增加人數
		addPeople()

		// 推送當前人數
		sendOnlineCountToAllPeople()

		defer sendOnlineCountToAllPeople()
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
			sendMsgToAllPeople(ACTION_MSG, username, text)
		}
	})

	server.Run("127.0.0.1:8081")
}
