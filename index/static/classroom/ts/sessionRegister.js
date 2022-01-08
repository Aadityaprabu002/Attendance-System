function sendRequest(obj){
    
    xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(){
        if(this.readyState == 4 && this.status == 200){
            var result = JSON.parse(this.responseText);
            switch(result.Status){
                case 0 : document.querySelector("#response").innerHTML =  result.Response; break;
                case 1 : window.location.href = "/teacher/dashboard/sessionDetails/"; break;
            }
        }
    }
    xhr.open("POST","/teacher/dashboard/sessionRegister/");
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
    console.log(st);
    st.setTime(st.getTime() - st.getTimezoneOffset()*60*1000)
    et = new Date(dt+" "+et);
    et.setTime(et.getTime() - et.getTimezoneOffset()*60*1000);

    let obj = JSON.stringify({
       "start_time": st,
       "end_time": et
    });

    console.log(obj);
    sendRequest(obj);
}
var btn = document.getElementById("submitBtn");
btn.addEventListener("click",submitForm)


