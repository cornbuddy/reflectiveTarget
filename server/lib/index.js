const cookieParser = require("cookie-parser");
const morgan = require("morgan");
const express = require("express");

const { healthRouter, shotsRouter } = require("../routers");

function makeApp(client) {
    const app = express();
    app.set("view engine", "ejs");
    app.set("views", "../views");
    app.use(express.json());
    app.use(morgan("tiny"));
    app.use(cookieParser());
    app.use("/health", healthRouter(client));
    app.use("/shots", shotsRouter(client));
    return app;
}

module.exports = { makeApp };
