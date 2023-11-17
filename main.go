package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/jeessy2/firewall-policy-api/ban"
)

func main() {
	//run the cmds in the switch, and get the execution results
	runWebServer()
}

func runWebServer() error {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/banByIP", ban.BanByIP)

	http.HandleFunc("/banByGrafana", ban.BanByGrafana)

	// 获取命令执行结果

	log.Println("监听", ":80", "...")

	l, err := net.Listen("tcp", ":80")
	if err != nil {
		return fmt.Errorf("监听端口发生异常, 请检查端口是否被占用: %w", err)
	}

	return http.Serve(l, nil)
}
