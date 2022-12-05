
addActive = (page) => $('#page-' + page).addClass('active')

setStatusBtn = (id) => $.get("status?id=" + id, 
    data => {
        if (data.status === "active"){
            $("#status_btn").text("Active")
        } else if (data.status === "inactive"){
            $("#status_btn").text("Stopped (start at boot)")
        } else {
            $("#status_btn").text("Stopped")
        }
    }
)

execAction = (id,type) => $.get("action?id=" + id + "&" + "type=" + type,
    data => {
        if (data.status === "active"){
            location.href='server?message='
        } else if (data.status === "inactive"){
            $("#status_btn").text("Stopped (start at boot)")
        } else {
            //location.href='server?message=Oh no! Something has gone wrong.'
            console.log(data)
        }
    }
)

