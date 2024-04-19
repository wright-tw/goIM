package main

import (
	"goIM/pkg/encode"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type UserConnsStruct struct {
	Lock        sync.Mutex
	ConnMap     map[int]*websocket.Conn
	ConnLockMap map[*websocket.Conn]*sync.Mutex
}

var UserConns = UserConnsStruct{}
var UserIncr int = 0

var ServerName = "Server"

// 歷史訊息
var HistoryMsgs = []string{}

func registerConn(ws *websocket.Conn, username string) int {
	UserConns.Lock.Lock()
	defer UserConns.Lock.Unlock()

	if UserConns.ConnMap == nil {
		UserConns.ConnMap = make(map[int]*websocket.Conn)
	}
	if UserConns.ConnLockMap == nil {
		UserConns.ConnLockMap = make(map[*websocket.Conn]*sync.Mutex)
	}
	UserIncr++
	UserConns.ConnMap[UserIncr] = ws
	UserConns.ConnLockMap[ws] = &sync.Mutex{}

	sendMsgToAllPeople(ACTION_SYSTEM_MSG, ServerName, username+" 已進入房間")
	return UserIncr
}

func unregisterConn(ws *websocket.Conn, userId int, username string) {
	if UserConns.ConnMap != nil {
		delete(UserConns.ConnMap, userId)
	}
	if UserConns.ConnLockMap != nil {
		delete(UserConns.ConnLockMap, ws)
	}
	sendMsgToAllPeople(ACTION_SYSTEM_MSG, ServerName, username+" 已離開房間")
}

func getMsgMap() map[string]interface{} {
	return map[string]interface{}{
		"action":   nil,
		"username": nil,
		"msg":      nil,
		"time":     time.Now().Format("2006-01-02 15:04:05"),
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

	// 普通訊息存起來
	if action == ACTION_MSG {
		if len(HistoryMsgs) > 100 {
			HistoryMsgs = HistoryMsgs[1:] // 刪除切片中的第一個元素（最舊的數據）
		}
		HistoryMsgs = append(HistoryMsgs, jsonMsg)
	}

	sendText(ws, jsonMsg)
}

func sendText(ws *websocket.Conn, msg string) {

	// 同一個 socket client 只能寫一次
	UserConns.ConnLockMap[ws].Lock()
	defer UserConns.ConnLockMap[ws].Unlock()
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

func sendHistoryMsg(ws *websocket.Conn) {
	for _, oldJsonMsg := range HistoryMsgs {
		sendText(ws, oldJsonMsg)
	}
}
