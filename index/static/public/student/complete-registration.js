var subbtn = document.getElementById("submit-btn")
subbtn.addEventListener('click',completeRegistration)

var camdiv = document.getElementById("camera");
var camera = loadCamera(camdiv);
startVideo();

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
    xhr.open("POST","/student/signin/complete_registration");
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
    var re = /\S+@\S+\.\S+/;
    if(!re.test(email)){
        document.getElementById('response').innerHTML = 'Enter valid email!';
        return;
    }
    for(let i=0;i<p.length;i++){
        password[i] = String(p[i].value);
    }
    re = /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[^a-zA-Z0-9])(?!.*\s).{8,15}$/;

    if (password[0]!=password[1]){
        document.getElementById('response').innerHTML = "Password doesnt match!";
        return;
    }
    if(!password[0].match(re)){
        document.getElementById('response').innerHTML = "Password must be between 8 to 15 characters which contain at least one lowercase letter, one uppercase letter, one numeric digit, and one special character";
        return;
    }
   
    let obj = JSON.stringify({
        "email":email,
        "password":password,
        "image64":image64
    });
    console.log(obj);
    sendRequest(obj);
}




