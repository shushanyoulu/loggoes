package main

import (
	"regexp"
	"runtime"
	"strings"
	"time"

	"gopkg.in/olivere/elastic.v5"
)

// use cpu
func uCPUNum() int {
	u := useCPUS()
	if u == 0 {
		u = runtime.NumCPU()
		return u
	}
	return u
}

//将日志文件分割成时间，信息2部分；
func extract(line string) (string, string, string) {
	logtimeFormat := regexp.MustCompile(`(?P<datetime>\d\d\d\d[-|/]\d\d[-|/]\d\d\s\d\d:\d\d:\d\d\.\d+)`)
	logdatetime := logtimeFormat.FindString(line)
	loglineFormat := regexp.MustCompile(`".*`)
	a := loglineFormat.FindString(line)
	if len(a)-1 > 1 {
		loginfo := a[1 : len(a)-1]
		ldate := strings.Fields(logdatetime)[0]
		ltime := strings.Fields(logdatetime)[1]
		return ldate, ltime, loginfo
	}
	goesl.Println("error")
	return "", "", ""
}
func getTime(line string) string {
	logTimeFormat := regexp.MustCompile(`(?P<datetime>\d\d\d\d[-|/]\d\d[-|/]\d\d\s\d\d:\d\d:\d\d\.\d+)`)
	logDateTime := logTimeFormat.FindString(line)
	if len(logDateTime) > 19 {
		return logDateTime[:19]
	}
	goesl.Println("log date  error")
	return ""
}
func timeChangeToEsFormat(s string) string {
	schange := strings.Replace(s, "/", "-", -1)
	a := strings.Fields(schange)
	return a[0] + "T" + a[1] + "+0800"
}
func checkerr(err error) {
	if err != nil {
		// fmt.Println(err)
		goesl.Println(err)
	}
}

func gotNum(n uint32) {
	logi = n / 10000
	if logi > logj {
		logj = logi
		goesl.Printf(" had receviced %d !\n", logj*10000)
	}
}
func trimf(s string) string {
	a := strings.Index(s, "(")
	b := strings.Index(s, ")")
	if a < 0 || b < 0 {
		return ""
	}
	return s[a+1 : b]

}
func period(ld, lt string) string {
	ldchange := strings.Replace(ld, "/", "-", -1)
	p := ldchange + "T" + lt[:8] + "+0800"
	return p
}

func connetEs() *elastic.BulkService {
	goesl.Println("连接elasticsearch地址为： ", esAddr())
	client, err := elastic.NewClient(elastic.SetURL(esAddr()))
	if err != nil {
		goesl.Println(err)
		panic(err)
	}
	bulkRequest := client.Bulk()
	return bulkRequest

}

//两个时间做差
func timeDifference(a, b string) float64 {
	timeLayout := "2006/01/02 15:04:05"                     //待转化为时间戳的字符串 注意 这里的小时和分钟还要秒必须写 因为是跟着模板走的 修改模板的话也可以不写
	loc, _ := time.LoadLocation("Local")                    //重要：获取时区
	theTimeA, _ := time.ParseInLocation(timeLayout, a, loc) //使用模板在对应时区转化为time.time类型
	theTimeB, _ := time.ParseInLocation(timeLayout, b, loc) //使用模板在对应时区转化为time.time类型
	aTob := theTimeB.Sub(theTimeA)
	ab := aTob.Seconds()
	return ab
}
