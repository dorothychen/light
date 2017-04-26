function validateColor(col) {
    if (col.length != 6) {
        return false;
    }
    return true;
}

function sendRequest(c) {
    var xhr = new XMLHttpRequest();
    xhr.open('GET', '/send-mood/' + c);
    xhr.onload = function() {
        if (xhr.status === 200) {
            console.log("success")
        }
        else {
            console.log('Request failed. Returned status of ' + xhr.status);
        }
    };
    xhr.send();
}

function submitColor() {
    var c = document.getElementById("color-hex-input").value;
    if (!validateColor(c)) {
        return false;
    }

    sendRequest(c);
}