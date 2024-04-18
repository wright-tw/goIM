package main

import (
	"fmt"
	"strconv"
	"time"
)

type ServerStats struct {
	People int
}

var Stats = ServerStats{}

func subPeople() {
	Stats.People--
}

func addPeople() {
	Stats.People++
}

func logPeople() {
	for {
		fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " 目前有 " + strconv.Itoa(Stats.People) + " 人")
		time.Sleep(time.Second * 3)
	}
}
