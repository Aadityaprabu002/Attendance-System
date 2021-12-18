var btn = document.getElementById("submit");
btn.addEventListener("click",submitForm);


function sendRequest(obj){
    xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function(){
        if(this.readyState == 4 && this.status == 200){
            var result  = JSON.parse(this.responseText);
            switch(result.Status){
                case 0 : document.querySelector("#response").innerHTML =  result.Response; break;
                case 1 : window.location.href = "/classroom"; break;
            }
        }
    }
    xhr.open("POST","/joinclassroom");
    xhr.setRequestHeader("content-type","application/json")
    xhr.send(obj);
}

function submitForm(){
   
    let obj = JSON.stringify({
       "classroomid": document.getElementById("classroom-id").value,
       "joiningtime": new Date().toJSON()
    })
    console.log(obj);
    sendRequest(obj);
}

