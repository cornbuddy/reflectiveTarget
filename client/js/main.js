'use strict';

let tapCounter = 0;
let tapsCoordinates = [];

const targetImage = document.querySelector('img');
const clearTargetButton = document.getElementById('clear');
const sendDataButton = document.getElementById('send');

let generateData = function(data, httpMethod = 'POST') {
    const headers = {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
    };
    let obj = {
        method: httpMethod,
        body: JSON.stringify(data),
        headers: headers
    };
    httpMethod === 'GET' && delete obj.body;
    return obj;
};

let drawShot = function(tap) {
    let shot = document.createElement('div');
    let imageWrapper = document.getElementById('image-wrapper');
    shot.className = 'shot';
    shot.style.marginTop = `${tap.y}px`;
    shot.style.marginLeft = `${tap.x}px`;
    imageWrapper.insertBefore(shot, targetImage);
};

let removeShots = function() {
    let shots = document.querySelectorAll('.shot');
    Array.prototype.forEach.call(shots, function(shot) {
        shot.remove();
    });
};

let displayShots = function(response) {
    return response.json().then((objects) => {
        for (let studentResult of objects)
            for (let shot of studentResult)
                drawShot(shot);
    });
};

targetImage.addEventListener('click', function(event) {
    const tap = {
        x: event.offsetX,
        y: event.offsetY
    };
    tapsCoordinates.push(tap);
    drawShot(tap);
    tapCounter++;
    if (tapCounter === 4)
        sendDataButton.className = '';
});

clearTargetButton.addEventListener('click', function() {
    tapCounter = 0;
    tapsCoordinates = [];
    removeShots();
    sendDataButton.className = 'hidden';
});

sendDataButton.addEventListener('click', function() {
    removeShots();
    const init = generateData(tapsCoordinates);
    fetch('/shots', init);
    const getData = generateData(null, 'GET');
    fetch('/shots', getData)
        .then((response) => displayShots(response))
        .catch((error) => console.log(error));
    sendDataButton.innerHTML = 'Обновить';
});
