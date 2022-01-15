
var timer;
var camera;
const MINUTE_IN_MS = 60000;

function sendAttendance(){}


function executeTimer(cur,res){
    if(res.StartTime.getTime() <= cur.getTime() && cur.getTime() < res.EndTime.getTime()){ 
        //assuming both have same day
        let diff = res.EndTime.getTime() - cur.getTimer();
        let hours = Math.floor( ( diff / (3600 * 1000) ) % 24 );
        let minutes = Math.floor( (diff % (1000 * 60 ) ) % 60 );
        let seconds = Math.floor( (diff / 1000) %  60 );
        document.getElementById("response").innerHTML = hours+"-"+minutes+"-"+"-"+seconds;
        if(res.Popup1.getTime() <= cur.getTime() && cur.getTime() < res.Popup1.getTime() + 3 * MINUTE_IN_MS){
            if(modelsLoaded & !cameraLoaded){
                camera = loadCamera(document.getElementsByTagName("body"));
                startVideo();
            }
        }else if(res.Popup1.getTime() + 3* MINUTE_IN_MS < cur.getTime()){
            if(modelsLoaded & cameraLoaded){
               stopVideo();
               unloadCamera(camera)
            }
            
        }

        if(res.Popup2.getTime() <= cur.getTime() && cur.getTime() < res.Popup2.getTime() + 3 * MINUTE_IN_MS){
            if(modelsLoaded & !cameraLoaded){
                camera = loadCamera(document.getElementsByTagName("body"));
                startVideo();
            }
        }else if(res.Popup2.getTime() + 3* MINUTE_IN_MS < cur.getTime()){
            if(modelsLoaded & cameraLoaded){
               stopVideo();
               unloadCamera(camera)
            }    
        }

        if(res.Popup3.getTime() <= cur.getTime() && cur.getTime() < res.Popup3.getTime() + 3 * MINUTE_IN_MS){
            if(modelsLoaded & !cameraLoaded){
                camera = loadCamera(document.getElementsByTagName("body"));
                startVideo();
            }
        }else if(res.Popup3.getTime() + 3* MINUTE_IN_MS < cur.getTime()){
            if(modelsLoaded & cameraLoaded){
               stopVideo();
               unloadCamera(camera)
            }    
        }
        
        



    }else if(cur.getTime() < res.Start.getTime()){
        document.getElementById("response").innerHTML = "Session not yet started!";
    }
    else if(res.EndTime.getTime() < cur.getTime()){
        document.getElementById("response").innerHTML = "Session Ended! Redirecting in 5 seconds!";
        setTimeout(function(){
            window.location.href = "/student/dashboard/";
        },1000*5)
    }
}


function setTimer(res){
    cur = new Date();
    executeTimer(cur,res);
    timer = setInterval(
        function(){
            cur = new Date();
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


var SendRequestEverySecond;
const MINUTE = 60*1000

function sendRequest(){
    
    xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(){
        if(this.readyState == 4 && this.status == 200){
            var result = JSON.parse(this.responseText);
            switch(result.Status){
                case -1:
                    console.log("Error retrieving session status!")
                    window.location.href = "/student/dashboard";
                    break;
                case 0 :
                    document.getElementById("response").innerHTML = result.Response
                    break;
                case 1 :
                    document.getElementById("reponse").innerHTML = result.Response;
                    break;
                case 2:
                    clearInterval(SendRequestEverySecond)
                    document.getElementById("reponse").innerHTML = result.Response;
                    setTimeout(function(){
                        console.log("Redirecting.....")
                        window.location.href = "/student/dashboard";
                    },5000)

                    break;
                          
            }
        }
    }
    xhr.open("POST",window.location.href);
    xhr.setRequestHeader("content-type","application/json")
    xhr.send(null);
}

SendRequestEverySecond = setInterval(sendRequest,MINUTE);
