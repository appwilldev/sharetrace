function set_stiframe(){
    var iframe = document.createElement('iframe');
    console.log(iframe);
    var stiframe_url="http://st.apptao.com/1/st/webbeacon?share_url=" + encodeURIComponent(window.location);
    iframe.src = stiframe_url;
    console.log(stiframe_url);
    document.body.appendChild(iframe);
}

set_stiframe()
