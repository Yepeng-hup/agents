package conf

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

var (
	cfgAgent *Agent = nil
	ConfAgent *Agent
	cfgProcess *Process = nil
	ConfProcess *Process
)

type (
	Agent struct {
		LocalIp string `json:"local_ip"`
		ServerIp string `json:"server_ip"`
		ServerPort string `json:"server_port"`
		MonitorSecond int64 `json:"monitor_second"`
	}
	Process struct {
		ProcessName []string `json:"process_name"`
	}
)


func InitAgentConfig() (error,) {
	file, err := os.Open("conf/agent.json")
	if err != nil {
		log.Fatal("open json file: ", err.Error())
	}
	defer file.Close()
	f := bufio.NewReader(file)
	configObj := json.NewDecoder(f)
	if err = configObj.Decode(&cfgAgent); err != nil {
		return err
	}
	ConfAgent = cfgAgent
	return nil
}


func InitProcessConfig() (error,) {
	file, err := os.Open("conf/process.json")
	if err != nil {
		log.Fatal("open json file: ", err.Error())
	}
	defer file.Close()
	f := bufio.NewReader(file)
	configObj := json.NewDecoder(f)
	if err = configObj.Decode(&cfgProcess); err != nil {
		return err
	}
	ConfProcess = cfgProcess
	return nil
}