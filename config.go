package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type config struct {
	Director    direct
	SERVER      server      `toml:"server"`
	CONCURRENCY concurrency `toml:"concurrency"`
	Topics      map[string]topic
	LogConfig   logConfig
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
}
type logConfig struct {
	LogPath string
}

var l config
var configPath string = "conf/goes.toml"
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
func esLoginInfoIndex() string {
	if _, err := toml.DecodeFile(configPath, &c); err != nil {
		fmt.Println(err)

	}
	return c.SERVER.LoginInfoIdx
}
func esLogoutIndex() string {
	if _, err := toml.DecodeFile(configPath, &c); err != nil {
		fmt.Println(err)

	}
	return c.SERVER.LogOutIdx
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

// func main() {
// 	fmt.Println(logPath())
// }
