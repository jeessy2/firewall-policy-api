package ban

import (
	"fmt"
	"log"
	"os"

	ssh "github.com/shenbowei/switch-ssh-go"
)

// banIP
func banIP(ips []string) error {

	if len(ips) == 0 {
		return fmt.Errorf("IP不能为空")
	}

	cmds := make([]string, 0)
	cmds = append(cmds, "sys")
	cmds = append(cmds, "security-policy ip")
	cmds = append(cmds, "rule name deny-auto")
	for i := 0; i < len(ips); i++ {
		cmds = append(cmds, "source-ip-host "+ips[i])
	}
	cmds = append(cmds, "quit")
	cmds = append(cmds, "quit")
	return execH3CCmd(cmds...)
}

// 执行H3C防火墙指令
func execH3CCmd(cmds ...string) error {

	user := os.Getenv("SSH_USER")
	password := os.Getenv("SSH_PASSWORD")
	ipPort := os.Getenv("SSH_IP_PORT")

	if condition := user == "" || password == "" || ipPort == ""; condition {
		return fmt.Errorf("请设置SSH_USER, SSH_PASSWORD, SSH_IP_PORT环境变量")
	}

	result, err := ssh.RunCommandsWithBrand(user, password, ipPort, ssh.H3C, cmds...)
	log.Println("执行结果:", result)

	return err
}
