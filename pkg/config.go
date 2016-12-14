package pkg

import (
	"fmt"
	"github.com/BurntSushi/toml"
	// "strconv"
	// "time"
)

type config struct {
	Director direct
	DB       database `toml:"database"`
	LOG      loginfo  `toml:"loginfo"`
	ECHATLOG echatlog `toml:echatlog`
	Nodes    map[string]node
	Topics   map[string]topic
}
type direct struct {
	ReadFileFunction int
	UseKafkaStream   int
}
type database struct {
	Hostname string
	Port     string
	Driver   string
	Dbname   string
	Username string
	Password string
	ConnMax  int
	Enabled  bool
}
type echatlog struct {
	Echatlog string
}
type loginfo struct {
	Filename string
}
type node struct {
	Logpath string
	Loged   string
}
type topic struct {
	KafkaTopics   string
	ConsumerGroup string
	// clientPoolSize int
}

var URL string = "conf/base.toml"
var c config

func HostName() string {
	if _, err := toml.DecodeFile(URL, &c); err != nil {
		fmt.Println(err)
	}
	return c.DB.Hostname
}
func Port() string {
	if _, err := toml.DecodeFile(URL, &c); err != nil {
		fmt.Println(err)

	}
	return c.DB.Port
}
func Driver() string {
	if _, err := toml.DecodeFile(URL, &c); err != nil {
		fmt.Println(err)

	}
	return c.DB.Driver
}
func Dbname() string {
	if _, err := toml.DecodeFile(URL, &c); err != nil {
		fmt.Println(err)
	}
	return c.DB.Dbname
}
func Username() string {
	if _, err := toml.DecodeFile(URL, &c); err != nil {
		fmt.Println(err)
	}
	return c.DB.Username
}
func Password() string {
	if _, err := toml.DecodeFile(URL, &c); err != nil {
		fmt.Println(err)
	}
	return c.DB.Password
}
func ConnMax() int {
	if _, err := toml.DecodeFile(URL, &c); err != nil {
		fmt.Println(err)
	}
	return c.DB.ConnMax
}
func Enabled() bool {
	if _, err := toml.DecodeFile(URL, &c); err != nil {
		fmt.Println(err)
	}
	return c.DB.Enabled
}
func EchatLogPath() string {
	if _, err := toml.DecodeFile(URL, &c); err != nil {
		fmt.Println(err)
	}
	return c.ECHATLOG.Echatlog
}

var l config

func Nodes() (int, map[string]map[string]string) {
	ns := make(map[string]map[string]string)
	if _, err := toml.DecodeFile(URL, &l); err != nil {
		fmt.Println(err)
	}
	c := len(l.Nodes)
	for k, v := range l.Nodes {
		n := make(map[string]string)
		n["fh"] = v.Logpath
		n["ah"] = v.Loged
		ns[k] = n
	}
	return c, ns
}

func NodeCount() int {
	a, _ := Nodes()
	return a
}
func NodeInfo() map[string]map[string]string {
	_, b := Nodes()
	return b
}
func KafkaFunc() (int, map[string]map[string]string) {
	ns := make(map[string]map[string]string)
	if _, err := toml.DecodeFile(URL, &l); err != nil {
		fmt.Println(err)
	}
	c := len(l.Topics)
	for k, v := range l.Topics {
		n := make(map[string]string)
		n["kafkaTopics"] = v.KafkaTopics
		n["consumerGroup"] = v.ConsumerGroup
		ns[k] = n
	}
	return c, ns
}
func KafkaCount() int {
	a, _ := KafkaFunc()
	return a
}
func KafkaInfo() map[string]map[string]string {
	_, b := KafkaFunc()
	return b
}
