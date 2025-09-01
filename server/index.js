"use strict";

const path = require("path");

const {
    renderMainPage,
    getDataFromClient,
    sendResult,
} = require("./routes");

const express = require("express");
const app = express();

const path = path.join(__dirname, "..", "client");
app.use(express.static(path))
app.get("/", renderMainPage);
app.get("/shots", sendResult);
app.post("/shots", getDataFromClient);

const port = process.env.PORT || 8080;
app.listen(port);
