package ban

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// checkIP 检查IP地址是否合法
func checkIP(ips []string) (toBeBanIPs []string, err error) {
	whiteList := getIPWhiteList()

	var errStr string
	for i := 0; i < len(ips); i++ {
		if err != nil {
			errStr = err.Error()
		}

		if condition := ips[i] == ""; condition {
			err = fmt.Errorf(errStr + "IP地址不能为空!")
			log.Println(err)
			continue
		}
		ip := net.ParseIP(ips[i])
		if ip == nil {
			err = fmt.Errorf(errStr+"IP地址: %s 非法, 将不会被Ban! ", ips[i])
			continue
		}
		if ip.IsPrivate() {
			err = fmt.Errorf(errStr+"IP地址: %s 不能为私网地址, 将不会被Ban! ", ips[i])
			continue
		}
		if ip.IsLoopback() {
			err = fmt.Errorf(errStr+"IP地址: %s 不能为Loopback, 将不会被Ban! ", ips[i])
			continue
		}
		// 如果在白名单中, 将不会被Ban
		var find = false
		for j := 0; j < len(whiteList); j++ {
			if whiteList[j].Contains(ip) {
				find = true
				err = fmt.Errorf(errStr+"IP地址: %s 在白名单中, 将不会被Ban! ", ips[i])
				break
			}
		}
		if !find {
			toBeBanIPs = append(toBeBanIPs, ip.String())
		}
	}
	return
}

// getIPWhiteList 获取白名单
func getIPWhiteList() (ipWhiteList []net.IPNet) {
	white := os.Getenv("IP_WHITE_LIST")
	if white == "" {
		return
	}
	sp := strings.Split(white, ",")
	for i := 0; i < len(sp); i++ {
		_, ipNet, err := net.ParseCIDR(sp[i])
		if err != nil {
			ip := net.ParseIP(sp[i])
			if ip != nil {
				ipNet = &net.IPNet{
					IP:   ip,
					Mask: net.CIDRMask(32, 32),
				}
				ipWhiteList = append(ipWhiteList, *ipNet)
				continue
			} else {
				log.Println(err.Error())
				continue
			}
		}
		ipWhiteList = append(ipWhiteList, *ipNet)
	}
	return
}
