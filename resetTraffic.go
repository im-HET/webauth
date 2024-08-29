package main

import (
	"fmt"
	"time"
)

func resetTraffic() {
	for {
		<-time.After(5 * time.Minute)
		t := time.Now()
		if t.Hour() == resetTime && t.Minute() < 5 {
			sql := "update employee set trafficin = 0, trafficout = 0"
			_, err := db.Exec(sql)
			if err != nil {
				fmt.Println("Ошибка сброса данных траффика ", err)
			} else {
				fmt.Println("Данные трафика сброшены на 0 ", time.Now())
			}
		}
	}
}
