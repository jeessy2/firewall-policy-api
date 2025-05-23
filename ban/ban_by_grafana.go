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

	// 临时封禁
	temporaryBan := r.URL.Query().Get("temporaryBan") == "true"

	// 解析IP
	ipBan := make([]string, 0)
	ipUnBan := make([]string, 0)
	for i := 0; i < len(ga.Alerts); i++ {
		alert := ga.Alerts[i]
		sp := strings.Split(alert.Annotations.Description, "\n")
		ips, err := checkIP(sp)
		if err != nil {
			// 只打印
			log.Println(err)
		}
		if alert.Status == "firing" {
			ipBan = append(ipBan, ips...)
		} else if temporaryBan {
			// 如果是临时封禁&&不是firing状态, 则解封
			ipUnBan = append(ipUnBan, ips...)
		}
	}

	// 判断是否为空
	if len(ipBan) == 0 && len(ipUnBan) == 0 {
		returnError(w, fmt.Errorf("IP为空, Request: %s", ga))
		return
	}

	if len(ipBan) > 0 {
		log.Printf("IP %s 将被封禁\n", ipBan)
		err = banIP(ipBan, true)
	}

	if len(ipUnBan) > 0 {
		log.Printf("IP %s 将被解封\n", ipUnBan)
		err = banIP(ipUnBan, false)
	}

	if err != nil {
		returnError(w, fmt.Errorf("操作失败: %s", err))
		return
	}

	// success
	returnOK(w, "成功", map[string][]string{"ipBan": ipBan, "ipUnBan": ipUnBan})
}
