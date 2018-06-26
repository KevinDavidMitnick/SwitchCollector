#### GetVisitLog

列出时间间隔内IP访问次数

#### 请求语法

```
GET /VisitLog HTTP/1.1
```

#### 请求元素(Request Elements)

|名称|类型|	描述|是否可选|
| ------------- |:-------------:|:-------------| ------------- |
|Period|Integer|获取多长时间间隔（以当前时间往前推）内的数据，最大取值30，单位：分钟|必填参数|


#### 响应元素(Response Elements)
|名称|类型|描述|是否非空|
| ------------- |:-------------: |:-------------:| ------------- |
|Data|Array|IP访问次数数据，集合，其中的元素为VisitLog|非空|
|StatisticsTime|Integer|统计时间点，单位：毫秒|非空|

#### VisitLog
|名称|类型|描述|是否非空|
| ------------- |:-------------: |:-------------:| ------------- |
|IP|String|IP地址|非空|
|VisitCount|Integer|访问次数|非空|

#### 请求示例

##### Request
```
GET /VisitLog?Period=5 HTTP/1.1
```

##### Response

```
	{
        "Data":[
		   {
    		    "IP":"218.90.171.226",
    		    "VisitCount":100
		   },
		   {
    		    "IP":"60.174.248.59",
    		    "VisitCount":66
		   }
        ],
        "StatisticsTime":18232893232
	}
```
