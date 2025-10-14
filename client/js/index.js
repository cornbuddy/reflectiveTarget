import "htmx";

const MAX_TAPS = 4;

const targetImage = document.querySelector("img");
const sendDataButton = document.getElementById("send");

let tapCounter = 0;
let tapsCoordinates = [];

function makeRequestObject(data, httpMethod = "POST") {
    const headers = {
        "Accept": "application/json",
        "Content-Type": "application/json",
    };
    let obj = {
        method: httpMethod,
        body: JSON.stringify(data),
        headers: headers,
        credentials: "include",
    };
    httpMethod === "GET" && delete obj.body;
    return obj;
};

function drawShot(tap) {
    let shot = document.createElement("div");
    let imageWrapper = document.getElementById("image-wrapper");
    shot.className = "shot";
    shot.style.marginTop = `${tap.y}px`;
    shot.style.marginLeft = `${tap.x}px`;
    imageWrapper.insertBefore(shot, targetImage);
};

function drawShots(response) {
    return response.json().then((objects) => {
        for (let studentResult of objects) {
            for (let shot of studentResult) {
                drawShot(shot);
            }
        }
    });
};

function removeShots() {
    const shots = document.querySelectorAll(".shot");
    Array.prototype.forEach.call(shots, function(shot) {
        shot.remove();
    });
};

targetImage.addEventListener("click", function(event) {
    if (tapCounter >= MAX_TAPS) {
        return;
    }

    const tap = {
        x: event.offsetX,
        y: event.offsetY,
    };
    tapsCoordinates.push(tap);
    drawShot(tap);
    tapCounter++;
    if (tapCounter === MAX_TAPS) {
        sendDataButton.className = "";
    }
});

sendDataButton.addEventListener("click", function() {
    removeShots();
    const init = makeRequestObject(tapsCoordinates);
    fetch("/api/shots", init)
        .then(console.log)
        .catch(console.error);
    const getData = makeRequestObject(null, "GET");
    fetch("/api/shots", getData)
        .then(drawShots)
        .catch(console.error);
    sendDataButton.innerHTML = "Обновить";
});

document.addEventListener("DOMContentLoaded", function() {
    const getData = makeRequestObject(tapsCoordinates, "GET");
    fetch("/api/shots", getData)
        .then(drawShots)
        .catch(console.error);
});
