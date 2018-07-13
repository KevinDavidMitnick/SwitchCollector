package funcs

import (
	"bytes"
	"fmt"
	"github.com/toolkits/sys"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

//Ping check,switch ping .
func Ping(ip string, timeout int) (int64, int64) {
	cmd := exec.Command("ping " + ip + " -c 2")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		return 0, 0
	}
	err, isTimeout := sys.CmdRunWithTimeout(cmd, time.Duration(timeout)*time.Millisecond)
	if isTimeout || err != nil {
		return 0, 0
	}

	data := stdout.Bytes()
	if len(data) == 0 {
		return 0, 0
	}

	fmt.Println(string(data))

	var rrt float64
	rrts := strings.Fields(string(data))
	for _, str := range rrts {
		if strings.Contains(str, "time") {
			temp := strings.Split(str, "=")[1]
			if value, err := strconv.ParseFloat(temp, 64); err == nil {
				rrt += value
			}
		}
	}

	return 0, int64(rrt / 2)
}
