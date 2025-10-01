"use strict";

const fs = require("fs");

const express = require("express");
const router = express.Router();

const indexHtml = fs.readFileSync("../client/index.html");
let reflectionResults = [];

router.get("/", (req, res) => {
    req.session.views = (req.session.views || 0) + 1;
    console.log(`views: ${req.session.views}`);
    const header = {
        "Content-Type": "text/html",
        "Content-Length": Buffer.byteLength(indexHtml),
    };
    res.writeHeader(200, header);
    res.end(indexHtml);
});

router.get("/shots", (_, res) => {
    const textResponse = JSON.stringify(reflectionResults);
    const header = {
        "Content-Type": "application/json",
        "Content-Length": Buffer.byteLength(textResponse),
    };
    res.writeHeader(200, header);
    res.end(textResponse);
});

router.post("/shots", (req, res) => {
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
});

module.exports = router;
