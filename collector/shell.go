package collector

import (
	"log"
	"os/exec"
	"strings"
)

// 执行 shell 命令返回输出
func ExecCommand(cmd string) (error, string) {
	var _cmd *exec.Cmd

	_cmd = exec.Command("/bin/bash", "-c", cmd)
	str, err := _cmd.Output()
	if err != nil {
		log.Println("exec.Command failed.", err)
		return err, ""
	}
	return nil, strings.Replace(string(str), "\n", "", 1)
}
