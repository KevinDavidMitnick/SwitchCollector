#### GetFlowQuantityBytes

列出时间间隔内的流量数据

#### 请求语法

```
GET /FlowQuantityBytes HTTP/1.1
```

#### 请求元素(Request Elements)

|名称|类型|	描述|是否可选|
| ------------- |:-------------:|:-------------| ------------- |
|Period|Integer|获取多长时间间隔（以当前时间往前推）内的数据，最大取值30，单位：分钟|必填参数|


#### 响应元素(Response Elements)
|名称|类型|描述|是否非空|
| ------------- |:-------------: |:-------------:| ------------- |
|Data|Array|流量列表数据，集合，其中的元素为FlowQuantity|非空|

#### FlowQuantity
|名称|类型|描述|是否非空|
| ------------- |:-------------: |:-------------:| ------------- |
|Time|Integer|统计时间点，单位：毫秒|非空|
|InFlowQuantity|Integer|入口流量数据，单位：KBps|非空|
|OutFlowQuantity|Integer|出口流量数据，单位：KBps|非空|

#### 请求示例

##### Request
```
GET /FlowQuantityBytes?Period=5 HTTP/1.1
```

##### Response

```
	{
        "Data":[
		   {
    		    "Time":18232893232323,
    		    "InFlowQuantity":174832,
    		    "OutFlowQuantity":1212432
		   },
		   {
    		    "Time":182328555553,
    		    "InFlowQuantity":195564,
    		    "OutFlowQuantity":4545244
		   }
        ]
	}
```
