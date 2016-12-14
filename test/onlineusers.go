package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

type onlineUsers struct {
	nodeName string
	TTL      int
}

// type statisticNodeUsers struct {
// 	node    string
// 	userMap *Map
// }
var statisticNodeUsers = make(map[string]*Map)

func main() {
	options := &Options{
		InitialCapacity: 1024,
		OnWillExpire: func(key string, item *Item) {
		},
		OnWillEvict: func(key string, item *Item) {
		},
	}
	m := New(options)
	// don't forget to drain the map when you don't need it
	// defer m.Drain()
	o, _ := os.Open("log.txt")
	defer o.Close()
	f := bufio.NewReader(o)
	for i := 0; i < 10; i++ {
		l, _ := f.ReadString('\n')
		uid := analysisUid(l)
		switch {
		case strings.Contains(l, "LOGIN"):
			m.SetNX(uid, NewItemWithTTL("test", 5*time.Second))
		case strings.Contains(l, "LOGOUT"):
			m.Delete(uid)
			statisticNodeUsers["test"] = m
		}
	}
	for k, _ := range statisticNodeUsers {
		fmt.Println(k, statisticNodeUsers[k].Len())
	}

}

func (m *Map) searchMap() {
	s := m.store
	for _, pqi := range s.pq {
		if s.onWillEvict != nil {
			fmt.Println(pqi.key, pqi.item.Value(), pqi.item.TTL())
		}
	}
}

func analysisUid(l string) string {
	uidt0 := regexp.MustCompile(`uid=\(\d+\)`)
	uidt1 := uidt0.FindString(l)
	uidt2 := regexp.MustCompile(`\d+`)
	uid := uidt2.FindString(uidt1)
	return uid
}

func analysisStu(l string) string {
	j := strings.Index(l, "uid")
	if j >= 2 {
		stu := l[:j-1]
		return stu
	}
	return "other"
}
