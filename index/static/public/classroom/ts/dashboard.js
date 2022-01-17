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
    xhr.open("POST","/teacher/dashboard");
    xhr.setRequestHeader("content-type","application/json")
    xhr.send(obj);
}

function submitForm(){

    let deptid = document.querySelector("#departmentid").value
    let courseid = document.querySelector("#courseid").value; 

    let obj = JSON.stringify({
       "departmentid": deptid,
       "courseid": courseid
    });
    console.log(obj);
    sendRequest(obj);
}
var btn = document.getElementById("submitBtn");
btn.addEventListener("click",submitForm)


