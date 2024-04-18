package main

import (
	"goIM/pkg/encode"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type UserConnsStruct struct {
	Lock    sync.Mutex
	ConnMap map[int]*websocket.Conn
}

var UserConns = UserConnsStruct{}
var UserIncr int = 0

var ServerName = "Server"

func registerConn(ws *websocket.Conn, username string) int {
	UserConns.Lock.Lock()

	if UserConns.ConnMap == nil {
		UserConns.ConnMap = make(map[int]*websocket.Conn)
	}
	UserIncr++
	UserConns.ConnMap[UserIncr] = ws
	UserConns.Lock.Unlock()

	sendMsgToAllPeople(ACTION_SYSTEM_MSG, ServerName, username+" 已進入房間")
	return UserIncr
}

func unregisterConn(userId int, username string) {
	if UserConns.ConnMap != nil {
		delete(UserConns.ConnMap, userId)
	}
	sendMsgToAllPeople(ACTION_SYSTEM_MSG, ServerName, username+" 已離開房間")
}

func getMsgMap() map[string]interface{} {
	return map[string]interface{}{
		"action":   nil,
		"username": nil,
		"msg":      nil,
	}
}
func sendMsgToAllPeople(action int, username string, msg string) {
	for _, conn := range UserConns.ConnMap {
		sendMsgToPeople(conn, action, username, msg)
	}
}

func sendMsgToPeople(ws *websocket.Conn, action int, username string, msg string) {
	mData := getMsgMap()
	mData["action"] = action
	mData["username"] = username
	mData["msg"] = msg
	jsonMsg := encode.JSONEncode(mData)
	sendText(ws, jsonMsg)
}

func sendText(ws *websocket.Conn, msg string) {
	ws.WriteMessage(websocket.TextMessage, []byte(msg))
}

func sendOnlineCount(ws *websocket.Conn) {
	sendMsgToPeople(ws, ACTION_ONLINE_PEOPLE, ServerName, strconv.Itoa(Stats.People))
}

func sendOnlineCountToAllPeople() {
	for _, conn := range UserConns.ConnMap {
		sendOnlineCount(conn)
	}
}
