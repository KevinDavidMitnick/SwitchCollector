**交换机临时项目支持**

* 实现websocket服务，用于实时通信
* 实现http server,用于获取历史1分钟的数据，目标为1小时
* 实现udp server，用于交换机log推送
* 实现交换机端口流量获取代码，缓存在内存中，定时刷新到磁盘
* 实现交换机日志解析，缓存ip列表到缓存中，定时刷新到磁盘


1. 防火墙的日志记录，src地址获取 + 最后一次更新的时间戳 + 统计次数，　所有src总统计次数, 缓存一段时间间隔。定时删除超过最大超时时间的数据.
map[ip]AccessIP

2. 流量数据直接存储在rrd中,根据起始和结束时间，返回一段时间的历史数据.

   * 定时获取入口流量和出口流量，分别写入各自的rrd文件中,根据intervald定时器触发获取的步骤
   * 查询流量数据时，同时分别获取出口流量和入口流量，并且整合成一个对象，返回给前台。
   
3.　同时支持流量存储在rrd和内存中，两种方式.
   
## 交换机（网络设备）设计##
1. 模板内容
   * cpu使用率   (默认%)
   * mem使用率   (默认%)
   * 上行端口速率（默认单位kbps）
   * 下行端口速率 (默认单位kbps)
   * 指定端口速率 (默认单位kbps)
   * 指定端口状态 (默认0/1)
   * ping 可达检测（默人0/1）
   * 指定端口错误包数(int)
   * 系统描述 (字符串)
   * ping 延时
  
2. 扩展模板内容

   
3. 主配置文件(config.json)

```
{
    "udp" : {
        "enable": true,
        "addr" : "0.0.0.0:514"
    },
    "http" : {
        "enable": true,
        "addr" : "0.0.0.0:8080"
    },
    "trap" : {
        "enable": true,
        "addr" : "0.0.0.0:1622"
    },
    "template" : {
        "dir" : "/opt/switch/templates"
    },
    "switch" : {
        "dir": "/opt/switch/devices"
    }
    "expire" : 3600,
    "interval": 10
}

```

## template 目录结构应当如下: ##
  templates
  ----h3c
  --------h3c-v1.json
  --------h3c-v2.json
  ----huawei
  --------huawei-v1.json
  ----ciso
  --------ciso-v1.json
  ----other
  --------other1.json

说明: 程序启动的时候，应该加载template配置目录下的所有交换机类型和对应的模板。


## template 文件结构应该如下: ##
```
{
    "class" : "h3c",
    "type" :  "h3c-v1",
    "metric" : {
        "cpu使用率(%)"  : {
            "oid" : ".1.1.1.1..1.1.1.1.1",
            "interval" : 10,
        },
        "内存使用率"(%)" : {
            "oid" : ".1.1.1.1.1.1.1.1.2",
            interval: 10,
        },
        "ping检测(0/1)" : {
            "oid" : ".1.1.1.1.1.1.1.1.1.15",
            "interval" : 10,
        },
        "ping延时(ms)"  : {
            "oid" : ".1.1.1.1.1.1.1.1.1.16",
            "interval" : 10,
        }
    },
    "info" : {
        "系统名称" : {
            "oid" : ".1.1.11.1..11.1.13",
            "interval" : 10
        },
        "系统描述" : {
            "oid" : ".1.1.1.1.1.1.11.17",
            "interval" : 10
        }
    },
    "multimetric":{
        "interface" : {
            "上行端口速率(kbps)" : {
                "oid" : ".1.1.1.1.1.1.1.13",
                "interval" : 10
            }
            "下行端口速率(kbps)" : {
                "oid" : ".1.1.1.1.1.1.1.14",
                "interval" : 10
            }
        }
    },
    "multiinfo":{
        "interface" : {
            "端口状态(0/1)" : {
                "oid" : ".1.1.1.1.1.1.1.13",
                "interval" : 10
            }
            "端口描述" : {
                "oid" : ".1.1.1.1.1.1.1.13",
                "interval" : 10
            }
        }
    }
   "timeout": 1000,
   "interval": 10
}
```



## 网络设备配置文件结构如下 ##

```
{
   "ip" : "10.0.20.254",
   "community": "huayun.2017",
   "template" : {
       "class" : "h3c",
       "type"  : "h3c-v1"
   }
   "extension" : {
       "enable" : true,
       "metric" : {
           "风扇状态" : {
               "oid" : ".1.1.1.1.1.1.1.13",
               "interval" : 10
           }
       },
       "info" : {
           "网络接口数" : {
               "oid" : ".1.1.1.15.1.2.1.1.4",
               "interval" : 10
           }
       },
       "multimetric":{
            "interface" : {
            }
        },
        "multiinfo": {
            "interface" : {
                "错误丢包数" : {
                    "oid" : ".1.1.1.1.1.1.1.14",
                    "interval" : 10
                }
            }
        },
        "multiinfo": {
            "interface" : {
                "端口速率" : {
                    "oid" : ".1.1.1.1.1.1.1.15",
                    "interval" : 10
                }
            }
        }
    }
}
```


## 设计流程 ##
1. 加载配置文件,获取主配置信息。
2. 加载模板文件，获取模板对应信息。
3. 加载网络设备文件，获取网络设备对应的配置信息。
4. 根据配置信息,采集相应的指标,存储到内存中
5. 启动相应的http/udp/websocket服务，拉取信息
6. 定时将内存中的数据刷到backendz中。

