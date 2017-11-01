var debug = true;
var sock = null;
var origin = window.location.origin;
var recInterval = null;
var parameter={};
var message={parameter: parameter};
var selVideo="";
var background="IMAGES/JukeboxBackground.jpg";
var backgroundcolor="#ffffff"
var videoPages = 1;
var videoPage = 1;
var data;

$( function(){
  	$( ".Reload" ).button({
        icon: "ui-icon-refresh",
        showLabel: false
  	}).click(function(){
		location.reload(true);
  	});
	
  	$( "#VideoPrevButton" ).button({
		label: "zurück"
  	}).hide();
  	$( "#VideoNextButton" ).button({
		label: "weiter"
  	}).hide();
  	
	$( "#VideoNextButton" ).click(function() { VideoNextButton(); });
  	$( "#VideoPrevButton" ).click(function() { VideoPrevButton(); });
  	$( "#Dialog" ).click(function(){
		var infotext = "width: "+screen.width + " x height: "+screen.height;
		infotext += "<br>";	
		popupInfoText("Info",infotext);
  	});
  	$( "#VideoList" ).button().click(function(){
		getVideoList();
  	});
 	//$(".menuButton").button();
  	$( "#Video" ).on( "click", function( event) {
		if (debug) { console.log("### click on SelectedVideo start ###"); }
		setMenuView("videoButtons");
		document.getElementById("SelectedVideo").pause();
	} );
  	$( "#Video" ).on( "select", function( event) {
		if (debug) { console.log("### select SelectedVideo start ###"); }
		setMenuView("videoButtons");
		document.getElementById("SelectedVideo").pause();
	} );
  	$( "#Video" ).on( "ended", function( event) {
		if (debug) { console.log("### select SelectedVideo start ###"); }
		setMenuView("videoButtons");
	} );
  	$( "#Button-radio" ).on( "click", function( event, ui ) {
		setMenuView("rockabillyradio");
	  	 $("#rockabillyradio").html('<object width=100% height=600 data="http://rockabilly-radio.net/"></object>');
  	} );
  	$( ".zurueckButton" ).on( "click", function( event, ui ) {
		setMenuView("menuButtons");
  	} );
  	$( "#Button-selectVideo" ).on( "click", function( event, ui ) {
		if (videoPages>1){
			$( "#VideoNextButton" ).button().show();
		}
		setMenuView("videoButtons");
  	});
  	$( "#Button-randomVideo" ).on( "click", function( event, ui ) {
		setMenuView("randomVideoButtons");
  	});
  	$( "#Button-randomVideoBack" ).on( "click", function( event, ui ) {
		setMenuView("menuButtons");
  	});
} );
function setMenuView(MenuID){
	document.getElementById("randomVideoButtons").style.display = "none";
	document.getElementById("videoButtons").style.display = "none";
	document.getElementById("rockabillyradio").style.display = "none";
	document.getElementById("menuButtons").style.display = "none";
	document.getElementById(MenuID).style.display = "block";
	
}


function VideoPrevButton(){
	if (debug) { console.log("### VideoPrevButton start ###"); }
	$("#videoButtonPage-"+videoPage).hide();	
	videoPage--;
	$("#videoButtonPage-"+videoPage).show();	
	if (videoPage==1){
		$( "#VideoPrevButton" ).button().hide();
	}
	if (videoPage<videoPages){
		$( "#VideoNextButton" ).button().show();
	}
	
}

function VideoNextButton(){
	if (debug) { console.log("### VideoNextButton start ###"); }
	$("#videoButtonPage-"+videoPage).hide();	
	videoPage++;
	$("#videoButtonPage-"+videoPage).show();	
	if (videoPage==videoPages){
		$( "#VideoNextButton" ).button().hide();
	}
	if (videoPage>1){
		$( "#VideoPrevButton" ).button().show();
	}
}

function videoZurueck(event, ui){
	if (debug) { console.log("### videoZurueck start ###"); }
	document.getElementById("videoButtons").style.display = "none";
	document.getElementById("rockabillyradio").style.display = "none";
	document.getElementById("menuButtons").style.display = "block";
}
function videoSelect(event, ui){
	if (debug) { console.log("### videoSelect start ###"); }
	document.getElementById("videoButtons").style.display = "none";
	document.getElementById("rockabillyradio").style.display = "none";
	document.getElementById("menuButtons").style.display = "block";
	document.getElementById(event.target.id).style.display = "none";
	document.getElementById(event.target.id).pause();
}


var getPing = function(){
	if (debug) { console.log("### setPing start ###"); }
	message.Typ="getPing";
	message.parameter.dummy1="dummyValue1"; //document.getElementById("dummy1").value;
	message.parameter.dummy2="dummyValue2"; //document.getElementById("dummy2").value;
	sock.send(JSON.stringify(message));
}

var getMediaData = function(){
	if (debug) { console.log("### getMediaData start ###"); }
	message.Typ="getMediaData";
	sock.send(JSON.stringify(message));
}

function playVideo(event, ui){
	if (debug) { console.log("### playVideo start ###"); }
	console.log(event)
	var videoId=event.target.id.split("-")[1];
	$( "#Video" ).contents().remove();
	//var video='<video style="display: block" id="SelectedVideo" controls autoplay class="FSvideo"><source src="http://player:8080/videos/'+data.data.VideoData.Videos[videoId-1].FileName+'"></video>';
	var video='<video id="SelectedVideo" controls autoplay class="FSvideo"><source src="http://player:8080/videos/'+data.data.VideoData.Mediainfos[videoId].FileName+'"></video>';
	$( "#Video" ).append(video);
	document.getElementById("videoButtons").style.display = "none";
	document.getElementById("Video").style.display = "block";
}

function popupInfoText(title,infohtml){
	if (debug) { console.log("### popupInfoText start ###"); }

	var infoPopUp = $( "#InfoPopUp" );
	if(infohtml!=""){
		var InfoPopUp = $( "#InfoPopUp" );
		$( "#InfoPopUp" ).dialog({
			title: title,
			autoOpen: false,
			width: 800
		});		
		InfoPopUp.contents().remove();
		InfoPopUp.append(infohtml);
		InfoPopUp.dialog( "open" );
	}
}

var sockonopen = function(InspireServer) {
	if (debug) { console.log("### sockonopen start ###"); }
	getMediaData();
}

function setMediaData(){
	if (debug) { console.log("### setMediaData start ###"); }
	console.log("data");
	console.log(data);
	// set Video Buttons
	$( "#videoButtonPages" ).contents().remove();
	var Videos = data.data.VideoData.Mediainfos;
	console.log("Videos");
	console.log(Videos);
	
	var jbbclass ="";
	videoPages = Math.ceil(Videos.length/10);
	videoPage=0;
	for(var i in Videos)
	{
		jbbclass="jbb-"+((i%10)+1);
		if(((i%10)+1)==1){
			videoPage++;
			$('<div id="videoButtonPage-'+videoPage+'" ></div>').appendTo("#videoButtons");
		}				
		var button='<Button>'+Videos[i].InterpretName+'<br><b>'+Videos[i].SongName+'</b></Button>';
		$(button).appendTo("#videoButtonPage-"+videoPage)
			.addClass("videoButton")
			.addClass(jbbclass)
			.attr("id","Button-"+i).on( "click", function(event, ui) {
				playVideo(event, ui);	
			});
		$("#videoButtonPage-"+videoPage).hide();
	}
	videoPage=1;
	$("#videoButtonPage-1").show();
}

function new_conn(){
	if (debug) { console.log("### new_conn start ###"); }
	sock = new SockJS(origin + '/ws/serverdata', {
		debug: true,
		devel: true,
		protocols_whitelist: "['websocket', 'xdr-streaming', 'xhr-streaming', 'iframe-eventsource', 'iframe-htmlfile', 'xdr-polling', 'xhr-polling', 'iframe-xhr-polling', 'jsonp-polling']",
	});
	clearInterval(recInterval);
 	sock.onopen = function(msg){
       	sockonopen();
    };
    sock.onmessage = function(msg){
		msg.data = JSON.parse(msg.data)
		console.log(msg);
		switch(msg.data.XMLName.Local){
			case "setPing":
				setPing(msg);
				break;
			case "setMediaData":
				data=msg;
				setMediaData(); 
				break;
			default:
				popupInfoText("error","unbekannte Message vom Server: "+msg.data.XMLName.Local); 
				console.log("### unbekannte Message: "+msg.data.XMLName.Local+" erhalten ###");
		}
    };
    sock.onclose = function(){
		recInterval = window.setInterval(function () {
			console.log("try to reconnect...");
	        new_conn();
	    }, 2000);    
	};
};
new_conn();

