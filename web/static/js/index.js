//定义全局变量
var isscroll_port = true; //是否滚动端口列表
var portlists = [];    //端口列表
var selected_Device = null;    //选中的交换机
var cpu = echarts.init(document.getElementById('cpu'));
var store = echarts.init(document.getElementById('store'));
var namelist = [];          //当前交换机所有的up状态的端口名称
var datalist = [];          //折线图列表
var option_axis = {
    backgroundColor: 'rgba(66,93,117,0)',
    tooltip: {
        trigger: 'axis',
    },
    legend: {
        data:['端口发送速率('+ SWITCH_FLOW_UNIT+')','端口接收速率('+ SWITCH_FLOW_UNIT+')'],
    },
    xAxis: {
        type: 'category',
        splitLine: {
            show: false
        },
        nameTextStyle:{
            color: '#9abed4',
        },
        axisLine: {
            lineStyle: {
                type: 'solid',
                color: '#9abed4',
                width:'1'
            }
        },
        axisLabel: {
            rotate:45,
            textStyle: {
                color: '#9abed4',
            }
        },
    },
    yAxis: {
        type: 'value',
        splitLine: {
            show: false
        },
        axisLine: {
            lineStyle: {
                type: 'solid',
                color: '#9abed4',
                width:'1'
            }
        },
        axisLabel: {
            textStyle: {
                color: '#9abed4',
            }
        }
    },
    series: [{
        name: '端口发送速率('+ SWITCH_FLOW_UNIT+')',
        type: 'line',
        showSymbol: false,
        hoverAnimation: false,
        data: [],
        itemStyle:{
            normal:{
                color: '#d2685c',
                lineStyle:{
                    color : '#d2685c'
                }
            }
        },
    },{
        name: '端口接收速率('+ SWITCH_FLOW_UNIT+')',
        type: 'line',
        showSymbol: false,
        hoverAnimation: false,
        data: [],
        itemStyle:{
            normal:{
                color: '#00f2f4',
                lineStyle:{
                    color:'#00f2f4'
                }
            }
        },
    }]
};

//初始化页面
function init() {
    isscroll_port = true;
    portlists = [];
    namelist = [];
    selected_Device = null;
    $("#portlist").empty();
    $("#devicelists").empty();
    $("#portinfo").empty();
    $("#axislist").empty();
    $("div").remove(".system-panel");
    $("#axislists").css("width",720*CHART_LIST_NUM+"px");
    $(".right-chart-list").css("width",720*CHART_LIST_NUM+"px");
    $(".data-container").css("width",720*CHART_LIST_NUM+380+"px");
    getDeviceList();
}

//获取交换机列表
function getDeviceList() {
    fetch(API_SERVER+'/GetDeviceList', {
        method: 'GET',
        mode: 'cors',
    }).then(function(response) {
        return response.json();
    }).then((data) => {
        if(data.Data.length>0){
            selected_Device = data.Data[0];
            $("#selectDevice").text(selected_Device);
            for (var i = 0; i < data.Data.length; i++) {
                $("#devicelists").append("<li id='device-"+i+"' onclick='portList_Init(\""+data.Data[i]+"\")'>"+data.Data[i]+"</li>");
            }
            portList_Init(selected_Device);
        }
    });
}

//展开交换机列表
function showDevices() {
    $("#device-status").removeClass("pulldown");
    $("#device-status").addClass("pullup");
    $("#devicelists").fadeIn(100);
}


//端口列表注入
function portList_Init(ip){
    selected_Device = ip;
    $("#selectDevice").text(selected_Device);
    portlists = [];
    namelist = [];
    $("#portlist").empty();
    fetch(API_SERVER+'/GetInterfaceInfo?ip='+ip, {
        method: 'GET',
        mode: 'cors',
    }).then(function(response) {
        return response.json();
    }).then((data) => {
        for(var key in data.Data){
            var temp = data.Data[key];
            portlists.push({
                name: key,
                status: (temp["端口状态(1/2)"] == 1)? 'up': 'down',
                up: temp['端口发送速率(字节/秒)'],
                down: temp['端口接收速率(字节/秒)'],
                port: '',
            })
            if(temp["端口状态(1/2)"] == 1){
                namelist.push(key);
            }
        };
        var rowlen = (portlists.length%PORT_LIST_NUM == 0)?(parseInt(portlists.length/PORT_LIST_NUM)):(parseInt(portlists.length/PORT_LIST_NUM) + 1);
        for(var i = 1;i<= rowlen;i++){
            $("#portlist").append("<div class='row' id='row-"+i+"'></div>");
            for(j = (i-1)*PORT_LIST_NUM;j< i*PORT_LIST_NUM;j++){
                if(j >= portlists.length) break;  
                var item = portlists[j];
                $("#row-"+i).append("<li class='running-resource-item'>"+
                                "<i class='iconfont vmtp-photo'>"+
                                    "<img width='24' src='./static/images/port_"+ item.status +".png' alt=''>"+
                                "</i>"+
                                "<div class='running-resource-name'>"+ item.name +
                                    "<span class='running-resource-acount'>"+ item.port +"</span></div>"+
                                "<span class='number-status'><img width='14' src='./static/images/up.png' alt=''> "+ FormatSpeed(item.up)+ SWITCH_FLOW_UNIT  +"</span>"+
                                "<span class='number-status'><img width='14' src='./static/images/down.png' alt=''> "+ FormatSpeed(item.down)+ SWITCH_FLOW_UNIT  +"</span>"+
                            "</li>");
            }
        }
        portInfo(ip);
    });
};

//滚动端口列表
function AutoScroll_Ports(obj) {
    if(!isscroll_port) return;
    $(obj).find("#portlist:first").animate({
        marginTop: "-82px"
    },
    1000,
    function() {
        $(this).css({
            marginTop: "0px"
        }).find(".row:first").appendTo(this);
    });
}

//展开收起全部端口
function showAll_Port() {
    isscroll_port = !isscroll_port;
    if(isscroll_port){
        $("#portstatus").toggleClass("pulldown");
        $(".top-header").animate({height:"217px"});
        $(".line-tab1").show();
        $(".data-container").show();
    }else{
        var $portlists = $("#portlist .row");
        $portlists.sort(function(a,b){
            var valveNumOfa = $(a).attr("id").split("row-")[1];
            var valveNumOfb = $(b).attr("id").split("row-")[1];
            if(parseInt(valveNumOfa) < parseInt(valveNumOfb)) 
                return -1;
            else 
                return 1;
        });
        $portlists.detach().appendTo("#portlist");
        $("#portstatus").toggleClass("pulldown");
        $(".line-tab1").hide();
        $(".data-container").hide();
        $(".top-header").animate({height:"100%"});
    }
}

//获取详情
function portInfo(ip) {
    fetch(API_SERVER+'/GetDeviceInfo?ip='+ip, {
        method: 'GET',
        mode: 'cors',
    }).then(function(response) {
        return response.json();
    }).then((data) => {
        $("#portinfo").empty();
        $("div").remove(".system-panel");
        for(var i = 0;i < data.Data.length;i++){
            var temp = data.Data[i];
            for(var key in temp){
                if( key === "cpu使用率(%)" ){
                    var cpu_option = {
                        series: [{
                            type: 'liquidFill',
                            shape: 'circle',
                            radius: '80%',
                            data: [{
                                value: temp[key]/100,
                                itemStyle:{
                                    normal:{
                                        color: getColors(temp[key])
                                    }
                                }
                            }],
                            label: {
                                normal: {
                                    textStyle: {
                                        color: getColors(temp[key]),
                                        insideColor: getColors(temp[key]),
                                        fontSize: 25
                                    }
                                }
                            }
                        }]
                    };
                    cpu.setOption(cpu_option, true); 
                } else if( key === "内存使用率(%)" ){
                    var store_option = {
                        series: [{
                            type: 'liquidFill',
                            shape: 'circle',
                            radius: '80%',
                            data: [{
                                value: temp[key]/100,
                                itemStyle:{
                                    normal:{
                                        color: getColors(temp[key])
                                    }
                                }
                            }],
                            label: {
                                normal: {
                                    textStyle: {
                                        color: getColors(temp[key]),
                                        insideColor: getColors(temp[key]),
                                        fontSize: 25
                                    }
                                }
                            }
                        }]
                    };
                    store.setOption(store_option, true);
                } else if( key === "系统描述" ){
                    $(".left-message-list").append("<div class='usage-panel system-panel'>"+
                        "<div class='usage-panel-header title-type-2'>"+
                            "<h4 class='title'>系统描述</h4>"+
                        "</div>"+
                        "<div class='usage-panel-body usage-success'>"+
                            "<div class='usage-info'>"+
                                "<ul class='usage-info-list'>"+
                                    "<li class='usage-info-item usage-info-item-used info-item'>"+temp[key]+
                                    "</li>"+
                              "</ul>"+
                            "</div>"+
                        "</div>"+
                    "</div>");
                }else{
                    $("#portinfo").append("<div class='usage-info'>"+
                                "<ul class='usage-info-list'>"+
                                    "<li class='usage-info-item usage-info-item-used info-item'>"+
                                        "<span class='fl info-item-name'>"+ key +"</span>"+
                                        "<span class='fr info-item-mess'>"+ temp[key] +"</span>"+
                                    "</li>"+
                              "</ul>"+
                            "</div>");
                }
            }
        }
        getChartsList(ip);
    });
}

//获取单个折线图
function Loop_chart(ip,key,j,i){
    fetch(API_SERVER+'/GetInterfaceMetric?ip='+ip+'&filter='+key+'&accurate=true&period='+CHAR_TIME_STAMP, {
        method: 'GET',
        mode: 'cors',
    }).then(function(response) {
        return response.json();
    }).then((data) => {
        //初始化的时候只加载CHART_LIST_NUM*2张
        //页面转动的时候再次加载（刷新）CHART_LIST_NUM*2张图
        var temp = data.Data[key];
        var data1 = [], data2 = [];
        for(k = 0;k< temp["端口发送速率(字节/秒)"].length;k++){
            var dmp = temp["端口发送速率(字节/秒)"][k];
            data1.unshift({
                name: new Date(dmp.timestamp*1000),
                value: [
                    timestampToTime(dmp.timestamp),
                    FormatSpeed(dmp.value)
                ]
            })
        }
        for(k = 0;k< temp["端口接收速率(字节/秒)"].length;k++){
            var dmp = temp["端口接收速率(字节/秒)"][k];
            data2.unshift({
                name: new Date(dmp.timestamp*1000),
                value: [
                    timestampToTime(dmp.timestamp),
                    FormatSpeed(dmp.value)
                ]
            })
        }
        datalist.unshift({data1: data1,data2: data2,name: key});
        option_axis.series[0].data = data1;
        option_axis.series[1].data= data2;
        var axis = echarts.init(document.getElementById('axis'+key), 'dark');
        axis.setOption(option_axis, true);
        j++;
        if((j < CHART_LIST_NUM*2) && (j < namelist.length) && (j< i*2)){
            Loop_chart(ip,namelist[j],j,i);
        }
    })
}

//获取折线图列表
function getChartsList(ip){
    $("#axislist").empty();
    datalist = [];
    var rowlen = (namelist.length%2 == 0)?(namelist.length/2):(namelist.length/2 + 1);
    for(var i = 1;i<= rowlen;i++){
        $("#axislist").append("<div class='axiscontain' id='axis-li-"+i+"' class='axiscontain'></div>");
        var j = (i-1)*2;
        for(var k = (i-1)*2;k<i*2;k++){
            if(k >= namelist.length) return;
            $("#axis-li-"+i).append("<div class='axis-item' id='axis-chart-"+k+"'>"+
                "<div class='axis-title'>"+ namelist[k] +"</div>"+
                "<div id='axis"+ namelist[k] +"' style='width:700px; height: 450px;'></div>"+
            "</div>");
        }
        if((j < CHART_LIST_NUM*2) && (j < namelist.length)){
           Loop_chart(ip,namelist[j],j,i); 
        }
    }
}

//单个刷新折线图
function Update_chart(ip,ids,i){
    if(!isscroll_port) return;
    var index = ids[i];
    var key = namelist[index];
    fetch(API_SERVER+'/GetInterfaceMetric?ip='+ip+'&filter='+key+'&accurate=true&period='+CHAR_TIME_STAMP, {
        method: 'GET',
        mode: 'cors',
    }).then(function(response) {
        return response.json();
    }).then((data) => {
        //页面转动的时候再次加载（刷新）CHART_LIST_NUM*2张图
        var temp = data.Data[key];
        var data1 = [], data2 = [];
        for(k = 0;k< temp["端口发送速率(字节/秒)"].length;k++){
            var dmp = temp["端口发送速率(字节/秒)"][k];
            data1.unshift({
                name: new Date(dmp.timestamp*1000),
                value: [
                    timestampToTime(dmp.timestamp),
                    FormatSpeed(dmp.value)
                ]
            })
        }
        for(k = 0;k< temp["端口接收速率(字节/秒)"].length;k++){
            var dmp = temp["端口接收速率(字节/秒)"][k];
            data2.unshift({
                name: new Date(dmp.timestamp*1000),
                value: [
                    timestampToTime(dmp.timestamp),
                    FormatSpeed(dmp.value)
                ]
            })
        }
        datalist.unshift({data1: data1,data2: data2,name: key});
        option_axis.series[0].data = data1;
        option_axis.series[1].data= data2;
        var myChart = echarts.init(document.getElementById('axis'+key), 'dark');
        myChart.setOption(option_axis, true);
        i++;
        if(i<ids.length){
            Update_chart(ip,ids,i);
        }
    })
}


//滚动折线图列表
function AutoScroll_Charts(obj) {
    if(!isscroll_port || (portlists.length == 0) || (namelist.length == 0)) return;
    $(obj).find("#axislist:first").animate({
        marginLeft: "-720px"
    },
    1000,
    function() {
        $(this).css({
            marginLeft: "0px"
        }).find(".axiscontain:first").appendTo(this);
        //根据dom的id截取来判断是namelist里的第几个元素，从而懒加载刷新当前的图表
        var index = $(this).find(".axiscontain:first")[0].id.split("axis-li-")[1];
        var ids = [];
        if($(this).find(".axiscontain:first").find(".axis-item")[0])
            ids.push($(this).find(".axiscontain:first").find(".axis-item")[0].id.split("axis-chart-")[1]);
        if($(this).find(".axiscontain:first").find(".axis-item")[1])
            ids.push($(this).find(".axiscontain:first").find(".axis-item")[1].id.split("axis-chart-")[1]);
        for(var i = 1;i<CHART_LIST_NUM;i++){
            var f_index =  parseInt(index) + parseInt(i-1);
            var l_index = parseInt(index) + parseInt(i+1);
            var temp = $("#axis-li-"+f_index).nextUntil("#axis-li-"+l_index);
            if(temp.children(".axis-item")[0])
                ids.push(temp.find(".axis-item")[0].id.split("axis-chart-")[1]);
            if(temp.find(".axis-item")[1])
                ids.push(temp.find(".axis-item")[1].id.split("axis-chart-")[1]);
        }
        Update_chart(selected_Device,ids,0);
    });
}

$(document).ready(function() {
    init();
    setInterval('AutoScroll_Ports("#portlists")', PORT_SCROLL_TIME);
    setInterval('AutoScroll_Charts("#axislists")', CHART_SCROLL_TIME);
    $(document).click(function(event){
        var event = event || window.event;
        if (event && event.stopPropagation) {
                event.stopPropagation();
                if((event.target.className != 'cd') && (event.target != document.getElementById('device-status'))){
                    $("#device-status").addClass("pulldown");
                    $("#device-status").removeClass("pullup");
                    $("#devicelists").fadeOut(100);
                }
        } else {
                event.cancelBubble = true;
        }
    });
});