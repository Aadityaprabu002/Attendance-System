var subbtn = document.getElementById("submit-btn")
subbtn.addEventListener('click',completeRegistration)

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
    xhr.open("POST","/teacher/signin/complete_registration");
    xhr.setRequestHeader("content-type","application/json")
    xhr.send(obj);
}
function completeRegistration(){
   
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
    });

    console.log(obj);
    sendRequest(obj);
}




