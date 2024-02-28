package ban

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jeessy2/firewall-policy-api/entity"
)

// returnError 返回错误信息
func returnError(w http.ResponseWriter, err error) {
	result := &entity.Result{}

	w.WriteHeader(http.StatusOK)
	result.Code = http.StatusOK
	result.Msg = fmt.Sprintf("Ban IP失败: %s", err)

	jsonData, _ := json.Marshal(result)
	log.Println(string(jsonData))

	json.NewEncoder(w).Encode(result)
}

// returnOK	返回成功信息
func returnOK(w http.ResponseWriter, msg string, data interface{}) {
	result := &entity.Result{}

	w.WriteHeader(http.StatusOK)
	result.Code = http.StatusOK
	result.Msg = msg
	result.Data = data

	jsonData, _ := json.Marshal(result)
	log.Println(string(jsonData))

	json.NewEncoder(w).Encode(result)
}
