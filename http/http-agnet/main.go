package main

import (
	"http-agnet/conf"
	"log"
	"time"
	"http-agnet/agent"
)

func main()  {
	if err := conf.InitAgentConfig(); err != nil {
		log.Fatalln(err.Error())
	}
	if err := conf.InitProcessConfig(); err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("INFO: agent start run.")
	monitorServerURL := conf.ConfAgent.ServerSsl+"://"+conf.ConfAgent.ServerIp+":"+conf.ConfAgent.ServerPort+conf.ConfAgent.ServerApi

	// start goroutine run corntab
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

	// blocking goroutine
	select {}
}
