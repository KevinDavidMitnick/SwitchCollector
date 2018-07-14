package funcs

import (
	"bytes"
	"github.com/toolkits/sys"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

//Ping check,switch ping .
func Ping(ip string, timeout int) (int64, int64) {
	cmd := exec.Command("/bin/bash", "-c", "ping "+ip+" -c 2")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		return 0, 0
	}
	err, isTimeout := sys.CmdRunWithTimeout(cmd, time.Duration(timeout)*time.Second)
	if isTimeout || err != nil {
		return 0, 0
	}

	data := stdout.Bytes()
	if len(data) == 0 {
		return 0, 0
	}

	var rrt float64
	rrts := strings.Fields(string(data))
	for _, str := range rrts {
		if strings.Contains(str, "time") {
			temp := strings.Split(str, "=")
			if len(temp) == 2 {
				if value, err := strconv.ParseFloat(strings.TrimSpace(temp[1]), 64); err == nil {
					rrt += value
				}
			}
		}
	}

	return 1, int64(rrt / 2)
}
