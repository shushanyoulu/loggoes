package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/wvanbergen/kafka/consumergroup"
	"github.com/wvanbergen/kazoo-go"
)

//并发队列数量
var dealTopicGroup sync.WaitGroup
var wg, waitGroup sync.WaitGroup
var bdata *broadLogData
var bulkRequest = connetEs() // es 连接地址
var ctx = context.Background()

type kafkaStruct struct {
	topicNames        string
	consumerGroupName string
}

type broadLogData struct {
	nodeName string
	data     string
}

type logTime struct {
	nodeName string
	dt       string
}

var ldt = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

var dataChan = make(chan broadLogData, 1000) //日志处理管道缓存数据量
var eUID = setExcludeUIDMap()                //要剔除的uid map

type dealLoger interface {
	periodOfUsersOffline()
	onlines()
	loginAndLogoutSendToEs()
	dailyUser()
}

func dealLog(d dealLoger) {
	d.periodOfUsersOffline()
	d.onlines()
	d.loginAndLogoutSendToEs()
	d.dailyUser()
}

func (k *kafkaStruct) dealNode() {
	go k.receiveTopicData()
	go dealTopicData()
}

// 并发处理每个topic数据
func analysisTopicGroup(nodesTopic map[string]topic) {
	var topicNames, consumerGroupName string
	var n *kafkaStruct
	for _, v := range nodesTopic {
		if topicNames == "" {
			topicNames = v.KafkaTopics
		} else {
			topicNames = topicNames + "," + v.KafkaTopics
		}
		consumerGroupName = v.ConsumerGroup
		n = &kafkaStruct{topicNames, consumerGroupName}
	}
	n.dealNode()
}

// 从每个topic 接受数据
func (k *kafkaStruct) receiveTopicData() {
	config := consumergroup.NewConfig()
	var nd broadLogData
	config.Offsets.Initial = sarama.OffsetNewest
	config.Offsets.ProcessingTimeout = 20 * time.Second
	zookeeperNodes, config.Zookeeper.Chroot = kazoo.ParseConnectionString(*zookeeper)
	kafkaTopics := strings.Split(k.topicNames, ",")
	goesl.Println(k.consumerGroupName, kafkaTopics, zookeeperNodes, config)
	consumer, consumerErr := consumergroup.JoinConsumerGroup(k.consumerGroupName, kafkaTopics, zookeeperNodes, config)
	if consumerErr != nil {
		goesl.Fatalln(consumerErr)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		if err := consumer.Close(); err != nil {
			sarama.Logger.Println("Error closing the consumer", err)
		}
	}()
	go func() {
		for err := range consumer.Errors() {
			goesl.Println(err)
		}
	}()

	eventCount := 0
	offsets := make(map[string]map[int32]int64)
	goesl.Println(consumer.InstanceRegistered())
	for message := range consumer.Messages() {
		if offsets[message.Topic] == nil {
			offsets[message.Topic] = make(map[int32]int64)
		}
		listenGetNum++
		nd.nodeName, nd.data = message.Topic, string(message.Value)
		dataChan <- nd
		eventCount++
		if offsets[message.Topic][message.Partition] != 0 && offsets[message.Topic][message.Partition] != message.Offset-1 {
			goesl.Printf("Unexpected offset on %s:%d. Expected %d, found %d, diff %d.\n", message.Topic, message.Partition, offsets[message.Topic][message.Partition]+1, message.Offset, message.Offset-offsets[message.Topic][message.Partition]+1)
		}
		offsets[message.Topic][message.Partition] = message.Offset
		consumer.CommitUpto(message)
	}
	goesl.Printf("Processed %d events.", eventCount)
	goesl.Printf("%+v", offsets)
	over <- "over"

}

// 处理节点数据
func dealTopicData() {
	for logData := range dataChan {
		logData.nodeLogLastTime()
		logData.classifyNodeLog()
	}
}
func (b broadLogData) nodeLogLastTime() {
	ldt.Lock()
	ldt.m[b.nodeName] = getTime(b.data)
	ldt.Unlock()
}

//将原始日志以类型分类
func (b broadLogData) classifyNodeLog() {
	log := b.data
	u := analysisUid(log)
	if _, ok := eUID[u]; ok == false {
		switch {
		case strings.Contains(log, "INFO"):
			dealLog(b)
		case strings.Contains(log, "DEBUG"):
		case strings.Contains(log, "ERROR"):
		case strings.Contains(log, "WARNING"):
		default:
			goesl.Println(log)
		}
	}
}
