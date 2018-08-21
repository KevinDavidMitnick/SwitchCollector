package g

import (
	log "github.com/sirupsen/logrus"
	"os"
	"runtime/pprof"
	"sync"
	"time"
)

var wg sync.WaitGroup

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
	wg.Done()
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
	log.Println("Heap Profile stopped")
	wg.Done()
}

func DebugReport() {
	for {
		wg.Add(2)
		go cpuProfile()
		go heapProfile()
		wg.Wait()
		time.Sleep(300 * time.Second)
	}
}
