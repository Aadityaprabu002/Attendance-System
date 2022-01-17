
function sendRequest(){
    xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(){
        if(this.readyState == 4 && this.status == 200){
            var result = JSON.parse(this.responseText);
            console.log(result.Reponse);
            // switch(result.Status){
            //     case 0 : document.querySelector("#response").innerHTML =  result.Response; break;
            //     case 1 : window.location.href = "/teacher/dashboard"; break;
            // }
        }
    }
    xhr.open("POST",window.location.href + "/postattendance");
    xhr.setRequestHeader("content-type","application/json")
    xhr.send(obj);
}

function postAttendance(){
  
    var student_list = document.querySelectorAll("#attendance");
    var d1 = {};
    for (student of student_list){
        if(d1[student.name]){
            d1[student.name].push(student.checked);
        }else{
            d1[student.name] = [student.checked];
        }
    }
    var d2={};
    var list = []
    for (student in d1){
        if(d1[student][0] == true && d1[student][1] == false){
            d2[student] = false;
        }else if( d1[student][0] == false && d1[student][1] == true){
            d2[student] = true;
        }else{
            document.getElementById("response").innerText = "Please Review all the students before submission!";
        }
        // continue from here
    }
    let obj = JSON.stringify(d2);
    console.log(obj);
    // sendRequest(d2);
    
}

var btn = document.getElementById("post-attendance-btn");

if( btn!=null && typeof(btn)!=undefined){
    btn.addEventListener("click",postAttendance)
}