//定义配置信息

var API_SERVER = 'http://10.21.1.225:8080';

var PORT_SCROLL_TIME = 10*1000;											//端口列表滚动时间（秒*1000）
var CHART_SCROLL_TIME = 15*1000;										//折线图滚动时间（秒*1000）
//var CHART_UPDATE_TIME = 15*1000;										//折线图刷新时间（秒*1000）

var PORT_LIST_NUM = 6;													//端口列表每行显示个数
var CHART_LIST_NUM = 3;													//折线图显示列数，默认一列2个折线图
var CHAR_TIME_STAMP = 60*60;											//折线图查询接口的时间period（秒）
var SWITCH_FLOW_UNIT = 'Mbps';											//交换机流量单位 Kbps, Mbps, KB/s, MB/s
