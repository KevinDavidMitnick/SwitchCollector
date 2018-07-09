package g

import (
	"encoding/json"
	"github.com/toolkits/file"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type UdpConfig struct {
	Enabled bool   `json:"enabled"`
	Addr    string `json:"addr"`
}

type HttpConfig struct {
	Enabled bool   `json:"enabled"`
	Addr    string `json:"addr"`
}

type TemplatesConfig struct {
	Dir string `json:"dir"`
}

type NetDevicesConfig struct {
	Dir string `json:"dir"`
}

type GlobalConfig struct {
	Udp        *UdpConfig        `json:"udp"`
	Http       *HttpConfig       `json:"http"`
	Templates  *TemplatesConfig  `json:"templates"`
	NetDevices *NetDevicesConfig `json:"netdevices"`
	Expire     int               `json:"expire"`
	Interval   int               `json:"interval"`
}

type Metric struct {
	Oid      string `json:"oid"`
	Interval int    `json:"interval"`
}

type MetricTemplate struct {
	Class        string            `json:"class"`
	Type         string            `json:"type"`
	Metrics      map[string]Metric `json:"metrics"`
	Infos        map[string]Metric `json:"infos"`
	MultiMetrics map[string]Metric `json:"multimetrics"`
	MultiInfos   map[string]Metric `json:"multiinfos"`
	Timeout      int               `json:"timeout"`
	Interval     int               `json:"interval"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	metricM    map[string]*MetricTemplate
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	lock.Lock()
	defer lock.Unlock()

	config = &c

	log.Println("read config file:", cfg, "successfully")

}

func listDir(dirpath string) (files []string, err error) {
	dirs, err := ioutil.ReadDir(dirpath)
	if err != nil {
		return
	}

	for _, dir := range dirs {
		pathname := dirpath + string(os.PathSeparator) + dir.Name()
		if dir.IsDir() {
			if childFiles, err := listDir(pathname); err == nil {
				files = append(files, childFiles...)
			}
		} else {
			files = append(files, pathname)
		}
	}
	return files, nil
}

func LoadTemplatesConfig() {
	lock.Lock()
	defer lock.Unlock()
	dir := config.Templates.Dir

	files, e := listDir(dir)
	if e != nil || len(files) == 0 {
		log.Fatalln("template dir load err...")
	}

	metricM = make(map[string]*MetricTemplate)

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			log.Fatalln("open file:", file, "failed")
		}
		defer f.Close()
		var metricTemplate MetricTemplate
		err1 := json.NewDecoder(f).Decode(&metricTemplate)
		if err1 != nil {
			log.Fatalln("template file decode error:", file)
		}
		metricM[metricTemplate.Type] = &metricTemplate
	}
	log.Println("load template config dir:", dir, "successfully")
	log.Println(metricM)
}

func ParseNetDevicesConfig(filename string) {

}
