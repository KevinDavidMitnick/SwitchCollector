package g

import (
	"encoding/json"
	"github.com/toolkits/file"
	"log"
	"sync"
)

type UdpConfig struct {
	Addr string `json:"addr"`
}

type HttpConfig struct {
	Addr string `json:"addr"`
}

type SwitchConfig struct {
	Ip         string `json:"ip"`
	Community  string `json:"community"`
	InFlowOid  string `json:"inFlowOid"`
	OutFlowOid string `json:"outFlowOid"`
	Timeout    int    `json:"timeout"`
}

type GlobalConfig struct {
	Udp      *UdpConfig    `json:"udp"`
	Http     *HttpConfig   `json:"http"`
	Expire   int           `json:"expire"`
	Interval int           `json:"interval"`
	Switch   *SwitchConfig `json:"switch"`
}

var (
	ConfigFile string
	config     *GlobalConfig
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
