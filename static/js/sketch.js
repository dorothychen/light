var button;
var uH, uS, uB;

//noprotect
function setup() { 
  createCanvas(400, 650);
	background(255);
	colorMode(HSB,100);
  noStroke();
  textSize(16);
  text("Mood", 180, 520);
	fill(0);
	uH = random(100);
	uS = random(100);
	uB = random(100);
	makeWheel();
	makeRec();
  makeMood();
  button = createButton('Submit');
  button.position(175, 590);
  button.mouseClicked(submit);

  isLive();

} 

function draw() { 
}

// hue and saturation selection
function makeWheel(){
  for(var h=-PI; h<PI; h+=0.01){
		for(var s=0; s<width/2; s+=1){
      var hue = map(h,-PI,PI,0,100);
      var sat = map(s,0,width/2,0,100);
      fill(hue,sat,100);
      ellipse(200+s*cos(h),230+s*sin(h),2,2);
    }
  }
}

// brightness selection
function makeRec() {
	for(var b = 0; b<400; b++) {
		fill(uH, uS, b/4);
		rect(b, 430, 10, 50);
	}
}

function makeMood() {
  fill(uH, uS, uB);
	rect(160, 530, 80, 50);
}

function mouseDragged(){
	if (dist(mouseX, mouseY, 200, 200) < 200) {
		h = Math.atan2((mouseY-200),(mouseX-200));
		s = dist(mouseX, mouseY, 200, 200);
		uH = round(map(h,-PI,PI,0,100));
  	uS = round(map(s,0,width/2,0,100));
    makeRec();
	}
	else if (mouseY > 430 && mouseY < 480) {
		uB = mouseX/4;
	}
	makeMood();
}

//hsb to rgb to hex
function submit() {
  print("HSB: " + uH + " ," + uS + ", " + uB);
  var c = color(uH, uS, uB);
  var r = round(red(c));
  var g = round(green(c));
  var b = round(blue(c));
  print("RGB: " + r + ", " + g + ", " + b);
  var hex = r.toString(16) + g.toString(16) + b.toString(16);
  print("hex: " + hex);
  submitColor(hex);
}

// TODOOOOOO
function validateColor(col) {
    if (col.length != 6) {
        return false;
    }
    return true;
}

function isLive() {
  var xhr = new XMLHttpRequest();
  xhr.open('GET', '/ctrl/is-live');
  xhr.onload = function() {
    if (xhr.status === 200) {
      var resp = JSON.parse(xhr.responseText);
      if (resp["OK"]) {
        document.getElementById("live-indicator").style.display = "none";
      }
      else {
       document.getElementById("live-indicator").style.display = "block"; 
      }
    }
    else {
      console.log('Request failed. Returned status of ' + xhr.status);
    }
  };
  xhr.send();  
}

function sendRequest(c) {
    var xhr = new XMLHttpRequest();
    xhr.open('GET', '/send-mood/' + c);
    xhr.onload = function() {
      if (xhr.status === 200) {
          console.log("success");
       window.location.href = 'thanks.html';
      }
      else {
          console.log('Request failed. Returned status of ' + xhr.status);
      }
    };
    xhr.send();
}

function submitColor(hex) {
    var c = hex;
    if (!validateColor(c)) {
        return false;
    }

    sendRequest(c);
}
