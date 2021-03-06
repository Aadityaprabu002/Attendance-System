const MAX_SESSION_HOUR = 4;
const MIN_SESSION_MINUTE = 20
function sendRequest(obj){
    
    xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(){
        if(this.readyState == 4 && this.status == 200){
            var result = JSON.parse(this.responseText);
            switch(result.Status){
                case 0 : document.querySelector("#response").innerHTML =  result.Response; break;
                case 1 : window.location.reload(); break;
            }
        }
    }
    xhr.open("POST",window.location.href);
    xhr.setRequestHeader("content-type","application/json")
    xhr.send(obj);
}

function submitForm(){

    let today = new Date();
    let dt = today.getFullYear()+'-'+(parseInt(today.getMonth())+1)+'-'+today.getDate();
    let st = document.querySelector("#start_time").value ;
    let et = document.querySelector("#end_time").value ; 
  
    if(st != "" || et != "") {
        st+=":00";
        et+=":00";
    }else {
        document.querySelector("#response").innerHTML = "Invalid start time or end time!";
        return;
    }

    st = new Date(dt+" "+st);
    et = new Date(dt+" "+et);
    dt = new Date(dt);
    console.log(dt);
    console.log(st,et);
    if (st.getTime() > et.getTime()){
        document.querySelector("#response").innerHTML = "Start time greater than end time! ";
        return;
    }else if( et.getTime() - st.getTime() > MAX_SESSION_HOUR * 60 * 60 * 1000 ){
        document.querySelector("#response").innerHTML = "Session time greater than 4 hours! ";
        return;
    }else if( et.getTime() - st.getTime() < MIN_SESSION_MINUTE * 60 * 1000){
        document.querySelector("#response").innerHTML = "Session time less than 20 minutes! ";
        return;
    }
    console.log('Sending session timings....');

    let obj = JSON.stringify({
       "date":dt,
       "start_time": st,
       "end_time": et
    });

    console.log(obj);
    sendRequest(obj);
}

var btn = document.getElementById("submitBtn");
btn.addEventListener("click",submitForm)



function addStudent(){
     let regnumber = document.getElementById("add-student").value;
     if (regnumber.length != 10){
         document.getElementById("add-student-response").innerText = "Regnumber length should be 10";
         return;
     }
     let obj = JSON.stringify({
         "regnumber":regnumber,
         "code":1
     });

     xhr = new XMLHttpRequest();
     xhr.onreadystatechange = function(){
         if(this.readyState == 4 && this.status == 200){
             var result = JSON.parse(this.responseText);
             console.log(result.Response);
             switch(result.Status){
                 case 0 : document.querySelector("#add-student-response").innerHTML =  result.Response; break;
                 case 1 : 
                 document.querySelector("#add-student-response").innerHTML =  result.Response;
                 setTimeout(function(){
                     window.location.reload(); 
                 },5*1000)
                 break;
             }
         }
     }
     xhr.open("POST",window.location.href+'/handlestudents');
     xhr.setRequestHeader("content-type","application/json")
     xhr.send(obj);     
}

function removeStudent(){
    let regnumber = document.getElementById("remove-student").value;
    if (regnumber.length != 10){
        document.getElementById("remove-student-response").innerText = "Regnumber length should be 10";
        return;
    }
    let obj = JSON.stringify({
        "regnumber":regnumber,
        "code":2
    });

    xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(){
        if(this.readyState == 4 && this.status == 200){
            var result = JSON.parse(this.responseText);
            console.log(result.Response);
            switch(result.Status){
                case 0 : document.querySelector("#remove-student-response").innerHTML =  result.Response; break;
                case 1 : 
                    document.querySelector("#remove-student-response").innerHTML =  result.Response;
                    setTimeout(function(){
                        window.location.reload(); 
                    },5*1000)
                    break;
            }
        }
    }
    xhr.open("POST",window.location.href+'/handlestudents');
    xhr.setRequestHeader("content-type","application/json")
    xhr.send(obj);

    
}

var addStudentBtn = document.getElementById("add-student-btn");
addStudentBtn.addEventListener("click",addStudent);

var removeStudentBtn = document.getElementById("remove-student-btn");
removeStudentBtn.addEventListener("click",removeStudent);