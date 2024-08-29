package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func openDB() {
	var err error
	connectionString := "user=" + dbUser + " password=" + dbPassword + " host=" + dbHost + " dbname=netlimiter sslmode=disable"
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Соединение с базой данных выполнено")
}

func closeDB() {
	fmt.Println("DB close")
	db.Close()
}

var trafficInInsertDBChan chan map[string]*trafficInfo
var trafficOutInsertDBChan chan map[string]*trafficInfo
var sumTrafficInInsertDBChan chan *employee
var sumTrafficOutInsertDBChan chan *employee
var requestChan chan requestLogin

func insertTrafficStatistics() {
	fmt.Println("Запущщено копирование статистики в базу данных")
	trafficInInsertDBChan = make(chan map[string]*trafficInfo, 1000)
	trafficOutInsertDBChan = make(chan map[string]*trafficInfo, 1000)
	sumTrafficInInsertDBChan = make(chan *employee)
	sumTrafficOutInsertDBChan = make(chan *employee)
	requestChan = make(chan requestLogin)
	for {
		select {
		case listSrcInfo := <-trafficInInsertDBChan:
			if len(trafficInInsertDBChan) > 300 {
				fmt.Println("buffer trafficInInsertDBChan > 300 :", len(trafficInInsertDBChan))
			}
			for _, elemTrafficInfo := range listSrcInfo {
				sql := "insert into statistictrafficin (datetime, login, employee_ip_address, src_ip_address, src_dns, traffic) values ($1, $2, $3, $4, $5, $6)"
				_, err := db.Exec(sql, elemTrafficInfo.DateTime, elemTrafficInfo.EmployeeLogin, elemTrafficInfo.EmployeeIpAddress, elemTrafficInfo.IpAddress, elemTrafficInfo.DnsName, elemTrafficInfo.Traffic)
				if err != nil {
					fmt.Println("Ошибка записи в базу ", err)
				}
			}
		case listDstInfo := <-trafficOutInsertDBChan:
			if len(trafficOutInsertDBChan) > 300 {
				fmt.Println("buffer trafficOutInsertDBChan > 300 :", len(trafficOutInsertDBChan))
			}
			for _, elemTrafficInfo := range listDstInfo {
				sql := "insert into statistictrafficout (datetime, login, employee_ip_address, dst_ip_address, dst_dns, traffic) values ($1, $2, $3, $4, $5, $6)"
				_, err := db.Exec(sql, elemTrafficInfo.DateTime, elemTrafficInfo.EmployeeLogin, elemTrafficInfo.EmployeeIpAddress, elemTrafficInfo.IpAddress, elemTrafficInfo.DnsName, elemTrafficInfo.Traffic)
				if err != nil {
					fmt.Println("Ошибка записи в базу ", err)
				}
			}
		case e := <-sumTrafficInInsertDBChan:
			sql := "update employee set trafficin = trafficin+1 where id=$1 returning trafficin"
			result := db.QueryRow(sql, e.id)
			err := result.Scan(&e.trafficIn)
			if err != nil {
				fmt.Println("Ошибка прибавки траффика ", e.login, " ", err)
			}
		case e := <-sumTrafficOutInsertDBChan:
			sql := "update employee set trafficout = trafficout+1 where id=$1 returning trafficout"
			result := db.QueryRow(sql, e.id)
			err := result.Scan(&e.trafficOut)
			if err != nil {
				fmt.Println("Ошибка прибавки траффика ", e.login, " ", err)
			}
		case r := <-requestChan:
			sql := "insert into requestlogin (fio, org, tel, datetime) values ($1, $2, $3, $4)"
			_, err := db.Exec(sql, r.fio, r.org, r.tel, time.Now())
			if err != nil {
				fmt.Println("Ошибка записи в базу ", err)
			}
		}
	}
}
