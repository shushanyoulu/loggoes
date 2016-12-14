package main

import "gopkg.in/olivere/elastic.v5"

type rate struct {
	dt    string
	value int
}

var nodeRateOnline = make(map[string]rate)
var nodeRateOffline = make(map[string]rate)

//计算节点掉线用户/在线用户
func reckonRate() {
	var c esScene
	var nodeRateIndex string
	var rateIndex *elastic.BulkIndexRequest
	for k, v := range nodeRateOnline {
		a := nodeRateOffline[k].value
		b := v.value
		r := float32(a) / float32(b) * 1000
		c.LdateTime = v.dt
		c.Value = int(r)
		nodeRateIndex = k + "-scene"
		rateIndex = elastic.NewBulkIndexRequest().Index(nodeRateIndex).Type("rate").Doc(c)
		bulkRequest = bulkRequest.Add(rateIndex)
	}
	_, err := bulkRequest.Do(ctx)
	checkerr(err)
}
