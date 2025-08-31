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

let drawPoint = function(tap) {
    let point = document.createElement('div');
    let imageWrapper = document.getElementById('image-wrapper');
    point.className = 'shot';
    point.style.marginTop = `${tap.y}px`;
    point.style.marginLeft = `${tap.x}px`;
    imageWrapper.insertBefore(point, targetImage);
};

let removePoints = function() {
    let points = document.querySelectorAll('.point');
    Array.prototype.forEach.call(points, function(point) {
        point.remove();
    });
};

let displayPoints = function(response) {
    return response.json().then((objects) => {
        for (let studentResult of objects)
            for (let point of studentResult)
                drawPoint(point);
    });
};

targetImage.addEventListener('click', function(event) {
    const tap = {
        x: event.offsetX,
        y: event.offsetY
    };
    tapsCoordinates.push(tap);
    drawPoint(tap);
    tapCounter++;
    if (tapCounter === 4)
        sendDataButton.className = '';
});

clearTargetButton.addEventListener('click', function() {
    tapCounter = 0;
    tapsCoordinates = [];
    removePoints();
    sendDataButton.className = 'hidden';
});

sendDataButton.addEventListener('click', function() {
    removePoints();
    const init = generateData(tapsCoordinates);
    fetch('/points', init);
    const getData = generateData(null, 'GET');
    fetch('/points', getData)
        .then((response) => displayPoints(response))
        .catch((error) => console.log(error));
    sendDataButton.innerHTML = 'Обновить';
});
