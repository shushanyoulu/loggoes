package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type config struct {
	Director    direct
	SERVER      server      `toml:"server"`
	CONCURRENCY concurrency `toml:"concurrency"`
	Topics      map[string]topic
	LogConfig   logConfig
	ExcludeUID  excludeUID
}
type direct struct {
	ReadFileFunction int
	UseKafkaStream   int
}
type concurrency struct {
	Tcount int
	// Putbuff     int
	PutToEsBuff int
}
type topic struct {
	// PttsvcName    string
	KafkaTopics   string
	ConsumerGroup string
	// clientPoolSize int
}
type server struct {
	Zookeeperaddr    string
	Esaddr           string
	LoginInfoIdx     string
	LogOutIdx        string
	Cpus             int
	UserStatisticTTL int
	OnlineuserUpdate int
	OfflineInterval  int
}
type logConfig struct {
	LogPath string
}
type excludeUID struct {
	Uids string
}

var l config
var configPath = "../conf/goes.toml"
var c config

func zookeeperAddr() *string {
	if s, err := toml.DecodeFile(configPath, &c); err != nil {
		fmt.Println(err)
		fmt.Println(s)
	}
	if c.SERVER.Zookeeperaddr == "" {
		fmt.Println("zookeeper is empty ! please check it !")
		os.Exit(2)
	}
	return &c.SERVER.Zookeeperaddr
}
func esAddr() string {
	if _, err := toml.DecodeFile(configPath, &c); err != nil {
		fmt.Println(err)

	}
	return c.SERVER.Esaddr
}

func useCPUS() int {
	if _, err := toml.DecodeFile(configPath, &c); err != nil {
		fmt.Println(err)
	}
	return c.SERVER.Cpus
}
func userStatisticTTL() int {
	if _, err := toml.DecodeFile(configPath, &c); err != nil {
		fmt.Println(err)
	}
	return c.SERVER.UserStatisticTTL
}
func onlineUserUpdate() int {
	if _, err := toml.DecodeFile(configPath, &c); err != nil {
		fmt.Println(err)
	}
	return c.SERVER.OnlineuserUpdate
}

// 掉线时间间隔设置
func offlineInterval() int {
	if _, err := toml.DecodeFile(configPath, &c); err != nil {
		fmt.Println(err)
	}
	return c.SERVER.OfflineInterval
}
func tcount() int {
	if _, err := toml.DecodeFile(configPath, &c); err != nil {
		fmt.Println(err)
	}
	return c.CONCURRENCY.Tcount
}

func putToEsBuff() int {
	if _, err := toml.DecodeFile(configPath, &c); err != nil {
		fmt.Println(err)
	}
	return c.CONCURRENCY.PutToEsBuff
}
func ktopic() map[string]topic {
	if _, err := toml.DecodeFile(configPath, &c); err != nil {
		fmt.Println(err)
	}
	return c.Topics
}
func logPath() string {
	if _, err := toml.DecodeFile(configPath, &c); err != nil {
		fmt.Println(err)
	}
	return c.LogConfig.LogPath
}
func deleteUID() string {
	if _, err := toml.DecodeFile(configPath, &c); err != nil {
		fmt.Println(err)
	}
	return c.ExcludeUID.Uids
}

// func main() {
// 	fmt.Println(deleteUID())
// }
func notAnalysisUID() []string {
	us := deleteUID()
	uidSlice := strings.Split(us, ";")
	return uidSlice
}
func setExcludeUIDMap() map[string]int {
	s := notAnalysisUID()
	m := len(s)
	excludeUIDMap := make(map[string]int)
	for i := 0; i < m; i++ {
		excludeUIDMap[s[i]] = i
	}
	return excludeUIDMap
}
