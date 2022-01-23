function sendRequest(opt,type,action) {
    xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200) {
            var result = JSON.parse(this.responseText);
            if (action == "add") {
                if (type == "department") {
                    document.querySelector("#add-dept-response").innerHTML = result.Response;
                } else if (type == "course") {
                    document.querySelector("#add-course-response").innerHTML = result.Response;
                }
            }else if(action == "delete"){
                if (type == "department") {
                    document.querySelector("#delete-dept-response").innerHTML = result.Response;
                } else if (type == "course") {
                    document.querySelector("#delete-course-response").innerHTML = result.Response;
                }
            }
        }
    }

    if (type == "department") {
        xhr.open("POST", "/admin/handledepartment");
    } else if (type == "course") {
        xhr.open("POST", "/admin/handlecourse");
    }

    xhr.setRequestHeader("content-type", "application/json")
    xhr.send(opt);
}

function deleteDept() {
    let d = document.getElementById("dept-id-delete").value;
    let obj = JSON.stringify({
        "deptid": d.toUpperCase(),
        "action": "delete",
    });

    sendRequest(obj,"department","delete")

}
function deleteCourse() {
    let c = document.getElementById("course-id-delete").value;
    let obj = JSON.stringify({
        "courseid": c.toUpperCase(),
        "action": "delete",
    });

    sendRequest(obj,"course","delete")

}


function addDept() {
    let d = document.getElementById("dept-id-add").value;
    let dn = document.getElementById("dept-name-add").value;


    if (d.length < 2) {
        document.getElementById("add-dept-response").innerHTML = "Department Id should be greater than 1 character"
        return;
    }
    if (dn.length < 20) {
        document.getElementById("add-dept-response").innerHTML = "Department name should be greater than 20 characters"
        return;
    }

    let obj = JSON.stringify({
        "deptid": d.toUpperCase(),
        "deptname": dn.toUpperCase(),
        "action": "add"
    });

    sendRequest(obj, "department","add");
}

function addCourse() {
    let c = document.getElementById("course-id-add").value;
    let cn = document.getElementById("course-name-add").value;

    if (c.length < 5) {
        document.getElementById("add-course-response").innerHTML = "Course Id should be greater than 5 characters"
        return;
    }
    if (cn.length < 20) {
        document.getElementById("add-course-response").innerHTML = "Course name should be greater than 20 characters"
        return;
    }

    let obj = JSON.stringify({
        "courseid": c.toUpperCase(),
        "coursename": cn.toUpperCase(),
        "action": "add"
    });

   console.log(obj);
    sendRequest(obj,"course","add");
}
