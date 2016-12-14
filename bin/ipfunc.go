package main

import (
	"strings"

	"github.com/yinheli/qqwry"
)

// ip change to addr
func ipToAddr(ipStr string) string {
	var ipAddr string
	q := qqwry.NewQQwry("../iplib/qqwry.dat")
	if ipStr == "500" {
		ipAddr = "unknown"
		return ipAddr
	}
	q.Find(ipStr)
	ipAddr = q.Country + "-(" + q.City + ")"
	return ipAddr
}

// ip  parse
func ipParse(s string) string {
	s = trimf(s)
	a := strings.Index(s, ":")
	if a > 0 {
		s = s[:a]
		return s
	}
	return "500"
}
