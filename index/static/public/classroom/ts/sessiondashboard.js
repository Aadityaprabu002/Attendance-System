
function sendRequest(obj){
    xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(){
        if(this.readyState == 4 && this.status == 200){
            var result = JSON.parse(this.responseText);
            console.log(result.Response);
            switch(result.Status){
                case 0 : document.querySelector("#response").innerHTML =  result.Response; break;
                case 1 : 
                    document.querySelector("#response").innerHTML =  result.Response + "Refreshing in 5 secs.."; 
                    setTimeout(function(){
                        window.location.reload();   
                    },5000);
                    break;
            }
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
    
    var list = []
    for (student in d1){
        let d2 = {};
        if(d1[student][0] == true && d1[student][1] == false){
            d2["regnumber"] = student
            d2["is_present"] = false;

        }else if( d1[student][0] == false && d1[student][1] == true){
            d2["regnumber"] = student
            d2["is_present"] = true;
        }else{
            document.getElementById("response").innerText = "Please Review all the students before submission!";
            return;
        }
        list.push(d2);
    }
    let obj = JSON.stringify(list);
    console.log(obj);
    let submit = confirm("Do you want to continue ?");
    if(submit){
        sendRequest(obj);
    }
     
}

var btn = document.getElementById("post-attendance-btn");

if( btn!=null && typeof(btn)!=undefined){
    btn.addEventListener("click",postAttendance)
}