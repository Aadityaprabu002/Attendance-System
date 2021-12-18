
// form and ajax

function submitForm(){
    let fname = document.querySelector("#fname").value;
    let lname = document.querySelector("#lname").value;
    let teacherid = document.querySelector("#teacherid").value
    let deptid = document.querySelector("#departmentid").value;
    let courseid =   document.querySelector("#courseid").value;
    let p = document.querySelectorAll("#password"); 
    let password = [];
    for(let i=0;i<p.length;i++){
        password[i] = p[i].value;
    }
    let email = document.querySelector("#email").value; 

    let obj = JSON.stringify({
       "firstname": fname,
       "lastname": lname,
       "password": password,
       "departmentid":deptid,
       "courseid":courseid,
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