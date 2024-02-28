package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type (
	SystemInfo struct {
		Hostname string    `json:"hostname"`
		SysType  string    `json:"sys_type"`
		SysAddr  string    `json:"sys_addr"`
		Mem      Mem       `json:"mem"`
		Disk     Disk      `json:"disk"`
		Cpu      Cpu       `json:"cpu"`
		Net      []Network `json:"net"`
		Process  []Process `json:"process"`
	}

	Mem struct {
		MemTotal     float64 `json:"mem_total"`
		MemUse       float64 `json:"mem_use"`
		MemFree      float64 `json:"mem_free"`
		MemAvailable float64 `json:"mem_available"`
	}

	Cpu struct {
		CpuFree int `json:"cpu_free"`
	}

	Network struct {
		NetName string `json:"net_name"`
		NetRecv uint64 `json:"net_input"`
		NetSent uint64 `json:"net_ouput"`
	}

	Disk struct {
		DiskUsed  int     `json:"disk_used"`
		DiskFree  float64 `json:"disk_free"`
		DiskTotal float64 `json:"disk_total"`
	}

	Process struct {
		ProcessName string `json:"process_name"`
		ProcessPid  int32  `json:"process_pid"`
	}
)

func receiveData(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var data SystemInfo
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading JSON", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}
	fmt.Println("----------------------------------------------------------------------------------------")
	fmt.Println(data)

	// 响应客户端
	//w.WriteHeader(http.StatusOK)
	//if _, err := w.Write([]byte("JSON received and processed")); err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}

}

func main(){
	http.HandleFunc("/sys/info", receiveData)

	// 启动HTTP服务器
	log.Println("Start http-server and listen port 0.0.0.0:10100")
	if err := http.ListenAndServe(":10100", nil); err != nil {
		log.Fatal("Start error: ", err)
	}
}
