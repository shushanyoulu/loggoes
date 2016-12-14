package main

import (
	"strconv"
	"sync"
	"time"

	"gopkg.in/olivere/elastic.v5"
)

type eventData struct {
	dt       string
	uid      string
	event    string
	gid      []string
	sign     int
	nodeName string
}
type offLineUserInfo struct {
	LdateTime string `json:"时间"`
	UID       string `json:"UID"`
	Info      string `json:"信息"`
	NodeName  string `json:"节点名称"`
}
type users struct {
	nodeName string
	m        map[string]int
}

var offLineList = New(eventOffLineData) //掉线数据缓存

var eventOffLineData = &Options{
	InitialCapacity: 1024,
	OnWillExpire: func(key string, item *Item) {
	},
	OnWillEvict: func(key string, item *Item) {
	},
}
var userLastDataList = struct {
	sync.RWMutex
	m map[string]streamData
}{m: make(map[string]streamData)}

var offlineDataLastTime = struct { //记录掉线缓存数据中的最后一次时间
	sync.RWMutex
	m int64
}{m: 0}

var offlineTimeSet = float64(offlineInterval())

//将分析出来的用户掉线数据插入缓存中
func (e eventData) InsertList() {
	if e.sign == 1 {
		t, err := time.Parse("2006/01/02 15:04:05", e.dt)
		checkerr(err)
		offlineDataLastTime.Lock()
		offlineDataLastTime.m = t.Unix()
		w := strconv.Itoa(int(offlineDataLastTime.m))
		offLineList.Set(w, NewItemWithTTL(e, 15*time.Minute)) //掉线事件缓存20分钟
		e.timesOfStatisticUserOffline()
		offlineDataLastTime.Unlock()
	}
}

//periodOfUsersOffline 插入用户掉线数据至缓存和es中
func (b broadLogData) periodOfUsersOffline() {
	var e eventData
	uid := analysisUid(b.data)
	userLastDataList.RLock()
	historyData := userLastDataList.m[uid]
	userLastDataList.RUnlock()
	newData := b.makeLogFormat()
	newData.dealStreamData(b.nodeName)
	e = analysisStreamStu(uid, b.nodeName, historyData, newData)
	e.InsertList()
}

//将原始日志进行格式化
func (b broadLogData) makeLogFormat() logFormat {
	var f logFormat
	d, t, info := extract(b.data)
	if len(t) > 8 {
		f.dateTime = d + " " + t[:8]
	} else {
		goesl.Println(l)
		return f
	}
	f.uid = analysisUid(b.data)
	f.stu = analysisStu(info)
	f.gid = analysisGid(b.data)
	f.nodeName = b.nodeName
	return f
}
func dealTime(tdt, i string) int64 {
	i = "-" + i
	n, err := time.Parse("2006/01/02 15:04:05", tdt)
	checkerr(err)
	m, err := time.ParseDuration(i)
	checkerr(err)
	n = n.Add(m)
	a := n.Unix()
	return a
}

//插入每个用户掉线
func (e eventData) timesOfStatisticUserOffline() {
	var c offLineUserInfo
	var userOfflineIndex string
	userOfflineIndex = e.nodeName + "-useroffline"
	c.LdateTime = timeChangeToEsFormat(e.dt)
	c.UID = e.uid
	c.NodeName = e.nodeName
	c.Info = e.event
	request := elastic.NewBulkIndexRequest().Index(userOfflineIndex).Type("INFO").Doc(c)
	bulkRequest = bulkRequest.Add(request)
	if bulkRequest.NumberOfActions() > 10 {
		_, err := bulkRequest.Do(ctx)
		checkerr(err)
	}
}

//将流数据分析处理
func (s logFormat) dealStreamData(nodeName string) {
	uid := s.uid
	var st streamData
	userLastDataList.Lock()
	if _, ok := userLastDataList.m[uid]; ok == true { // 如果缓存数据中数据存在，则分析数据产生事件
		st.dateTime = s.dateTime
		st.stu = stuNameModify(s.stu)
		st.gid = s.gid
		userLastDataList.m[uid] = st

	} else { //如果缓存数据流中数据不存在，则添加数据
		st.dateTime = s.dateTime
		st.stu = stuNameModify(s.stu)
		st.gid = s.gid
		userLastDataList.m[uid] = st

	}
	userLastDataList.Unlock()
}

//分析数据数据数据，返回状态事件
func analysisStreamStu(uid, nodeName string, last streamData, this logFormat) eventData {
	var u eventData
	u.nodeName = nodeName
	var k float64
	switch {
	case last.stu == "LOGIN" && this.stu == "LOGIN":
		u.uid, u.dt, u.event, u.sign = uid, this.dateTime, "login->login", 1
		return u
	case last.stu == "JOINGROUP" && this.stu == "LOGIN":
		u.uid, u.dt, u.gid, u.event, u.sign = uid, this.dateTime, append(u.gid, this.gid), "joinGroup->login", 1
		return u
	case last.stu == "RELOGIN" && this.stu == "LOGIN":
		u.uid, u.dt, u.event, u.sign = uid, this.dateTime, "relogin->login", 1
		return u
	case last.stu == "LOGOUT" && this.stu == "LOGIN":
		k = timeDifference(last.dateTime, this.dateTime)
		if k < offlineTimeSet {
			u.sign = 1
		} else {
			u.sign = 0
		}
		u.uid, u.dt, u.event = uid, this.dateTime, "logout->login"
		return u
	case last.stu == "" && this.stu == "LOGIN":
		u.uid, u.dt, u.event, u.sign = uid, this.dateTime, "->login", 0
		return u
	case last.stu == "LOGIN" && this.stu == "JOINGROUP":
		u.uid, u.dt, u.gid, u.event, u.sign = uid, this.dateTime, append(u.gid, this.gid), "login->joinGroup", 0
		return u
	case last.stu == "JOINGROUP" && this.stu == "JOINGROUP":
		u.gid, u.sign = append(u.gid, this.gid), 0
		return u
	case last.stu == "RELOGIN" && this.stu == "JOINGROUP":
		u.gid, u.sign = append(u.gid, this.gid), 0
		return u
	case last.stu == "LOGOUT" && this.stu == "JOINGROUP":
		k = timeDifference(last.dateTime, this.dateTime)
		if k < offlineTimeSet {
			u.sign = 1
		} else {
			u.sign = 0
		}
		u.uid, u.dt, u.gid, u.event = uid, this.dateTime, append(u.gid, this.gid), "logout->login"
		return u
	case last.stu == "" && this.stu == "JOINGROUP":
		u.uid, u.dt, u.event, u.sign = uid, this.dateTime, "->joinGroup", 0
		return u
	case last.stu == "LOGIN" && this.stu == "RELOGIN":
		u.uid, u.dt, u.event, u.sign = uid, this.dateTime, "login->relogin", 1
		return u
	case last.stu == "JOINGROUP" && this.stu == "RELOGIN":
		u.uid, u.dt, u.event, u.sign = uid, this.dateTime, "joinGroup->relogin", 1
		return u
	case last.stu == "RELOGIN" && this.stu == "RELOGIN":
		u.uid, u.dt, u.event, u.sign = uid, this.dateTime, "relogin->relogin", 1
		return u
	case last.stu == "LOGOUT" && this.stu == "RELOGIN":
		u.uid, u.dt, u.event, u.sign = uid, this.dateTime, "logout->relogin", 1
		return u
	case last.stu == "" && this.stu == "RELOGIN":
		u.uid, u.dt, u.event, u.sign = uid, this.dateTime, "->relogin", 1
		return u
	case last.stu == "LOGIN" && this.stu == "LOGOUT":
		u.uid, u.dt, u.event, u.sign = uid, this.dateTime, "login->logout", 0
		return u
	case last.stu == "JOINGROUP" && this.stu == "LOGOUT":
		u.uid, u.dt, u.event, u.sign = uid, this.dateTime, "joinGroup->logout", 0
		return u
	case last.stu == "RELOGIN" && this.stu == "LOGOUT":
		u.uid, u.dt, u.event, u.sign = uid, this.dateTime, "relogin->logout", 0
		return u
	case last.stu == "LOGOUT" && this.stu == "LOGOUT":
		u.uid, u.dt, u.event, u.sign = uid, this.dateTime, "logout->logout", 1
		return u
	case last.stu == "" && this.stu == "LOGOUT":
		u.uid, u.dt, u.event, u.sign = uid, this.dateTime, "->logout", 0
		return u
	}
	// fmt.Println(u)
	return u
}
func stuNameModify(stu string) string {
	if stu == "LOGOUT BROKEN" {
		return "LOGOUT"
	} else if stu == "JOIN GROUP" {
		return "JOINGROUP"
	}
	return stu

}
