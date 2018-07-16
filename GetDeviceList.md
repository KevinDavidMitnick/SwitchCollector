#### GetDeviceList

列出所有网络设备列表
#### 请求语法

```
GET /GetDeviceList HTTP/1.1
```

#### 响应元素(Response Elements)
|名称|类型|描述|是否非空|
| ------------- |:-------------: |:-------------:| ------------- |
|Data|Array|返回网络设备名称，IP地址         |非空|
|StatisticsTime|Integer|统计时间点，单位：秒 |非空|

#### DeviceList
|名称|类型|描述|是否非空|
| ------------- |:-------------: |:-------------:| ------------- |
|IP|String|IP地址|非空|

#### 请求示例

##### Request
```
GET /GetDeviceList HTTP/1.1
```

##### Response

```
	{
        "Data":["218.90.171.226", "IP":"60.174.248.59"],
        "StatisticsTime":18232893232
	}
```
