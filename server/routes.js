"use strict";

const fs = require("fs");

const index = fs.readFileSync("../client/index.html");
let reflectionResults = [];

exports.index = (req, res) => {
    req.session.views = (req.session.views || 0) + 1;
    console.log(`views: ${req.session.views}`);
    const header = {
        "Content-Type": "text/html",
        "Content-Length": Buffer.byteLength(index),
    };
    res.writeHeader(200, header);
    res.end(index);
};

exports.saveShots = (req, res) => {
    const shots = +req.cookies.shots;
    if (shots > 0) {
        const msg = "not ok";
        const header = {
            "Content-Type": "text/html",
            "Content-Length": Buffer.byteLength(msg),
        };
        res.writeHeader(406, header);
        res.end(msg);
        return;
    }

    const obj = req.body;
    reflectionResults.push(obj);

    const msg = "ok";
    const header = {
        "Content-Type": "text/html",
        "Content-Length": Buffer.byteLength(msg),
        "Set-Cookie": `shots=${obj.length}`,
    };
    res.writeHeader(200, header);
    res.end(msg);
};

exports.getShots = (_, res) => {
    const textResponse = JSON.stringify(reflectionResults);
    const header = {
        "Content-Type": "application/json",
        "Content-Length": Buffer.byteLength(textResponse),
    };
    res.writeHeader(200, header);
    res.end(textResponse);
};
