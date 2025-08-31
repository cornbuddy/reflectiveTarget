'use strict';

const fs = require('fs');
const parse = require('url').parse;
const join = require('path').join;
const createReadStream = require('fs').createReadStream;
const stat = require('fs').stat;

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

exports.clearShots = (_, res) => {
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

exports.serveFile = (req, resp) => {
    const parsedUrl = parse(req.url);
    const root = __dirname;
    const path = join(root, '..', 'client', parsedUrl.pathname);
    stat(path, (err, stat) => {
        if (err) {
            if (err.code === 'ENOENT')
                exports.show404Page(err, req, resp);
            else
                exports.show500Page(req, resp);
        } else {
            resp.setHeader('Content-Length', stat.size);
            let stream = createReadStream(path);
            stream.pipe(resp);
            stream.on('error', (error) => exports.show500Page(resp, error));
        }
    });
};
