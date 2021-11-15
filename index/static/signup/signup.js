

//face detection
var canSet = false; // bool to set isFaceVisible
var isFaceVisible = false; // bool to detect face

Promise.all(
    [
        faceapi.nets.tinyFaceDetector.loadFromUri('/static/signup/models'),
        faceapi.nets.faceLandmark68Net.loadFromUri('/static/signup/models'),
        faceapi.nets.faceRecognitionNet.loadFromUri('/static/signup/models'),
        faceapi.nets.faceExpressionNet.loadFromUri('/static/signup/models')
    ]
).then(startVideo);


// video 

const player = document.getElementById('player');
player.removeAttribute('controls') ;  

//canvas
const canvas = document.getElementById('canvas');
const context = canvas.getContext('2d');

//capture button
const captureButton = document.getElementById('capture');

const constraints = {
    video: true,
};

// start video
function startVideo(){
    canSet = true;
    navigator.mediaDevices.getUserMedia(constraints)
        .then((stream) => {
        player.srcObject = stream;
    });
}

player.addEventListener('play',()=>{
    console.log("Video started");
    setInterval(async () =>{
        const detections = await faceapi.detectAllFaces(player,
        new faceapi.TinyFaceDetectorOptions()).withFaceLandmarks()
        if(detections.length == 0 && canSet ){
            isFaceVisible = false;
            document.getElementById('face-detection').innerHTML = 'No face detected!';
        }
        else if(canSet){
            isFaceVisible = true;
            document.getElementById('face-detection').innerHTML = 'Face detected!';
        }
    },100)
});


// stop video
function stopVideo(){
    canSet = false;
    player.srcObject.getVideoTracks().forEach(track => track.stop());
    player.pause();
    player.style.display = 'none';
}
player.addEventListener('pause',()=>{
    console.log('Video stopped');
})



captureButton.addEventListener('click', () => {
    // Draw the video frame to the canvas.
    if(captureButton.value == "1"){
        console.log('Camera enabled');
        startVideo();
        player.style.display = '';
        captureButton.innerText = 'Capture';
        captureButton.value = "0";
    }else{
        console.log('Camera disabled');
        context.drawImage(player, 0, 0, canvas.width, canvas.height);
        stopVideo();
        captureButton.innerText = 'Retake';
        captureButton.value = "1";
    }
});


// form and ajax
function submitForm(){
    if(!isFaceVisible){
        document.getElementById('response').innerHTML = 'Form can not be submitted since no face was detected!';
        return;
    }
    let image64 = document.getElementById("canvas").toDataURL('image/png');
    let fname = document.querySelector("#fname").value;
    let lname = document.querySelector("#lname").value;
    let regnumber = document.querySelector("#rollno").value
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
       "email": email,
       "image64":image64,
       "regnumber":regnumber
    })

    console.log(obj);
    sendRequest(obj);
}

function sendRequest(obj){
    xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(){
        if(this.readyState == 4 && this.status == 200){
            var result = JSON.parse(this.responseText);
            document.querySelector("#response").innerHTML =  result.Response; 
        }
    }
    xhr.open("POST","/signup");
    xhr.setRequestHeader("content-type","application/json")
    xhr.send(obj);
}