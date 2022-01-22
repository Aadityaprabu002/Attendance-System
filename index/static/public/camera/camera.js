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
var isPhotoTaken = false;
var modelsLoaded = false; // bool to detect whether models have loaded
var cameraLoaded = false; // bool to see whether camera loaded or not
var player;
var canvas;
var captureButton;

const constraints = {
    video: true,
};

Promise.all(
    [
        faceapi.nets.tinyFaceDetector.loadFromUri('/static/camera/facedetection/models'),
        faceapi.nets.faceLandmark68Net.loadFromUri('/static/camera/facedetection/models'),
        faceapi.nets.faceRecognitionNet.loadFromUri('/static/camera/facedetection/models'),
        faceapi.nets.faceExpressionNet.loadFromUri('/static/camera/facedetection/models')
    ]
).then(
    function(){
        modelsLoaded = true;
    }
);



// start video
function startVideo(){
    canSet = true;
    navigator.mediaDevices.getUserMedia(constraints)
        .then((stream) => {
        player.srcObject = stream;
    });
}



// stop video
function stopVideo(){
    canSet = false;
    player.srcObject.getVideoTracks().forEach(track => track.stop());
    player.pause();
    player.style.display = 'none';
}


function loadCamera(e){
    let video = document.createElement("video");
        video.setAttribute("id","player");
        video.setAttribute("controls",false)
        video.setAttribute("autoplay",true);

    let button = document.createElement("button");
        button.setAttribute("id","capture");
        button.setAttribute("value","0");
        button.setAttribute("type","button");
        button.setAttribute("class","btn btn-danger");
        button.innerHTML = "Capture";

    let div = document.createElement("div");
        div.setAttribute("id","face-detection");

    let canvas = document.createElement("canvas");
        canvas.setAttribute("id","canvas"); 
        canvas.width = 320;
        canvas.height = 240;
    
    let parent = document.createElement("div");
    parent.setAttribute("id","camera");
    parent.appendChild(video);
    parent.appendChild(button);
    parent.appendChild(div);
    parent.appendChild(canvas);

    e.appendChild(parent);
    cameraLoaded = true;

    // video 
    player = document.getElementById('player');
    player.removeAttribute('controls') ;  
    player.addEventListener('pause',()=>{
        console.log('Video stopped');
    })
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
    
    
    //canvas
    canvas = document.getElementById('canvas');
    context = canvas.getContext('2d');

    //capture button
    captureButton = document.getElementById('capture');
    captureButton.addEventListener('click', () => {
        if(!modelsLoaded){
            console.log('Waiting to load models!')
            return;
        }
        // Draw the video frame to the canvas.
        if(captureButton.value == "1"){
            console.log('Camera enabled');
            startVideo();
            player.style.display = '';
            captureButton.innerText = 'Capture';
            captureButton.value = "0";
        }else{
            isPhotoTaken = true;
            console.log('Camera disabled');
            context.drawImage(player, 0, 0, canvas.width, canvas.height);
            stopVideo();
            captureButton.innerText = 'Retake';
            captureButton.value = "1";
        }
    });
    
    return parent;

}

function unloadCamera(parent,camera){
    cameraLoaded = false;
    parent.removeChild(camera);
}