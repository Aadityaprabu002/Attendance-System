var subbtn = document.getElementById("submit-btn")
subbtn.addEventListener('click',completeRegistration)

var camdiv = document.getElementById("camera");
var camera = loadCamera(camdiv);


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
function completeRegistration(){
    if(!modelsLoaded){
        document.getElementById('response').innerHTML = 'Error Wait for a minute and try again';
        return;
    }
    if(!isPhotoTaken){
        document.getElementById('response').innerHTML = 'Take a photo first!';
        return;
    }
    if(!isFaceVisible){
        document.getElementById('response').innerHTML = 'Form can not be submitted since no face was detected!';
        return;
    }
    let image64 = document.getElementById("canvas").toDataURL('image/png');
    let email = document.querySelector("#email").value; 
    let p = document.querySelectorAll("#password"); 
    let password = [];
    for(let i=0;i<p.length;i++){
        password[i] = p[i].value;
    }
    let obj = {
        "email":email,
        "password":password,
        "image64":image64
    }
    sendRequest(obj);
}




