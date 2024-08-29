package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
)

func getStatistic(w http.ResponseWriter, r *http.Request) {
	var listTrafInfo [][]string
	var sqlstring []string
	login := r.FormValue("login")
	pass := r.FormValue("pwd")

	e, err := newEmployee(login, pass)
	if err != nil {
		http.ServeFile(w, r, "./html/reject.html")
		fmt.Println("Авторизация сотрудника неуспешна\n", err)
	} else {
		e.ip = ReadUserIP(r)
		sqlstring = []string{
			"SELECT ",
			"$1 as n, ",
			"max(datetime) as datetime, ",
			"max(employee_ip_address) as employee_ip_address, ",
			"src_ip_address, ",
			"sum(traffic) as sumtraffic ",
			"FROM public.statistictrafficin ",
			"WHERE ",
			"login=$1 ",
			"and ",
			"datetime > date_trunc('day', current_date) ",
			"GROUP BY src_ip_address ",
			"UNION ALL ",
			"SELECT ",
			"'Итог' as n, ",
			"current_date as datetime, ",
			"'' as employee_ip_address, ",
			"'' as src_ip_address, ",
			"sum(traffic) as sumtraffic ",
			"FROM public.statistictrafficin ",
			"WHERE ",
			"login=$1 ",
			"and ",
			"datetime > date_trunc('day', current_date) ",
			"ORDER BY sumtraffic DESC",
		}
		rows, err := db.Query(strings.Join(sqlstring, ""), e.login)
		if err != nil {
			fmt.Println("Error request statistic ", err)
		}
		for rows.Next() { // result request db
			elemTrafficInfo := trafficInfo{time.Now(), "", "", "", "", 0}
			err = rows.Scan(&elemTrafficInfo.EmployeeLogin, &elemTrafficInfo.DateTime, &elemTrafficInfo.EmployeeIpAddress, &elemTrafficInfo.IpAddress, &elemTrafficInfo.Traffic)
			if err != nil {
				fmt.Println("Error Scan result row ", err)
			}
			var f float64
			f = float64(elemTrafficInfo.Traffic)
			f = f / 1024 / 1024
			elemSlice := []string{elemTrafficInfo.EmployeeLogin, elemTrafficInfo.EmployeeIpAddress, elemTrafficInfo.IpAddress, fmt.Sprintf("%.2f", f)}
			listTrafInfo = append(listTrafInfo, elemSlice)
		}
		rows.Close()

		t, _ := template.ParseFiles("./html/statisticDay.html")
		data := ViewData{
			FIO:      e.name,
			IP:       e.ip,
			Mb:       e.limittraffic - e.trafficIn,
			TrafStat: listTrafInfo,
		}
		//for _, elemTrafficInfo := range listTrafInfo {
		//	fmt.Println(elemTrafficInfo.IpAddress, "", elemTrafficInfo.traffic)
		//}

		t.Execute(w, data)
	}

}
