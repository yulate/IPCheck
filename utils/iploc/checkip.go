package iploc

import (
	_ "embed"
)

var loc *Locator

func IpLocInit() {
	loc1, err := Open("data/czutf8.dat")
	if err != nil {
		panic(err)
	}
	loc = loc1
}

// Check 返回IP的运营商、当是云主机IP的时候会返回云主机运营商名字，如：腾讯云
func Check(ip string) string {
	detail := loc.Find(ip)
	return detail.Region
}
