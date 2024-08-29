package main

import (
	"fmt"
)

var clients map[string]*employee

var addEmployeeChannel chan *employee
var delEmployeeChannel chan *employee

func clientsManager() {
	clients = make(map[string]*employee)
	addEmployeeChannel = make(chan *employee)
	delEmployeeChannel = make(chan *employee)
	for {
		select {
		case addIpChan := <-addEmployeeChannel:
			clients[addIpChan.ip] = addIpChan
		case delIpChan := <-delEmployeeChannel:
			fmt.Println("Удаляем клиента ", delIpChan.login, " от отслеживания")
			delete(clients, delIpChan.ip)
		}
	}
}
