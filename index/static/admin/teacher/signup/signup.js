
// form and ajax

function submitForm(){
    let fname = document.querySelector("#fname").value;
    let lname = document.querySelector("#lname").value;
    let teacherid = document.querySelector("#teacherid").value
    let p = document.querySelectorAll("#password"); 
    let password = [];
    for(let i=0;i<p.length;i++){
        password[i] = p[i].value;
    }
    let email = document.querySelector("#email").value; 
    var re = /^[A-Za-z]+$/;
    if(fname.length == 0 && fname.length <= 2 || fname.length > 20){
        document.getElementById('response').innerHTML = "First name can't be empty and should be of length between 2 to 20";
        return;
    }
    if(!re.test(fname)){
        document.getElementById('response').innerHTML = "First name can only have characters";
        return;
    } 

    if(lname.length!=0){
        if(lname.length <= 2){
            document.getElementById('response').innerHTML = "Last name should be of length between 2 to 20 if there";
            return;
        }
        if(!re.test(lname)){
            document.getElementById('response').innerHTML = "Last name can only have characters";
            return
        } 
    }
    if (teacherid.length !=10){
        document.getElementById('response').innerHTML = "Teacher id should be having 10 digits";
        return;
    }
    for(let i=0;i<p.length;i++){
        password[i] = String(p[i].value);
    }
    var re = /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z])(?=.*[^a-zA-Z0-9])(?!.*\s).{8,15}$/;

    if (password[0]!=password[1]){
        document.getElementById('response').innerHTML = "Password doesnt match!";
        return;
    }
    if(!password[0].match(re)){
        document.getElementById('response').innerHTML = "Password must be between 8 to 15 characters which contain at least one lowercase letter, one uppercase letter, one numeric digit, and one special character";
        return;
    }
    var re = /\S+@\S+\.\S+/;
    if(!re.test(email)){
        document.getElementById('response').innerHTML = 'Enter valid email!';
        return;
    }

    let obj = JSON.stringify({
       "firstname": fname,
       "lastname": lname,
       "password": password,
       "email": email,
       "teacherid":teacherid
    })

    console.log(obj);
    sendRequest(obj);
}

function sendRequest(obj){
    xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(){
        if(this.readyState == 4 && this.status == 200){
            var result = JSON.parse(this.responseText);
            switch(result.Status){
                case 0 : document.querySelector("#response").innerHTML =  result.Response; break;
                case 1 : window.location.href = "/teacher/signin"; break;
            }
        }
    }
    xhr.open("POST","/admin/teacher/signup");
    xhr.setRequestHeader("content-type","application/json")
    xhr.send(obj);
}