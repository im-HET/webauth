package main

import (
	"fmt"
	"maps"
	"time"

	"github.com/google/gopacket/layers"
)

type trafficStatistics struct {
	listSrcInfo map[string]*trafficInfo
	listDstInfo map[string]*trafficInfo
}

func newTrafficStatistics() trafficStatistics {
	var trafStat trafficStatistics
	trafStat.listSrcInfo = make(map[string]*trafficInfo)
	trafStat.listDstInfo = make(map[string]*trafficInfo)
	return trafStat
}

func (t trafficStatistics) addTrafficIn(employeeLogin string, employeeIpAddress string, ip4 *layers.IPv4) {
	elemTrafficInfo, ok := t.listSrcInfo[ip4.SrcIP.String()]
	if ok == true {
		//fmt.Println("ip yes : ", ip4.SrcIP.String(), "", elemTrafficInfo.IpAddress, " ", uint32(ip4.Length), "", t.sum())
		elemTrafficInfo.Traffic += uint32(ip4.Length)
		elemTrafficInfo.DateTime = time.Now()

	} else {
		//fmt.Println("ip no : ", ip4.SrcIP.String(), "", elemTrafficInfo.IpAddress, " ", uint32(ip4.Length), "", t.sum())
		t.listSrcInfo[ip4.SrcIP.String()] = &trafficInfo{time.Now(), employeeLogin, employeeIpAddress, ip4.SrcIP.String(), "", uint32(ip4.Length)}
	}
}

func (t trafficStatistics) addTrafficOut(employeeLogin string, employeeIpAddress string, ip4 *layers.IPv4) {
	elemTrafficInfo, ok := t.listDstInfo[ip4.DstIP.String()]
	if ok == true {
		elemTrafficInfo.Traffic += uint32(ip4.Length)
		elemTrafficInfo.DateTime = time.Now()
	} else {
		t.listDstInfo[ip4.DstIP.String()] = &trafficInfo{time.Now(), employeeLogin, employeeIpAddress, ip4.DstIP.String(), "", uint32(ip4.Length)}
	}
}

func (t trafficStatistics) sum() uint32 {
	var result uint32
	result = 0
	for _, elem := range t.listSrcInfo {
		result += elem.Traffic
	}
	return result
}

func (t trafficStatistics) flushIn() {
	l := make(map[string]*trafficInfo)
	maps.Copy(l, t.listSrcInfo)
	trafficInInsertDBChan <- l
	clear(t.listSrcInfo)
}

func (t trafficStatistics) flushOut() {
	l := make(map[string]*trafficInfo)
	maps.Copy(l, t.listDstInfo)
	trafficOutInsertDBChan <- l
	clear(t.listDstInfo)
}

func (t trafficStatistics) print() {
	fmt.Println("--------------- Статистика за 1mb---------------------")
	for _, v := range t.listSrcInfo {
		fmt.Println(v.DateTime, v.IpAddress, v.DnsName, v.Traffic)
	}
	fmt.Println("------------------------------------------------------")
}

type trafficInfo struct {
	DateTime          time.Time
	EmployeeLogin     string
	EmployeeIpAddress string
	IpAddress         string
	DnsName           string
	Traffic           uint32
}
