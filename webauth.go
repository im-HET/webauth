package main

import (
	"time"
)

// Конфигурационные переменные
var timeoutClient time.Duration
var dbUser string
var dbPassword string
var dbHost string
var webServerPort string
var localhostIP []string
var resetTime int
var bogonNet string

func main() {
	configure()
	openDB()
	go resetTraffic()
	go insertTrafficStatistics()
	go clientsManager()
	go netlimiter()
	start()
	<-make(chan struct{})
	closeDB()
}

func configure() {
	timeoutClient = 100
	dbUser = "netlimiter"
	dbPassword = "netlimiter"
	dbHost = "127.0.0.1"
	webServerPort = "8080"
	localhostIP = []string{"10.100.0.1"}
	resetTime = 4
	bogonNet = "192.168.13.0/24"
}
