"use strict";

const cookieParser = require("cookie-parser");
const morgan = require("morgan");
const express = require("express");

const router = require("../router");

function makeApp() {
    const app = express();
    app.use(express.json());
    app.use(morgan("tiny"));
    app.use(cookieParser());
    app.use(router);
    return app;
}

module.exports = { makeApp };
