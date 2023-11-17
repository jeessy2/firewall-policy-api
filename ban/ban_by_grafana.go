package ban

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/jeessy2/firewall-policy-api/entity"
)

// BanByGrafana 通过Grafana告警Ban IP
func BanByGrafana(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var ga entity.GrafanaAlert
	err := json.NewDecoder(r.Body).Decode(&ga)
	if err != nil {
		returnError(w, fmt.Errorf("数据格式不正确: %s", err))
		return
	}

	ips := make([]string, 0)
	for i := 0; i < len(ga.Alerts); i++ {
		alert := ga.Alerts[i]
		if alert.Status == "firing" {
			sp := strings.Split(alert.Annotations.Description, "\n")
			banIPs, err := checkIP(sp)
			if err != nil {
				log.Println(err)
			}
			ips = append(ips, banIPs...)
		}
	}

	// ban IP
	if len(ips) == 0 {
		returnError(w, fmt.Errorf("Ban的IP为空, Request: %s", ga))
		return
	}
	log.Printf("将要Ban的IP有: %s\n", ips)
	err = banIP(ips)

	if err != nil {
		returnError(w, fmt.Errorf("禁止IP失败: %s", err))
		return
	}

	// success
	returnOK(w, "成功BanIP", ips)
}
