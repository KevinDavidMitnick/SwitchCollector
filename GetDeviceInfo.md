#### GetDeviceInfo

列出所有网络设备列表
#### 请求语法

```
GET /GetDeviceInfo HTTP/1.1
```

#### 响应元素(Response Elements)
|名称|类型|描述|是否非空|
| ------------- |:-------------: |:-------------:| ------------- |
|Data|Array|返回网络设备名称，IP地址         |非空|
|StatisticsTime|Integer|统计时间点，单位：秒 |非空|

#### DeviceInfo
|名称|类型|描述|是否非空|
| ------------- |:-------------: |:-------------:| ------------- |
|ip |String|IP地址|非空|

#### 请求示例

##### Request
```
GET /GetDeviceInfo?ip=10.0.20.254 HTTP/1.1
```

##### Response

```
{
    "Data":[
        {
            "系统描述":"H3C Comware Platform Software, Software Version 7.1.045, Release 1121
H3C S5560-54C-EI
Copyright (c) 2004-2016 Hangzhou H3C Tech. Co., Ltd. All rights reserved."
        },
        {
            "网络接口数":124
        },
        {
            "ping检测(0/1)":1
        },
        {
            "风扇口状态":1
        },
        {
            "系统名称":"H3C-5560"
        },
        {
            "ping延时(ms)":1
        }
    ],
    "StatisticsTime":1531730823
}
```
