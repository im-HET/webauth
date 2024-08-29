package main

import (
	"fmt"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func netlimiter() {

	fmt.Println("Запушено копирование пакетов ")

	//с помощю библиотеки libpcap (устанавливаем в дебиан apt install libpcap-dev)
	//настраиваем подключение к ядру для считываения пакетов
	//параметры : устройство, количество копируемых байт с каждого пакета, и таймер работы
	//pcap.BlockForever означает ждать пакеты бесконечно
	h, err := pcap.OpenLive("enp0s20f1", 1500, false, pcap.BlockForever)
	if err != nil {
		fmt.Println("Ошибка создания хандлера ", err)
	}
	defer h.Close()
	//указываем параметры фильтра (синтаксис такой же как tcpdump)
	//в даннос случае отбираем только IPv4 пакеты
	if bogonNet != "" {
		err = h.SetBPFFilter("ip src net not " + bogonNet)
	} else {
		err = h.SetBPFFilter("ip")
	}
	if err != nil {
		fmt.Println("Ошибка создания фильтра BPF ", err)
	}

	//Запускаем считывание пакетов
	source := gopacket.NewPacketSource(h, h.LinkType())
	//в цикле считываем пакет из вышесозданного источника и обрабатываем его
	for p := range source.Packets() {
		//Проверяем есть ли в пакете слой с типом IPv4
		ip := p.Layer(layers.LayerTypeIPv4)
		if ip != nil {
			//Если такой слой есть то возвращаем его
			ip4, _ := ip.(*layers.IPv4)
			//если в списке наблюдаемых клиентов есть такой адрес то отправляем пакет в соответсвующий канал
			if e, ok := clients[ip4.DstIP.String()]; ok == true {
				e.trafficInChan <- ip4
			}
			if e, ok := clients[ip4.SrcIP.String()]; ok == true {
				e.trafficOutChan <- ip4
			}
		}
	}
	fmt.Println("Netlimiter завершился ")
}
