var colors = ['#15bf95','#ffbd49','#d2685c'];

function getColors(number) {
	if (number < 60){
		return colors[0];
	} else if ( (number>=60) &&(number<=90) ){
		return colors[1];
	} else {
		return colors[2];
	}
}

function FormatSpeed(number) {
	if(SWITCH_FLOW_UNIT == "Kbps"){
		return (number*8/1000).toFixed(2);
	} else if (SWITCH_FLOW_UNIT == "Mbps") {
		return (number*8/1000/1000).toFixed(2);
	} else if (SWITCH_FLOW_UNIT == "KB/s") {
		return (number/1024).toFixed(2);
	} else {
		return (number/1024/1024).toFixed(2);
	}
}

function timestampToTime(timestamp) {
    var date = new Date(timestamp * 1000);		//时间戳为10位需*1000，时间戳为13位的话不需乘1000
    Y = date.getFullYear() + '-';
    M = (date.getMonth()+1 < 10 ? '0'+(date.getMonth()+1) : date.getMonth()+1) + '-';
    D = date.getDate() + ' ';
    h = date.getHours() + ':';
    m = (date.getMinutes()+1 < 10 ? '0'+(date.getMinutes()+1) : date.getMinutes()+1) + ':';
    s = date.getSeconds();
    ss = date.getMilliseconds();
    return h+m+s;
}