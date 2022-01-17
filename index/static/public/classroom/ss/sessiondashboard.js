
var timer;
var camera;
var isAttBtnLoaded = false;
var isAttPosted = false;

const MINUTE_IN_MS = 60000;

function humanReadableTimeFormat(diff){
    let hours = Math.floor( ( diff / (3600 * 1000) ) % 24 );
    let minutes = Math.floor( (diff / (1000 * 60 ) ) % 60 );
    let seconds = Math.floor( (diff / 1000) %  60 );
    return hours+"-"+minutes+"-"+seconds;
}



function exitSession(exitCount){
    xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(){
        if(this.readyState == 4 && this.status == 200){
            var result = JSON.parse(this.responseText);
            switch (result.Status){
                case 0 :
                    if(exitCount < 5){
                        setTimeout(function(){
                            exitSession(++exitCount);
                        },1000*5); // try again to exit session
                    }else{
                        document.getElementById("response").innerText = "Failed exiting session after 5 tries...Logging out in 5 seconds";
                        window.location.href = "/student/signout" 
                    }
                    break;
                case 1:
                        window.location.href = "/student/dashboard"
                        break;
            }
        }
    }
    xhr.open("POST","/student/dashboard/session/exitsession");
    xhr.setRequestHeader("content-type","application/json")
    xhr.send(null);
}
function postAttendance(obj){

   
    console.log(obj);

    xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(){
        if(this.readyState == 4 && this.status == 200){
            var result = JSON.parse(this.responseText);
            switch (result.Status){
                case 0 :
                        document.getElementById("attendance-response").innerText = result.Response + "...Try again after 5 seconds";
                        break;
                case 1:
                        isAttPosted = true;
                        let attBtn = document.getElementById("attendance-button");
                        let capBtn = document.getElementById("capture");
                        let parent = attBtn.parentNode;
                        parent.removeChild(attBtn);

                        parent = capBtn.parentNode;
                        parent.removeChild(capBtn)

                        document.getElementById("attendance-response").innerText = result.Response;
                        break;
            }
        }
    }
    xhr.open("POST","/student/dashboard/session/postattendance");
    xhr.setRequestHeader("content-type","application/json")
    xhr.send(obj);

}

function loadPostAttendance(camera,attNumber){
    attBtn = document.createElement("button")
    attBtn.innerHTML = "Post Attendance"
    attBtn.type = "button";
    attBtn.id = "attendance-button"
    attBtn.addEventListener("click",function(){
        let image64 = document.getElementById("canvas").toDataURL('image/png');
        let time = new Date();  
        let obj = JSON.stringify({
            "attendance_number":attNumber,
            "attendance_time":time,
            "image64":image64,
        });
        postAttendance(obj)
    })
    div = document.createElement("div");
    div.id = "attendance-response";
    camera.appendChild(attBtn);
    camera.appendChild(div);
    isAttBtnLoaded = true;
}

function unloadPostAttendance(){
    let attBtn = document.getElementById("attendance-button");
    if (attBtn != null && typeof(attBtn) != undefined){
        let parent = attBtn.parentNode;
        parent.removeChild(attBtn);
    }

    let attRes = document.getElementById("attendance-response");
    let parent = attRes.parentNode;
    parent.removeChild(attRes);
    isAttBtnLoaded = false;
}

function executeTimer(cur,res){

    if(res.StartTime.getTime() <= cur.getTime() && cur.getTime() < res.EndTime.getTime()){ 
      
        let diff = res.EndTime.getTime() - cur.getTime();
        document.getElementById("response").innerHTML = humanReadableTimeFormat(diff);
        let body = document.getElementsByTagName("body")[0];
        if(cur.getTime() < res.Popup1.getTime()){
            return;
        }
        else if(res.Popup1.getTime() <= cur.getTime() && cur.getTime() < res.Popup1.getTime() + 3 * MINUTE_IN_MS){
            if(modelsLoaded & !cameraLoaded){
                camera = loadCamera(body);
                startVideo();
                if(isFaceVisible && isPhotoTaken && !isAttPosted && !isAttBtnLoaded){ 
                    loadPostAttendance(camera,1);  
                }
            }
        }else if(res.Popup1.getTime() + 3* MINUTE_IN_MS < cur.getTime() && cur.getTime() < res.Popup2.getTime()){
            if(modelsLoaded & cameraLoaded){
               stopVideo();
               unloadCamera(body,camera);
            }
            if(isAttBtnLoaded){
                unloadPostAttendance();
            }  

        }else if(res.Popup2.getTime() <= cur.getTime() && cur.getTime() < res.Popup2.getTime() + 3 * MINUTE_IN_MS){
            if(modelsLoaded & !cameraLoaded){
                camera = loadCamera(body);
                startVideo();
                if(isFaceVisible && isPhotoTaken && !isAttPosted && !isAttBtnLoaded){ 
                    loadPostAttendance(camera,2);  
                }
            }
        }else if(res.Popup2.getTime() + 3* MINUTE_IN_MS < cur.getTime() && cur.getTime() < res.Popup3.getTime()){
            if(modelsLoaded & cameraLoaded){
               stopVideo();
               unloadCamera(body,camera);
            }    
            if(isAttBtnLoaded){
                unloadPostAttendance();
            }  
        }else if(res.Popup3.getTime() <= cur.getTime() && cur.getTime() < res.Popup3.getTime() + 3 * MINUTE_IN_MS){
            if(modelsLoaded & !cameraLoaded){
                camera = loadCamera(body);
                startVideo();
                if(isFaceVisible && isPhotoTaken && !isAttPosted && !isAttBtnLoaded){ 
                    loadPostAttendance(camera,2);  
                }
            }
        }else if(res.Popup3.getTime() + 3* MINUTE_IN_MS < cur.getTime()){
            if(modelsLoaded & cameraLoaded){
               stopVideo();
               unloadCamera(body,camera);
            }    
            if(isAttBtnLoaded){
                unloadPostAttendance();
            }  
        }
    }else if(cur.getTime() < res.StartTime.getTime()){
        document.getElementById("response").innerText = "Session not yet started!";
    }else if(res.EndTime.getTime() < cur.getTime()){
        document.getElementById("response").innerText = "Session Ended! Redirecting in 5 seconds!";
        document.getElementById("status").innerText = "CLOSED";

        // setTimeout(function(){
        //     window.location.href = "/student/dashboard";
        // },1000*5)
    }
}


function setTimer(result){
    cur = new Date();
    res = {
        StartTime : new Date(result.StartTime),
        EndTime : new Date(result.EndTime),
        Popup1 : new Date(result.PopUp1),
        Popup2 : new Date(result.PopUp2),
        Popup3 :  new Date(result.PopUp3)
    }
    console.log(res);
    executeTimer(cur,res);
    timer = setInterval(
        function(){
            cur = new Date();
            // console.log(cur);
            executeTimer(cur,res)
        },1000
    );

}

function getSessionDetails(){
    xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(){
        if(this.readyState == 4 && this.status == 200){
            var result = JSON.parse(this.responseText);
            setTimer(result)
        }
    }
    xhr.open("GET","/student/dashboard/session/timerdetails");
    xhr.setRequestHeader("content-type","application/json")
    xhr.send(null);
} 

getSessionDetails();
