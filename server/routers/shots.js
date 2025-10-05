const { Router } = require("express");

let reflectionResults = [];

function router() {
    const route = Router();

    route.get("/", (_, res) => {
        const textResponse = JSON.stringify(reflectionResults);
        const header = {
            "Content-Type": "application/json",
            "Content-Length": Buffer.byteLength(textResponse),
        };
        res.writeHeader(200, header);
        res.end(textResponse);
    });

    route.post("/", (req, res) => {
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

    return route;
}


module.exports = router;
