package main

import (
	"errors"
	"fmt"
	"time"

	"os/exec"

	"github.com/google/gopacket/layers"
)

type employee struct {
	id             uint32
	name           string
	role           string
	login          string
	pass           string
	limittraffic   uint32
	limitbitrate   uint32
	trafficIn      uint32
	trafficOut     uint32
	ip             string
	trafficInChan  chan *layers.IPv4
	trafficOutChan chan *layers.IPv4
	commandChan    chan string
	trafficStatIn  map[string]trafficInfo
	trafficStatOut map[string]trafficInfo
}

// конструктор
func newEmployee(login string, pass string) (employee, error) {
	var e employee
	e.ip = ""
	e.trafficInChan = make(chan *layers.IPv4, 1000)
	e.trafficOutChan = make(chan *layers.IPv4, 1000)
	e.trafficStatIn = make(map[string]trafficInfo)
	e.trafficStatOut = make(map[string]trafficInfo)
	e.commandChan = make(chan string)
	fmt.Println(login, " ", pass)
	row := db.QueryRow("select id, login, pass, limittraffic, limitbitrate, trafficin, name, role, trafficout from employee where login=$1 and pass=$2;", login, pass)
	err := row.Scan(&e.id, &e.login, &e.pass, &e.limittraffic, &e.limitbitrate, &e.trafficIn, &e.name, &e.role, &e.trafficOut)
	if err != nil {
		return e, errors.New("Ошибка, в базе данных нет учетных данных " + login + " " + pass + " : " + err.Error())
	}
	return e, nil
}

func (e employee) headerTraffic() {
	trafStat := newTrafficStatistics()
	e.addEmployeeToNftRule()
	var sumTrafficIn1M uint32 //переменная для сумирования длин пакетов
	var sumTrafficOut1M uint32
	sumTrafficIn1M = 0
	sumTrafficOut1M = 0
	fmt.Println("Доступно трафика ", e.limittraffic-e.trafficIn)
	end_for := true
	for end_for {
		select {
		case command := <-e.commandChan:
			if command == "disconnected" {
				fmt.Println("Клиент отключился ", e.ip, " ", e.login)
				end_for = false
				trafStat.flushIn()
				trafStat.flushOut()
			}
		case ip4 := <-e.trafficInChan:
			sumTrafficIn1M += uint32(ip4.Length)
			trafStat.addTrafficIn(e.login, e.ip, ip4)
			if sumTrafficIn1M > 1048576 { //если сумма длин более 1Мб
				sumTrafficIn1M = 0
				trafStat.flushIn()
				sumTrafficInInsertDBChan <- &e
			}
			if e.trafficIn >= e.limittraffic {
				fmt.Println("Лимит трафика исчерпан ", e.ip, " ", e.login, " ", e.trafficIn, "Mb", " из ", e.limittraffic)
				end_for = false
				trafStat.flushIn()
				trafStat.flushOut()
			}
		case ip4 := <-e.trafficOutChan:
			sumTrafficOut1M += uint32(ip4.Length)
			trafStat.addTrafficOut(e.login, e.ip, ip4)
			if sumTrafficOut1M > 1048576 { //если сумма длин более 1Мб
				sumTrafficOut1M = 0
				trafStat.flushOut()
				sumTrafficOutInsertDBChan <- &e
			}
		case <-time.After(timeoutClient * time.Minute): //<-time.After срабатывает через заданное количество времени
			fmt.Println("Нет пакетов более ", timeoutClient, " секунд с клиента ", e.ip)
			end_for = false
			trafStat.flushIn()
			trafStat.flushOut()
		}
	}
	e.delEmployeeFromNftRule()
	delEmployeeChannel <- &e
}

func (e employee) addEmployeeToNftRule() {
	//добавление ip адреса в множество nftables
	nft_add_ip := exec.Command("/sbin/nft", "add element ip nat ip_accept_list {"+e.ip+"}")
	err := nft_add_ip.Run()
	if err != nil {
		fmt.Println("Ошибка при добавлении адреса в множество nft ", err)
	}
}

func (e employee) delEmployeeFromNftRule() {
	//удаление ip адреса из множества в nftables
	nft_add_ip := exec.Command("/sbin/nft", "delete element ip nat ip_accept_list {"+e.ip+"}")
	err := nft_add_ip.Run()
	if err != nil {
		fmt.Println("Ошибка при удаления адреса из множество nft ", err)
	}
}
