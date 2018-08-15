package g

import (
	log "github.com/sirupsen/logrus"
	"os"
	"runtime/pprof"
	"time"
)

// 生成 CPU 报告
func cpuProfile() {
	os.Remove("cpu.prof")
	f, err := os.OpenFile("cpu.prof", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	log.Println("CPU Profile started")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	time.Sleep(300 * time.Second)
	log.Println("CPU Profile stopped")
}

// 生成堆内存报告
func heapProfile() {
	os.Remove("heap.prof")
	f, err := os.OpenFile("heap.prof", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	time.Sleep(300 * time.Second)

	pprof.WriteHeapProfile(f)
	log.Println("Heap Profile generated")
}

func DebugReport() {
	for {
		cpuProfile()
		heapProfile()

		time.Sleep(time.Duration(Config().Expire) * time.Second)
	}
}
