package main

import (
	"fmt"
	"time"
)

var nodeSceneTimer = time.NewTicker(time.Minute)

func tikers() {
	fmt.Println("启动定时器")
	for {
		select {
		case <-nodeSceneTimer.C: //每分钟统计在线用户数，掉线用户数，掉线比率
			nodeScene()
		case <-getDailyuserTiker.C: //每10分钟统计，本日已登录用户数。
			tikerGetNodeUsers()
		}
	}
}

// 统计用户在线用户，掉线用户，掉线比率
func nodeScene() {
	statisticNodeOnline()
	offLineList.nodeOfflineScene()
	reckonRate()
}
