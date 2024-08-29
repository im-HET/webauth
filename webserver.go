package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
)

func start() {
	muxHs := http.NewServeMux()
	muxHs.HandleFunc("/", index)
	muxHs.HandleFunc("/auth", auth)
	muxHs.HandleFunc("/request", request)
	muxHs.HandleFunc("/disconnect", disconnect)
	muxHs.HandleFunc("/getstatisticday", getStatistic)

	go startHttpServerWithCatchErr(muxHs)
	fmt.Println("Запущен WebServer на порту :8080")
}

func auth(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	pass := r.FormValue("pwd")

	e, err := newEmployee(login, pass)
	if err != nil {
		http.ServeFile(w, r, "./html/reject.html")
		fmt.Println("Авторизация сотрудника неуспешна\n", err)
	} else {
		e.ip = ReadUserIP(r)
		if emp, ok := clients[e.ip]; ok == true {
			fmt.Println("На ip адресе ", emp.ip, " уже есть авторизованный пользователь ", emp.login)
			emp.commandChan <- "disconnected"
			time.Sleep(3 * time.Second)
		}
		addEmployeeChannel <- &e //add employee to netlimiter
		go e.headerTraffic()
		fmt.Println("Сотрудник ", e.login, " успешно авторизовался с адресом ", e.ip)
		t, _ := template.ParseFiles("./html/access.html")
		data := ViewData{
			FIO: e.name,
			IP:  e.ip,
			Mb:  e.limittraffic - e.trafficIn,
		}
		t.Execute(w, data)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Подключился клиент с адреса ", ReadUserIP(r))
	http.ServeFile(w, r, "./html/index.html")
}

func request(w http.ResponseWriter, r *http.Request) {
	fio := r.FormValue("fio")
	org := r.FormValue("org")
	tel := r.FormValue("tel")
	fmt.Println("Поступил запрос от ", fio, " Организация: ", org, "Тел: ", tel)
	requestChan <- requestLogin{fio, org, tel}
	http.ServeFile(w, r, "./html/request.html")
}

func disconnect(w http.ResponseWriter, r *http.Request) {
	ip := ReadUserIP(r)
	data := ViewData{
		FIO: "not connected",
		IP:  ip,
		Mb:  0,
	}
	t, _ := template.ParseFiles("./html/disconnect.html")
	if e, ok := clients[ip]; ok == true {
		fmt.Println("request disconnect ", e.ip, " ", e.login)
		e.commandChan <- "disconnected"
		time.Sleep(3 * time.Second)
		data = ViewData{
			FIO: e.name,
			IP:  e.ip,
			Mb:  e.limittraffic - e.trafficIn,
		}
	}
	t.Execute(w, data)
	http.ServeFile(w, r, "./html/disconnect.html")
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	port := strings.Index(IPAddress, ":")
	return IPAddress[:port]
}

func startHttpServerWithCatchErr(mux *http.ServeMux) {
	err := http.ListenAndServe(":"+webServerPort, mux)
	if err != nil {
		fmt.Println(err)
	}
}

type requestLogin struct {
	fio string
	org string
	tel string
}

type ViewData struct {
	FIO      string
	IP       string
	Mb       uint32
	TrafStat [][]string
}
