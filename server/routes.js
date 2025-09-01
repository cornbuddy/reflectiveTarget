'use strict';

const fs = require('fs');

const index = fs.readFileSync('../client/index.html');
let reflectionResults = [];

exports.renderMainPage = (_, res) => {
    const header = {
        'Content-Type': 'text/html',
        'Content-Length': Buffer.byteLength(index)
    };
    res.writeHeader(200, header);
    res.end(index);
};

exports.getDataFromClient = (req, res) => {
    let reflectionResult = '';
    req.setEncoding('utf-8');
    req.on('data', (chunk) => reflectionResult += chunk);
    req.on('end', () => {
        reflectionResults.push(JSON.parse(reflectionResult));
        const header = {
            'Content-Type': 'text/html',
            'Content-Length': Buffer.byteLength('OK')
        };
        res.writeHeader(200, header);
        res.end('OK');
    });
};

exports.sendResult = (_, res) => {
    const textResponse = JSON.stringify(reflectionResults);
    const header = {
        'Content-Type': 'application/json',
        'Content-Length': Buffer.byteLength(textResponse)
    };
    res.writeHeader(200, header);
    res.end(textResponse);
};
