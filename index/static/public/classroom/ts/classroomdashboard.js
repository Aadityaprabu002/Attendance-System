const MAX_SESSION_HOUR = 4;
const MIN_SESSION_MINUTE = 20
function sendRequest(obj){
    
    xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(){
        if(this.readyState == 4 && this.status == 200){
            var result = JSON.parse(this.responseText);
            switch(result.Status){
                case 0 : document.querySelector("#response").innerHTML =  result.Response; break;
                case 1 : window.location.reload(); break;
            }
        }
    }
    xhr.open("POST",window.location.href);
    xhr.setRequestHeader("content-type","application/json")
    xhr.send(obj);
}

function submitForm(){

    let today = new Date();
    let dt = today.getFullYear()+'-'+today.getMonth()+1+'-'+today.getDate();
    console.log(dt);
    let st = document.querySelector("#start_time").value
    let et = document.querySelector("#end_time").value; 
    
    st = new Date(dt+" "+st);
    et = new Date(dt+" "+et);
    dt = new Date(dt);

    if (st.getTime() > et.getTime()){
        document.querySelector("#response").innerHTML = "Start time greater than end time! javascript";
        return;
    }else if( et.getTime() - st.getTime() > MAX_SESSION_HOUR * 60 * 60 * 1000 ){
        document.querySelector("#response").innerHTML = "Session time greater than 4 hours! javascript";
        return;
    }else if( et.getTime() - st.getTime() < MIN_SESSION_MINUTE * 60 * 1000){
        document.querySelector("#response").innerHTML = "Session time less than 20 minutes! javascript";
        return;
    }

    // st.setTime(st.getTime() - st.getTimezoneOffset()*60*1000)
    // et.setTime(et.getTime() - et.getTimezoneOffset()*60*1000);

    let obj = JSON.stringify({
       "date":dt,
       "start_time": st,
       "end_time": et
    });

    console.log(obj);
    sendRequest(obj);
}

var btn = document.getElementById("submitBtn");
btn.addEventListener("click",submitForm)


