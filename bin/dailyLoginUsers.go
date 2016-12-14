package main

import (
	"sync"
	"time"

	"gopkg.in/olivere/elastic.v5"
)

type dailyOnline struct {
	LdateTime string `json:"时间"`
	NodeName  string `json:"节点名称"`
	Sum       int    `json:"已登录人数"`
}

//每个节点当前每日已登录用户数 [nodeName][]uid
var nodeUserLogin = struct {
	sync.RWMutex
	m map[string][]string
}{m: make(map[string][]string)}

var dailyUserList = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

var getDailyuserTiker = time.NewTicker(1 * time.Minute)

func (nodeLog broadLogData) dailyUser() {
	node, uid := nodeLog.nodeName, analysisUid(nodeLog.data)
	dailyUserList.Lock()
	if _, ok := dailyUserList.m[uid]; ok == true {
		// fmt.Println(ok)
	} else {
		dailyUserList.m[uid] = node
	}
	dailyUserList.Unlock()
	// time.Sleep(10e8)
}

func tikerGetNodeUsers() {
	var c dailyOnline
	var nodeDailySignUp string
	var dailyOnlineIndex *elastic.BulkIndexRequest
	dailyUserList.RLock()
	nodeUserLogin.Lock()
	nodeUserLogin.m = make(map[string][]string)
	for k, v := range dailyUserList.m {
		nodeUserLogin.m[v] = append(nodeUserLogin.m[v], k)
	}
	dailyUserList.RUnlock()
	for k, v := range nodeUserLogin.m {
		ldt.RLock()

		dt := timeChangeToEsFormat(ldt.m[k])
		ldt.RUnlock()
		c = dailyOnline{dt, k, len(v)}
		nodeDailySignUp = k + "-dailysignup"
		dailyOnlineIndex = elastic.NewBulkIndexRequest().Index(nodeDailySignUp).Type("INFO").Doc(c)
		bulkRequest = bulkRequest.Add(dailyOnlineIndex)
	}
	_, err := bulkRequest.Do(ctx)
	checkerr(err)
	nodeUserLogin.Unlock()
}
