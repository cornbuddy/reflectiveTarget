'use strict';

const http = require('http');

const renderMainPage = require('./routes').renderMainPage;
const getDataFromClient = require('./routes').getDataFromClient;
const sendResult = require('./routes').sendResult;
const show404Page = require('./routes').show404Page;
const show500Page = require('./routes').show500Page;
const serveFile = require('./routes').serveFile;
const clearShots = require('./routes').clearShots;


let isStatic = (url) => {
    const formats = ['js', 'gif', 'css', 'ico'];
    for (const format of formats)
        if (url.endsWith(format))
            return true;
    return false;
};

http.createServer((request, response) => {
    const method = request.method;
    const url = request.url;
    request.on('error', (error) => show500Page(error, response));
    // get methods
    if (method === 'GET' && url === '/')
        renderMainPage(request, response);
    else if (method === 'GET' && url === '/shots')
        sendResult(response);
    else if (method === 'GET' && url === '/clear-shots')
        clearShots(request, response);
    // post methods
    else if (method === 'POST' && url === '/shots')
        getDataFromClient(request, response);
    // serve file
    else if (method === 'GET' && isStatic(url))
        serveFile(request, response);
    else
        show404Page('no router', request, response);
}).listen(process.env.PORT || 8080);
