package pkg

import (
	"github.com/yinheli/qqwry"
	"strings"
)

func ipToAddr(ipStr string) string {
	var ipAddr string
	q := qqwry.NewQQwry("D://golang/golang/GOWORK/src/ip/qqwry.dat")
	if ipStr == "500" {
		ipAddr = "unknown"
	}
	q.Find(ipStr)
	ipAddr = q.Country + "-(" + q.City + ")"
	return ipAddr
}
func ipParse(s string) string {
	var m string
	a := strings.Index(s, ":")
	if a > 0 {
		m = s[:a]
	} else {
		m = "500"
	}
	return m
}
