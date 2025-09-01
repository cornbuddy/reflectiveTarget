"use strict";

const MAX_TAPS = 4;

const targetImage = document.querySelector("img");
const sendDataButton = document.getElementById("send");

let tapCounter = 0;
let tapsCoordinates = [];

function generateData(data, httpMethod = "POST") {
    const headers = {
        "Accept": "application/json",
        "Content-Type": "application/json"
    };
    let obj = {
        method: httpMethod,
        body: JSON.stringify(data),
        headers: headers
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

function removeShots() {
    let shots = document.querySelectorAll(".shot");
    Array.prototype.forEach.call(shots, function(shot) {
        shot.remove();
    });
};

function displayShots(response) {
    return response.json().then((objects) => {
        for (let studentResult of objects) {
            for (let shot of studentResult) {
                drawShot(shot);
            }
        }
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
    const init = generateData(tapsCoordinates);
    fetch("/shots", init);
    const getData = generateData(null, "GET");
    fetch("/shots", getData)
        .then((response) => displayShots(response))
        .catch((error) => console.log(error));
    sendDataButton.innerHTML = "Обновить";
});
