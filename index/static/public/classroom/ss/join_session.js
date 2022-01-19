var btn = document.getElementById("submit");
btn.addEventListener("click",submitForm);

function sendRequest(obj){
    xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(){
        if(this.readyState == 4 && this.status == 200){
            var result  = JSON.parse(this.responseText);
            switch(result.Status){
                case 0 : document.querySelector("#response").innerHTML =  result.Response; break;
                case 1 : window.location.href = "/student/dashboard/session"; break;
            }
        }
    }
    xhr.open("POST",window.location.href );
    xhr.setRequestHeader("content-type","application/json")
    xhr.send(obj);
}

function submitForm(){
   
    let obj = JSON.stringify({
       "session_key": document.getElementById("session-key").value,
    })
    sendRequest(obj);
}

