function set_stiframe(){
    var iframe = document.createElement('iframe');
    var stiframe_url="https://st.apptao.com/1/st/webbeacon?share_url=" + encodeURIComponent(window.location);
    iframe.src = stiframe_url;
    document.body.appendChild(iframe);
}

set_stiframe()


function gotoAppStore(buttonid){
    if (buttonid == "undefined"){
		buttonid = 1;
    }
    var iframe = document.createElement('iframe');
    var stiframe_url="https://st.apptao.com/1/st/webbeaconbutton?buttonid=" + buttonid ;
    iframe.src = stiframe_url;
    document.body.appendChild(iframe);
}
