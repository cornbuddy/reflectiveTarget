"use strict";

const join = require("path").join;

const {
    index,
    saveShots,
    getShots,
} = require("./routes");

const cookieParser = require("cookie-parser");
const morgan = require("morgan");
const express = require("express");
const app = express();

const path = join(__dirname, "..", "client");
app.use(express.static(path));
app.use(express.json());
app.use(morgan("tiny"));
app.use(cookieParser());

app.get("/", index);
app.get("/shots", getShots);
app.post("/shots", saveShots);

const port = process.env.PORT || 8080;
app.listen(port);
