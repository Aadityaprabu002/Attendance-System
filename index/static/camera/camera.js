/*
    Needed tags
    <video id="player" controls autoplay></video>
    <button type="button" value = "0" id="capture" >Capture</button>
    <div id="face-detection"></div>
    <canvas id="canvas" width=320 height=240></canvas>
*/





//face detection
var canSet = false; // bool to set isFaceVisible
var isFaceVisible = false; // bool to detect face
var modelsLoaded = false; // bool to detect whether models have loaded
var cameraLoaded = false; // bool to see whether camera loaded or not
Promise.all(
    [
        faceapi.nets.tinyFaceDetector.loadFromUri('/static/student/signup/models'),
        faceapi.nets.faceLandmark68Net.loadFromUri('/static/student/signup/models'),
        faceapi.nets.faceRecognitionNet.loadFromUri('/static/student/signup/models'),
        faceapi.nets.faceExpressionNet.loadFromUri('/static/student/signup/models')
    ]
).then(
    function(){
        modelsLoaded = true;
    }
);


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


function loadCamera(e){
    let video = document.createElement("video");
        video.setAttribute("id","player");
        video.setAttribute("controls",false)
        video.setAttribute("autoplay",true);

    let button = document.createElement("button");
        button.setAttribute("id","capture");
        button.setAttribute("value","0");
        button.setAttribute("type","button");
        button.innerHTML = "CAPTURE";

    let div = document.createElement("div");
        div.setAttribute("id","face-detection");

    let canvas = document.createElement("canvas");
        canvas.setAttribute("id",canvas); 
        canvas.width = 320;
        canvas.height = 240;
    
    let parent = document.createElement("div");
    parent.appendChild(video);
    parent.appendChild(button);
    parent.appendChild(div);
    parent.appendChild(canvas);

    e.appendChild(parent);
    cameraLoaded = true;
    return parent;

}

function unloadCamera(camera){
    cameraLoaded = false;
    document.removeChild(camera);
}