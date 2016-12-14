package main

import (
	"strings"
	"time"

	"sync"

	"gopkg.in/olivere/elastic.v5"
)

// es中即时在线json结构
type esScene struct {
	LdateTime string `json:"即时时间"`
	Value     int    `json:"值"`
}

var options = &Options{
	InitialCapacity: 1024,
	OnWillExpire: func(key string, item *Item) {
	},
	OnWillEvict: func(key string, item *Item) {
	},
}

// 每个节点当前在线用户缓存[nodeName][]uid
var nodeOnlineUsers = struct {
	sync.RWMutex
	m map[string][]string
}{m: make(map[string][]string)}

var userOnline = New(options)

// online 实时统计各个节点在线用户数据
func (nodeLog broadLogData) onlines() {
	nodeName := nodeLog.nodeName
	l := nodeLog.data
	uid := analysisUid(l)
	// fmt.Println(nodeName, uid)
	switch {
	case strings.Contains(l, "LOGOUT") == false:
		userOnline.SetNX(uid, NewItemWithTTL(nodeName, 12*time.Hour))
	case strings.Contains(l, "LOGOUT"):
		userOnline.Delete(uid)

	}
}

// 统计在线用户数
func searchMap(nodeOnlineUsers map[string][]string) {
	s := userOnline.store
	for _, pqi := range s.pq {
		if s.onWillEvict != nil {
			ier := pqi.item.Value().(string)
			nodeOnlineUsers[ier] = append(nodeOnlineUsers[ier], pqi.key)
		}
	}
}

//插入在线用户数到es
func statisticNodeOnline() {
	var c esScene
	var r rate
	var nodeOnlineIndex string
	var onlineIndex *elastic.BulkIndexRequest
	nodeOnlineUsers.Lock()
	nodeOnlineUsers.m = make(map[string][]string)
	searchMap(nodeOnlineUsers.m)
	for k, v := range nodeOnlineUsers.m {
		ldt.RLock()
		dt := timeChangeToEsFormat(ldt.m[k])
		ldt.RUnlock()
		c = esScene{dt, len(v)}
		r.dt, r.value = dt, len(v)
		nodeRateOnline[k] = r
		nodeOnlineIndex = k + "-scene"
		onlineIndex = elastic.NewBulkIndexRequest().Index(nodeOnlineIndex).Type("online").Doc(c)
		bulkRequest = bulkRequest.Add(onlineIndex)
	}
	_, err := bulkRequest.Do(ctx)
	checkerr(err)
	nodeOnlineUsers.Unlock()
}
