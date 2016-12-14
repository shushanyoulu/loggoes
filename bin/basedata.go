package main

type Login struct {
	Ldate         string `json:"ldate"`
	Ltime         string `json:"ltime"`
	Stu           string `json:"stu"`
	Uid           string `json:"uid"`
	Address       string `json:"address"`
	IpAddr        string `json:"ipAddr"`
	Version       string `json:"version"`
	Platform      string `json:"Platform"`
	Device        string `json:"device"`
	ExpectPayload string `json:"expect_payload"`
	Os            string `json:"os"`
	Esn           string `json:"esn"`
	Meid          string `json:"meid"`
	Sn            string `json:"sn"`
	Imsi          string `json:"imsi"`
	System        string `json:"system"`
	Context       string `json:"context"`
	ServerNode    string `json:"ServerNode"`
	PerOfTime     string `json:"用户登陆"`
}
type onAndOff struct {
	PerOfTime  string `json:"用户退出"`
	UID        string `json:"uid"`
	ServerNode string `json:"ServerNode"`
	Ldate      string `json:"ldate"`
	Ltime      string `json:"ltime"`
	Stu        string `json:"stu"`
}
type onlineUsers struct {
	nodeName string
	TTL      int
}

//流数据，用于记录每个用户最近一次的状态
type streamData struct {
	dateTime string
	stu      string
	gid      string
	nodeName string
}

type perOneHour struct {
	ldatetime string
	gid       string
}
type perSixHour struct {
	ldatetime string
	gid       string
}
type perTwelevHour struct {
	ldatetime string
	gid       string
}
type logFormat struct {
	dateTime string
	uid      string
	stu      string
	gid      string
	nodeName string
}

// 用户最近一次数据状态
var userStreamData = make(map[string]streamData)
