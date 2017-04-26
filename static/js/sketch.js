var slider, button;
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
	uH = 0;
	uS = 0;
	uB = 0;
  textSize(20);
  text("Pick a color that represents your mood!", 40, 20);
	makeWheel();
	makeRec();
  slider = createSlider(0, 100, 0);
	slider.style('width','400px');
	slider.position(0, 470);
  makeMood();
  button = createButton('Submit');
  button.position(175, 590);
  button.mousePressed(submit);
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

function mouseClicked(){
	if (dist(mouseX, mouseY, 200, 200) < 200) {
		h = Math.atan2((mouseY-200),(mouseX-200));
		s = dist(mouseX, mouseY, 200, 200);
		uH = round(map(h,-PI,PI,0,100));
  	uS = round(map(s,0,width/2,0,100));
    makeRec();
	}
	uB = slider.value();
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
