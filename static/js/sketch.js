var button;
var uH, uS, uB;
var centerX, centerY; // color wheel center

//noprotect
function setup() { 
  var canvas = createCanvas(400, 650);
  // Move the canvas so it's inside our <div id="sketch">.
  canvas.parent('sketch');

  background(255);
  colorMode(HSB,100);
  noStroke();
  textSize(16);
  text("Mood", 180, 500);
	fill(0);
	centerX = 200;
	centerY = 230;
	uH = random(100);
	uS = random(100);
	uB = 100;
	makeWheel();
  makeMood();
  button = createButton('Submit');
  button.parent('button-submit');
  button.mouseOver(function(e) {
    var c = color(uH, uS, uB);
    var rgb = "rgb(" + round(red(c)) + ", " + round(green(c)) + ", " + round(blue(c)) + ")";
    e.toElement.style.borderColor = "gray";
    e.toElement.style.color = "gray";
  });
  button.mouseOut(function(e) {
    e.fromElement.style.borderColor = "black";
    e.fromElement.style.color = "black";
    e.fromElement.style.background = "white";
  });
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
      ellipse(centerX+s*cos(h),centerY+s*sin(h),2,2);
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
	rect(160, 510, 80, 50);
}

function mouseDragged(){
	if (dist(mouseX, mouseY, centerX, centerY) < 200) {
		h = Math.atan2((mouseY-centerY),(mouseX-centerX));
		s = dist(mouseX, mouseY, centerX, centerY);
		uH = round(map(h,-PI,PI,0,100));
  	        uS = round(map(s,0,width/2,0,100));
		makeMood();
	}
	
}

function mouseClicked() {
	if (dist(mouseX, mouseY, centerX, centerY) < 200) {
		h = Math.atan2((mouseY-centerY),(mouseX-centerX));
		s = dist(mouseX, mouseY, centerX, centerY);
		uH = round(map(h,-PI,PI,0,100));
  	        uS = round(map(s,0,width/2,0,100));
		makeMood();
	}
}

function vToHex(v_int) {
  v = v_int.toString(16);
  if (v.length == 1) return "0" + v;
  else return v;
}

//hsb to rgb to hex
function submit() {
  print("HSB: " + uH + " ," + uS + ", " + uB);
  var c = color(uH, uS, uB);
  var r = round(red(c));
  var g = round(green(c));
  var b = round(blue(c));
  print("RGB: " + r + ", " + g + ", " + b);
  var hex = vToHex(r) + vToHex(g) + vToHex(b);
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
      console.log("invalid color");
      return false;
    }

    sendRequest(c);
}
