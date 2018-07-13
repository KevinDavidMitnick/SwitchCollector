package funcs

import (
	"fmt"
	"github.com/gaochao1/sw"
)

func Ping(ip string, timeout int, fastPingMode bool) (int64, int64) {
	rtt, err := sw.PingRtt(ip, timeout, fastPingMode)
	if err != nil {
		fmt.Println(err.Error())
		return 0, 0
	}
	return 1, int64(rtt)
}
