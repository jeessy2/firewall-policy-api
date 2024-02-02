package ban

import (
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
)

// Ipv4Reg IPv4正则
var ipv4Reg = regexp.MustCompile(`((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])`)

// Ipv6Reg IPv6正则
var ipv6Reg = regexp.MustCompile(`((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))`)

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
		// 使用正则查找IPv4/IPv6
		result := ipv4Reg.FindString(ips[i])
		if result == "" {
			result = ipv6Reg.FindString(ips[i])
		}
		ip := net.ParseIP(result)
		if ip == nil {
			err = fmt.Errorf(errStr+"IP地址: %s 非法, 将不会被Ban! ", ips[i])
			continue
		}
		if ip.IsPrivate() {
			err = fmt.Errorf(errStr+"IP地址: %s 不能为私网地址, 将不会被Ban! ", ip.String())
			continue
		}
		if ip.IsLoopback() {
			err = fmt.Errorf(errStr+"IP地址: %s 不能为Loopback, 将不会被Ban! ", ip.String())
			continue
		}
		// 如果在白名单中, 将不会被Ban
		var find = false
		for j := 0; j < len(whiteList); j++ {
			if whiteList[j].Contains(ip) {
				find = true
				err = fmt.Errorf(errStr+"IP地址: %s 在白名单中, 将不会被Ban! ", ip.String())
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
