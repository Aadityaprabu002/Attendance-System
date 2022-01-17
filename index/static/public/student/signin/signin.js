function sendRequest(obj){
    
    xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(){
        if(this.readyState == 4 && this.status == 200){
            var result = JSON.parse(this.responseText);
            switch(result.Status){
                case 0 : document.querySelector("#response").innerHTML =  result.Response; break;
                case 1 : window.location.href = "/student/dashboard"; break;
            }
        }
    }
    xhr.open("POST","/student/signin");
    xhr.setRequestHeader("content-type","application/json")
    xhr.send(obj);
}

function submitForm(){

    let regnumber = document.querySelector("#rollno").value
    let password = document.querySelector("#password").value; 
    let email = document.querySelector("#email").value; 
    let obj = JSON.stringify({
       "email": email,
       "regnumber":regnumber,
       "password": password
    });
    console.log(obj);
    sendRequest(obj);
}

var btn = document.querySelector("#submitBtn");
btn.addEventListener("click",submitForm);




