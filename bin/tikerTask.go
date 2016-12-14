package main

import (
	"fmt"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"

func getTargerTime(hour, minute, second int, offset int64) int64 {
	utcTime := time.Now().UTC()
	targetTime := time.Date(utcTime.Year(), utcTime.Month(), utcTime.Day(),
		hour, minute, second, 0, utcTime.Location())
	fmt.Println(targetTime)
	return targetTime.Unix() + offset
}

type refreshConfig struct {
	TargetHour      int
	TargetMinute    int
	Targetsecond    int
	Offset          int64
	lastRefreshTime int64
}

var zoneToOffset = map[string]int64{
	"Z0":  0,
	"E1":  -1 * 3600,
	"E2":  -2 * 3600,
	"E3":  -3 * 3600,
	"E4":  -4 * 3600,
	"E5":  -5 * 3600,
	"E6":  -6 * 3600,
	"E7":  -7 * 3600,
	"E8":  -8 * 3600,
	"E9":  -9 * 3600,
	"E10": -10 * 3600,
	"E11": -11 * 3600,
	"E12": 12 * 3600,
	"W1":  1 * 3600,
	"W2":  2 * 3600,
	"W3":  3 * 3600,
	"W4":  4 * 3600,
	"W5":  5 * 3600,
	"W6":  6 * 3600,
	"W7":  7 * 3600,
	"W8":  8 * 3600,
	"W9":  9 * 3600,
	"W10": 10 * 3600,
	"W11": 11 * 3600,
	"W12": 12 * 3600,
}

func timeIsUp(refresh *refreshConfig) bool {

	targetTime := getTargerTime(refresh.TargetHour,
		refresh.TargetMinute,
		refresh.Targetsecond,
		refresh.Offset)

	return refresh.lastRefreshTime < targetTime &&
		time.Now().Unix() >= targetTime
}

func judgeZeroTime() {

	refreshConfigs := []*refreshConfig{}

	refreshConfigs = append(refreshConfigs, &refreshConfig{TargetHour: 0,
		TargetMinute: 1,
		Targetsecond: 0,
		Offset:       zoneToOffset["E8"]})

	for {
		fmt.Println("server Time:", time.Now().Format(timeFormat))

		for _, r := range refreshConfigs {
			if timeIsUp(r) {
				deleteYesterdayUsers()
				r.lastRefreshTime = time.Now().Unix()
			}
			time.Sleep(time.Second * 59)
		}

	}

}
func deleteYesterdayUsers() {
	dailyUserList.Lock()
	dailyUserList.m = make(map[string]string)
	dailyUserList.Unlock()
	nodeUserLogin.Lock()
	nodeOnlineUsers.Lock()
	nodeUserLogin = nodeOnlineUsers
	nodeUserLogin.Unlock()
	nodeOnlineUsers.Unlock()
}
