package ban

import (
	"fmt"
	"net/http"
)

// BanByIP 通过IP地址Ban IP
func BanByIP(w http.ResponseWriter, r *http.Request) {

	ip := r.FormValue("ip")
	if ip == "" {
		returnError(w, fmt.Errorf("IP不能为空"))
		return
	}

	ips, err := checkIP([]string{ip})
	if err != nil {
		returnError(w, err)
		return
	}

	err = banIP(ips)
	if err != nil {
		returnError(w, err)
		return
	}

	// success
	returnOK(w, "成功BanIP", ips)
}
