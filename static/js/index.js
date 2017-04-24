function validateColor(col) {
    return true;
}

function sendRequest() {
    var xhr = new XMLHttpRequest();
    xhr.open('GET', '/send-mood/{color}');
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
    var col = document.getElementById("color-hex-input").value;
    if (!validateColor(col)) {
        return false;
    }

    sendRequest();
}