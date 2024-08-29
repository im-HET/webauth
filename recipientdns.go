package main

/*
import (
	"fmt"
	// "os/exec"
	// "regexp"
	// "strings"
)

var getDNSChannel chan map[string]*srcInfo

func recipientdns() {
	fmt.Println("Опросник ДНС запущен")
	listDNS := make(map[string]string)
	for _, localhost := range localhostIP {
		listDNS[localhost] = "gateway.buhta-sever.ru"
	}
	getDNSChannel = make(chan map[string]*srcInfo, 1000)
	for {
		listSrcInfo := <-getDNSChannel
		//for ipsrc, elemSrcInfo := range listSrcInfo {
		//for _, elemSrcInfo := range listSrcInfo {
		//DNS, ok := listDNS[ipsrc]
		//if ok == true {
		//	elemSrcInfo.dnsName = DNS
		//} else {
		//	out, err := exec.Command("nslookup", ipsrc).Output()
		//	if err != nil {
		//		fmt.Println("Ошибка nslookup ", ipsrc, " buf ", len(getDNSChannel), err)

		//	} else {
		//		var re = regexp.MustCompile("[[:space:]]")
		//		DNS = re.ReplaceAllString(string(out[strings.LastIndex(string(out), "name")+7:]), "")
		//		elemSrcInfo.dnsName = DNS
		//		listDNS[ipsrc] = DNS
		//	}
		//}
		//trafficStatisticsInsertDBChannel <- elemSrcInfo
		//}
		//fmt.Println("Буфер записи в базу ", len(getDNSChannel))
	}
}
*/
