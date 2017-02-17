$(document).ready(function(){
	for (var i = $('div').length - 1; i >= 0; i--) {
		var div = $($('div')[i]);
		var span = div.children('span').length ? div.children('span')[0] : 0;
		var divHeigth = div[0].clientHeight;
		var spanHeigth = span.clientHeight;
		if( divHeigth < spanHeigth ) {
			div[0].style.height = spanHeigth + 'px';
		}
		console.log(div[0].clientHeight, spanHeigth);
	}
});