package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"http-agnet/conf"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
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

func showMemInfo() Mem {
	m, _ := mem.VirtualMemory()
	total := float64(m.Total/1024/1024) / float64(1024)
	used := float64(m.Used/1024/1024) / float64(1024)
	free := float64(m.Free/1024/1024) / float64(1024)
	available := float64(m.Available/1024/1024) / float64(1024)
	memory := Mem{
		MemTotal:     float64(int(total*100)) / 100,
		MemUse:       float64(int(used*100)) / 100,
		MemFree:      float64(int(free*100)) / 100,
		MemAvailable: float64(int(available*100)) / 100,
	}
	return memory
}

func showDiskInfo() Disk {
	d, _ := disk.Usage("/")
	disks := Disk{
		DiskUsed:  int(d.Used / 1024 / 1024 / 1024),
		DiskFree:  float64(d.Free / 1024 / 1024 / 1024),
		DiskTotal: float64(d.Total / 1024 / 1024 / 1024),
	}
	return disks
}

func showCpuInfo() Cpu {
	c2, _ := cpu.Percent(time.Duration(time.Second), false)
	c := Cpu{
		CpuFree: 100 - int(c2[0]),
	}
	return c
}

func showNetworkInfo() []Network {
	allNetworkList := make([]Network, 0)
	netIO, err := net.IOCounters(true)
	if err != nil {
		fmt.Println("Failed to obtain network interface information:", err)
		return nil
	}

	// show all network interface info and status
	for _, io := range netIO {
		n := Network{
			NetName: io.Name,
			NetRecv: io.BytesRecv,
			NetSent: io.BytesSent,
		}
		allNetworkList = append(allNetworkList, n)
	}
	return allNetworkList
}

func showProcessInfo() []Process {
	initProcessNameList := make([]Process, 0)
	processNameList := conf.ConfProcess.ProcessName
	processes, err := process.Processes()
	if err != nil {
		log.Println("ERROR: get process info fail,", err.Error())
	}

	for _, p := range processes {
		name, _ := p.Name()
		pid := p.Pid
		for _, v := range processNameList {
			if name == v {
				p := Process{
					ProcessName: name,
					ProcessPid:  pid,
				}
				initProcessNameList = append(initProcessNameList, p)
				break
			}
		}

	}
	return initProcessNameList
}

func CollectSystemInfo() SystemInfo {
	hostName, _ := os.Hostname()
	osName := runtime.GOOS
	sysAddr := conf.ConfAgent.LocalIp
	mem_ := showMemInfo()
	cpu_ := showCpuInfo()
	disk_ := showDiskInfo()
	net_ := showNetworkInfo()
	process_ := showProcessInfo()

	info := SystemInfo{
		Hostname: hostName,
		SysType:  osName,
		SysAddr:  sysAddr,
		Mem:      mem_,
		Cpu:      cpu_,
		Disk:     disk_,
		Net:      net_,
		Process:  process_,
	}
	return info
}

func SendSystemInfo(url string, info SystemInfo) error {
	// convert json format
	jsonData, err := json.Marshal(info)
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	return nil
}
