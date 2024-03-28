package core

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var lowLevel []int
var middenLevel []int
var highLevel []int

func init() {
	lowLevel = []int{80, 443}
	middenLevel = []int{80, 443, 7000, 8080, 8081, 8443}
	highLevel = []int{21, 22, 23, 80, 81, 82, 88, 8000, 8888, 888, 443, 8443, 5000, 7000}
}

func checkPorts(ip string, ports []int) []int {
	openPorts := []int{}
	var wg sync.WaitGroup
	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			if checkPort(ip, p) {
				openPorts = append(openPorts, p)
			}
		}(port)
	}
	wg.Wait()
	return openPorts
}

func checkPort(ip string, port int) bool {
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

// CheckPort 传入IP、扫描等级：low、midden、high，返回开放端口
func CheckPort(ip string, level string) []int {
	var ports []int
	switch level {
	case "midden":
		ports = middenLevel
		break
	case "high":
		ports = highLevel
		break
	default:
		ports = lowLevel
	}
	return checkPorts(ip, ports)
}
