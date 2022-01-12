var SendRequestEverySecond;

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

SendRequestEverySecond = setInterval(sendRequest,1000);
