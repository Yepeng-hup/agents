package main

import (
	"http-agnet/conf"
	"log"
	"time"
	"http-agnet/agent"
)

func main()  {
	err0 := conf.InitAgentConfig()
	if err0 != nil {
		log.Fatalln(err0.Error())
	}
	err1 := conf.InitProcessConfig()
	if err1 != nil {
		log.Fatalln(err1.Error())
	}
	log.Println("INFO: agent start run.")
	monitorServerURL := conf.ConfAgent.ServerSsl+"://"+conf.ConfAgent.ServerIp+":"+conf.ConfAgent.ServerPort+conf.ConfAgent.ServerApi

	go func() {
		ticker := time.NewTicker(time.Duration(conf.ConfAgent.MonitorSecond) * time.Second)
		for range ticker.C {
			info := agent.CollectSystemInfo()
			err := agent.SendSystemInfo(monitorServerURL, info)
			if err != nil {
				log.Println("ERROR: sending system info:", err.Error())
			}else {
				log.Println("INFO: send server success")
			}
		}
	}()

	// 阻塞
	select {}
}
