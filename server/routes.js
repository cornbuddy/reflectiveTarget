'use strict';

const fs = require('fs');
const parse = require('url').parse;
const join = require('path').join;
const createReadStream = require('fs').createReadStream;
const stat = require('fs').stat;

const index = fs.readFileSync('../client/index.html');
let reflectionResults = [];
let blitzResult = [];

exports.renderMainPage = (req, res) => {
    const header = {
        'Content-Type': 'text/html',
        'Content-Length': Buffer.byteLength(index)
    };
    res.writeHeader(200, header);
    res.end(index);
};

exports.sendBlitz = (req, res) => {
    const textResponse = JSON.stringify(blitzResult);
    const header = {
        'Content-Type': 'application/json',
        'Content-Length': Buffer.byteLength(textResponse)
    };
    res.writeHeader(200, header);
    res.end(textResponse);
};

exports.clearBlitz = (req, res) => {
    blitzResult = [];
    const header = {
        'Content-Type': 'text/html',
        'Content-Length': Buffer.byteLength('OK')
    };
    res.writeHeader(200, header);
    res.end('OK');
};

exports.getBlitzResult = (req, res) => {
    let blitz = '';
    req.setEncoding('utf-8');
    req.on('data', (chunk) => blitz += chunk);
    req.on('end', () => {
        blitzResult.push(JSON.parse(blitz));
        const header = {
            'Content-Type': 'text/html',
            'Content-Length': Buffer.byteLength('OK')
        };
        res.writeHeader(200, header);
        res.end('OK');
    });
};

exports.clearPoints = (req, res) => {
    reflectionResults = [];
    const header = {
        'Content-Type': 'text/html',
        'Content-Length': Buffer.byteLength('OK')
    };
    res.writeHeader(200, header);
    res.end('OK');
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

exports.sendResult = (res) => {
    const textResponse = JSON.stringify(reflectionResults);
    const header = {
        'Content-Type': 'application/json',
        'Content-Length': Buffer.byteLength(textResponse)
    };
    res.writeHeader(200, header);
    res.end(textResponse);
};

exports.show404Page = (error, req, res) => {
    console.log(`404 url: ${req.method} ${req.url}, err: ${error};`);
    const header = {
        'Content-Type': 'text/html',
        'Content-Length': Buffer.byteLength('404')
    };
    res.writeHeader(404, header);
    res.end('404');
};

exports.show500Page = (error, res) => {
    console.log(`500 err: ${error};`);
    const header = {
        'Content-Type': 'text/html',
        'Content-Length': Buffer.byteLength('500')
    };
    res.writeHeader(500, header);
    res.end('500');
};

exports.serveFile = (request, response) => {
    const parsedUrl = parse(request.url);
    const root = __dirname;
    let path = join(root, '..', 'client', parsedUrl.pathname);
    stat(path, (err, stat) => {
        if (err) {
            if (err.code === 'ENOENT')
                exports.show404Page(err, request, response);
            else
                exports.show500Page(request, response);
        } else {
            response.setHeader('Content-Length', stat.size);
            let stream = createReadStream(path);
            stream.pipe(response);
            stream.on('error', (error) => exports.show500Page(response, error));
        }
    });
};
