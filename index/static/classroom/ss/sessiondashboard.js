
var timer;
var camera;
var attBtn
var isAttBtnLoaded = false;
const MINUTE_IN_MS = 60000;

function humanReadableTimeFormat(diff){
    let hours = Math.floor( ( diff / (3600 * 1000) ) % 24 );
    let minutes = Math.floor( (diff / (1000 * 60 ) ) % 60 );
    let seconds = Math.floor( (diff / 1000) %  60 );
    return hours+"-"+minutes+"-"+seconds;
}

function postAttendance(num){

    let image64 = document.getElementById("canvas").toDataURL('image/png');
   
    let time = new Date();

    obj = JSON.stringify({
        "attendance_number":num,
        "image64":image64,
        "attendance_time":time
    });

    xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(){
        if(this.readyState == 4 && this.status == 200){
            var result = JSON.parse(this.responseText);
        }
    }
    xhr.open("GET","/student/dashboard/session/postattendance");
    xhr.setRequestHeader("content-type","application/json")
    xhr.send(obj);

}



function loadPostAttendance(camera,attNumber){
    attBtn = document.createElement("button")
    attBtn.innerHTML = "Post Attendance"
    attBtn.type = "button";
    attBtn.id = "attbtn"
    attBtn.value = `${attNumber}`
    attBtn.addEventListener("click",postAttendance(this.value))

    div = document.createElement("div");
    div.id = "attendance-response";
    
    camera.appendChild(attBtn);
    camera.appendChild(div);
}

function executeTimer(cur,res){

    if(res.StartTime.getTime() <= cur.getTime() && cur.getTime() < res.EndTime.getTime()){ 
      
        let diff = res.EndTime.getTime() - cur.getTime();
        document.getElementById("response").innerHTML = humanReadableTimeFormat(diff);
        let body = document.getElementsByTagName("body")[0];
        if(cur.getTime() < res.Popup1.getTime()){
            return;
        }
        // else if(res.Popup1.getTime() <= cur.getTime() && cur.getTime() < res.Popup1.getTime() + 3 * MINUTE_IN_MS){
        //     if(modelsLoaded & !cameraLoaded){
        //         camera = loadCamera(body);
        //         startVideo();
        //         if(isFaceVisible && isPhotoTaken){ 
        //             if(isAttBtnLoaded){
        //                 loadPostAttendance();
        //             }
        //         }else{
        //             unloadPostAttendance();
        //         }
                
        //     }
        // }else if(res.Popup1.getTime() + 3* MINUTE_IN_MS < cur.getTime() && cur.getTime() < res.Popup2.getTime()){
        //     if(modelsLoaded & cameraLoaded){
        //        stopVideo();
        //        unloadCamera(body,camera);
        //     }  
        // }else if(res.Popup2.getTime() <= cur.getTime() && cur.getTime() < res.Popup2.getTime() + 3 * MINUTE_IN_MS){
        //     if(modelsLoaded & !cameraLoaded){
        //         camera = loadCamera(body);
        //         startVideo();
        //     }
        // }else if(res.Popup2.getTime() + 3* MINUTE_IN_MS < cur.getTime() && cur.getTime() < res.Popup3.getTime()){
        //     if(modelsLoaded & cameraLoaded){
        //        stopVideo();
        //        unloadCamera(body,camera);
        //     }    
        // }else if(res.Popup3.getTime() <= cur.getTime() && cur.getTime() < res.Popup3.getTime() + 3 * MINUTE_IN_MS){
        //     if(modelsLoaded & !cameraLoaded){
        //         camera = loadCamera(body);
        //         startVideo();
        //     }
        // }else if(res.Popup3.getTime() + 3* MINUTE_IN_MS < cur.getTime()){
        //     if(modelsLoaded & cameraLoaded){
        //        stopVideo();
        //        unloadCamera(body,camera);
        //     }    
        // }
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


// var SendRequestEverySecond;
// const MINUTE = 60*1000

// function sendRequest(){
    
//     xhr = new XMLHttpRequest();
//     xhr.onreadystatechange = function(){
//         if(this.readyState == 4 && this.status == 200){
//             var result = JSON.parse(this.responseText);
//             switch(result.Status){
//                 case -1:
//                     console.log("Error retrieving session status!")
//                     window.location.href = "/student/dashboard";
//                     break;
//                 case 0 :
//                     document.getElementById("response").innerHTML = result.Response
//                     break;
//                 case 1 :
//                     document.getElementById("reponse").innerHTML = result.Response;
//                     break;
//                 case 2:
//                     clearInterval(SendRequestEverySecond)
//                     document.getElementById("reponse").innerHTML = result.Response;
//                     setTimeout(function(){
//                         console.log("Redirecting.....")
//                         window.location.href = "/student/dashboard";
//                     },5000)

//                     break;
                          
//             }
//         }
//     }
//     xhr.open("POST",window.location.href);
//     xhr.setRequestHeader("content-type","application/json")
//     xhr.send(null);
// }

// SendRequestEverySecond = setInterval(sendRequest,MINUTE);
