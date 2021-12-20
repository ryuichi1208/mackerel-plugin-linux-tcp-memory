package linuxTcp

import (
	"bufio"
	"fmt"
	"os"
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
	"linutTcp.udpMem": {
		Label: "Linux udp Memory",
		Unit:  "integer",
		Metrics: []mp.Metrics{
			{Name: "udp_memory_size", Label: "udp_memory", Diff: false},
		},
	},
}

type LinuxTcpMemPlugin struct {
	Target string
}

var pagesize int = os.Getpagesize()

type ProtoMemPages struct {
	tcpMem int
	udpMem int
}

func parseSockstatMem() *ProtoMemPages {
	var txt string
	pmp := new(ProtoMemPages)

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
	fmt.Println(pmp)
	m := map[string]float64{
		"tcp_memory_size": float64(pmp.tcpMem * pagesize),
		"udp_memory_size": float64(pmp.udpMem * pagesize),
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
