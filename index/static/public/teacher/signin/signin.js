function sendRequest(obj){
    
    xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(){
        if(this.readyState == 4 && this.status == 200){
            var result = JSON.parse(this.responseText);
            switch(result.Status){
                case 0 : document.querySelector("#response").innerHTML =  result.Response; break;
                case 1 : window.location.href = "/teacher/dashboard"; break;
            }
        }
    }
    xhr.open("POST","/teacher/signin");
    xhr.setRequestHeader("content-type","application/json")
    xhr.send(obj);
}

function submitForm(){

    let teacherid = document.querySelector("#teacherid").value
    let password = document.querySelector("#password").value; 
    if(teacherid.length != 10){
        document.getElementById("response").innerText = "Enter a valid teacher id";
        return;
    }
    let obj = JSON.stringify({
       "teacherid": teacherid,
       "password": password
    });
    console.log(obj);
    sendRequest(obj);
}

var btn = document.querySelector("#submitBtn");
btn.addEventListener("click",submitForm);




