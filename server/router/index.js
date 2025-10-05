const express = require("express");
const router = express.Router();

let reflectionResults = [];

router.get("/health", (_, res) => {
    res.status(503).json({ connected: false });
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
