const MONTH = 1000*60*60*24*31;
function sendRequest(obj){
    
    xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(){
        if(this.readyState == 4 && this.status == 200){
            var result = JSON.parse(this.responseText);
            switch(result.Status){
                case 0 : document.querySelector("#classroom-response").innerText =  result.Response; break;
                case 1 : 
                window.location.reload(); 
                break;
            }
        }
    }
    xhr.open("POST","/teacher/dashboard");
    xhr.setRequestHeader("content-type","application/json")
    xhr.send(obj);
}

function submitForm(){

    let deptid = document.querySelector("#department-id").value
    let courseid = document.querySelector("#course-id").value; 
    let from = new Date(document.querySelector("#from-date").value);
    let to = new Date(document.querySelector("#to-date").value);
    let today = new Date();
    if (from.getTime() > to.getTime()){
        document.getElementById("classroom-response").innerText = "From date cant be lesser than to date";
        return;
    }else  if(to.getTime() < today.getTime() ){
        document.getElementById("classroom-response").innerText = "To date is already expried";
        return;
    }else if((to.getTime() - from.getTime()) < 5*MONTH ){
        document.getElementById("classroom-response").innerText = "Clasroom life should be atleast 5 month";
        return;
    }
    let obj = JSON.stringify({
       "departmentid": deptid,
       "courseid": courseid,
       "from":new Date(from),
       "to" : new Date(to)
    });
    console.log(obj);
    sendRequest(obj);
}
var btn = document.getElementById("submitBtn");
btn.addEventListener("click",submitForm)


