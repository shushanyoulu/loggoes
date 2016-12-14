package main

import (
	"strconv"
	"sync"

	"gopkg.in/olivere/elastic.v5"
)

//节点总掉线次数
var nodeOfflineTimes = struct {
	sync.RWMutex
	m map[string]int
}{m: make(map[string]int)}

func nodeOfflineSearch(nodeName string) {
	nodeOfflineTimes.Lock()
	nodeOfflineTimes.m[nodeName]++
	nodeOfflineTimes.Unlock()
}

// nodeOfflineScene 分析用户的掉线情况
func (m *Map) nodeOfflineScene() {
	var countTime = "1m" //统计时间
	nodeOfflineTimes.m = make(map[string]int)
	s := m.store
	for _, pqi := range s.pq {
		if s.onWillEvict != nil {
			nodeOfflineCount(countTime, pqi.key, pqi.item.Value())
		}
	}
	timesOfStatisticNodeOffline()
}

// 插入节点用户掉线次数
func timesOfStatisticNodeOffline() {
	var c esScene
	var r rate
	var nodeOfflineIndexName string
	var offlineIndex *elastic.BulkIndexRequest
	nodeOfflineTimes.Lock()
	for k, v := range nodeOfflineTimes.m {
		ldt.RLock()
		dt := timeChangeToEsFormat(ldt.m[k])
		ldt.RUnlock()
		c = esScene{dt, v}
		r.dt, r.value = dt, v
		nodeRateOffline[k] = r
		nodeOfflineIndexName = k + "-scene"
		offlineIndex = elastic.NewBulkIndexRequest().Index(nodeOfflineIndexName).Type("offline").Doc(c)
		bulkRequest = bulkRequest.Add(offlineIndex)
	}
	nodeOfflineTimes.Unlock()
	if bulkRequest.NumberOfActions() > 10 {
		_, err := bulkRequest.Do(ctx)
		checkerr(err)
	}
}
func nodeOfflineCount(countTime, key string, ier interface{}) {
	var getData eventData
	getData = ier.(eventData)
	a, err := strconv.Atoi(key)
	checkerr(err)
	dataTime := int64(a)
	ldt.RLock()
	before := dealTime(ldt.m[getData.nodeName], countTime)
	ldt.RUnlock()
	if dataTime > before {
		nodeOfflineSearch(getData.nodeName)
	}
}
