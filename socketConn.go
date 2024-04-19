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
	userIncrLock sync.Mutex
	ConnMap      sync.Map // user_incr => conn
	ConnLockMap  sync.Map // ws => mutx
}

var UserConns = UserConnsStruct{sync.Mutex{}, sync.Map{}, sync.Map{}}
var UserIncr int = 0

var ServerName = "Server"

// 歷史訊息
var HistoryMsgs = []string{}

func registerConn(ws *websocket.Conn, username string) int {

	UserConns.userIncrLock.Lock()
	UserIncr++
	UserConns.ConnMap.Store(UserIncr, ws)
	UserConns.userIncrLock.Unlock()
	UserConns.ConnLockMap.Store(ws, &sync.Mutex{})
	sendMsgToAllPeople(ACTION_SYSTEM_MSG, ServerName, username+" 已進入房間")
	return UserIncr
}

func unregisterConn(ws *websocket.Conn, userId int, username string) {
	UserConns.ConnMap.Delete(userId)
	UserConns.ConnLockMap.Delete(ws)
	sendMsgToAllPeople(ACTION_SYSTEM_MSG, ServerName, username+" 已離開房間")
}

func getMsgMap() map[string]interface{} {
	return map[string]interface{}{
		"action":   nil,
		"username": nil,
		"msg":      nil,
		// "time":     time.Now().Format("2006-01-02 15:04:05"),
		"time": time.Now().Format("15:04"),
	}
}
func sendMsgToAllPeople(action int, username string, msg string) {
	UserConns.ConnMap.Range(func(_, conn interface{}) bool {
		connTyped := conn.(*websocket.Conn)
		sendMsgToPeople(connTyped, action, username, msg)
		return true // 继续遍历
	})
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
	// 同一個 socket client 同時只能寫一次
	lockInterface, _ := UserConns.ConnLockMap.Load(ws)
	if lockInterface == nil {
		return
	}
	lock := lockInterface.(*sync.Mutex)
	lock.Lock()
	defer lock.Unlock()

	ws.WriteMessage(websocket.TextMessage, []byte(msg))
}

func sendOnlineCount(ws *websocket.Conn) {
	sendMsgToPeople(ws, ACTION_ONLINE_PEOPLE, ServerName, strconv.Itoa(Stats.People))
}

func sendOnlineCountToAllPeople() {
	UserConns.ConnMap.Range(func(_, conn interface{}) bool {
		connTyped := conn.(*websocket.Conn)
		sendOnlineCount(connTyped)
		return true // 继续遍历
	})
}

func sendHistoryMsg(ws *websocket.Conn) {
	for _, oldJsonMsg := range HistoryMsgs {
		sendText(ws, oldJsonMsg)
	}
}
