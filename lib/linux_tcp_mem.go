package linuxTcp

import (
	"bufio"
	"os"
	"os/exec"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin"
)

var graphdef map[string]mp.Graphs = map[string]mp.Graphs{
	"linutTcp.tcpMem": {
		Label: "Linux Tcp Memory",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "tcp_memory_size", Label: "tcp_memory", Diff: false},
		},
	},
	"linutTcp.tcpMemMin": {
		Label: "Linux Tcp Memory",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "tcp_memory_size_min", Label: "tcp_memory", Diff: false},
		},
	},
	"linutTcp.tcpMemPressure": {
		Label: "Linux Tcp Memory",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "tcp_memory_size_pressure", Label: "tcp_memory", Diff: false},
		},
	},
	"linutTcp.tcpMemMax": {
		Label: "Linux Tcp Memory",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "tcp_memory_size_max", Label: "tcp_memory", Diff: false},
		},
	},
	"linutTcp.udpMem": {
		Label: "Linux udp Memory",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "udp_memory_size", Label: "udp_memory", Diff: false},
		},
	},
	"linutTcp.udpMemMin": {
		Label: "Linux udp Memory",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "udp_memory_size_min", Label: "udp_memory", Diff: false},
		},
	},
	"linutTcp.udpPressure": {
		Label: "Linux udp Memory",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "udp_memory_size_pressure", Label: "udp_memory", Diff: false},
		},
	},
	"linutTcp.udpMemMax": {
		Label: "Linux udp Memory",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "udp_memory_size_max", Label: "udp_memory", Diff: false},
		},
	},
}

type LinuxTcpMemPlugin struct {
	Target string
}

var pagesize int = os.Getpagesize()

type ProtoMemPages struct {
	tcpMem         int
	tcpMemMin      int
	tcpMemPressure int
	tcpMemMax      int
	udpMem         int
	udpMemMin      int
	udpMemPressure int
	udpMemMax      int
}

func (pmp *ProtoMemPages) parseProtobufMemorytTreshhold() {
	var arr []string
	for _, param := range []string{"net.ipv4.tcp_mem", "net.ipv4.udp_mem"} {
		res, err := exec.Command("sysctl", param).Output()
		arr = strings.Split(string(res), "\t")
		if err != nil {
			panic(err)
		}
		switch param {
		case "net.ipv4.tcp_mem":
			pmp.tcpMemMin, _ = strconv.Atoi(strings.Split(arr[0], " ")[2])
			pmp.tcpMemPressure, _ = strconv.Atoi(arr[1])
			pmp.tcpMemMax, _ = strconv.Atoi(strings.TrimRight(arr[2], "\n"))
		case "net.ipv4.udp_mem":
			pmp.udpMemMin, _ = strconv.Atoi(strings.Split(arr[0], " ")[2])
			pmp.udpMemPressure, _ = strconv.Atoi(arr[1])
			pmp.udpMemMax, _ = strconv.Atoi(strings.TrimRight(arr[2], "\n"))
		}
	}
}

func parseSockstatMem() *ProtoMemPages {
	var txt string
	pmp := new(ProtoMemPages)
	pmp.parseProtobufMemorytTreshhold()

	fp, err := os.Open("/proc/net/sockstat")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		txt = scanner.Text()
		if strings.Contains(txt, "TCP") {
			arr := strings.Split(txt, " ")
			pmp.tcpMem, _ = strconv.Atoi(arr[len(arr)-1])
		} else if strings.Contains(txt, "UDP") && !strings.Contains(txt, "LITE") {
			arr := strings.Split(txt, " ")
			pmp.udpMem, _ = strconv.Atoi(arr[len(arr)-1])
		}
	}
	return pmp
}

func (ltmp LinuxTcpMemPlugin) FetchMetrics() (map[string]float64, error) {
	pmp := parseSockstatMem()
	m := map[string]float64{
		"tcp_memory_size":          float64(pmp.tcpMem * pagesize),
		"tcp_memory_size_min":      float64(pmp.tcpMemMin * pagesize),
		"tcp_memory_size_pressure": float64(pmp.tcpMemPressure * pagesize),
		"tcp_memory_size_max":      float64(pmp.tcpMemMax * pagesize),
		"udp_memory_size":          float64(pmp.udpMem * pagesize),
		"udp_memory_size_min":      float64(pmp.udpMemMin * pagesize),
		"udp_memory_size_pressure": float64(pmp.udpMemPressure * pagesize),
		"udp_memory_size_max":      float64(pmp.udpMemMax * pagesize),
	}
	return m, nil
}

func (ltmp LinuxTcpMemPlugin) GraphDefinition() map[string]mp.Graphs {
	return graphdef
}

func Do() {
	var ltmp LinuxTcpMemPlugin
	helper := mp.NewMackerelPlugin(ltmp)
	helper.Run()
}
