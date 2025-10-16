const path = require("path");

const cookieParser = require("cookie-parser");
const morgan = require("morgan");
const express = require("express");

const { healthRouter, shotsRouter } = require("../routers");

function makeApp(client) {
    const app = express();
    app.set("view engine", "ejs");
    app.set("views", path.join(__dirname, "..", "views"));
    app.use(express.json());
    app.use(morgan("tiny"));
    app.use(cookieParser());
    app.use("/api/health", healthRouter(client));
    app.use("/api/shots", shotsRouter(client));
    app.get("/", (_, res) => {
        res.render("index");
    });
    return app;
}

module.exports = { makeApp };
