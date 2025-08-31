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
    point.className = 'point';
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


let showTarget = function() {
    const blitzWrapper = document.getElementById('blitz-wrapper');
    blitzWrapper.className = 'hidden';
    const targetWrapper = document.getElementById('target-wrapper');
    targetWrapper.className = 'centered';
};

let getBlitzResult = function() {
    return {
        'lastName': document.getElementById('last-name').value,
        'group': document.getElementById('group').value,
        'first': document.getElementById('first').value,
        'second': document.getElementById('second').value,
        'third': document.getElementById('third').value,
        'fourth': document.getElementById('fourth').value,
        'fiveth': document.getElementById('fiveth').value
    };
};

const sendBlitzButton = document.getElementById('send-blitz');

sendBlitzButton.addEventListener('click', function() {
    const blitzResult = getBlitzResult();
    const init = generateData(blitzResult);
    fetch('/blitz', init);
    showTarget();
});
