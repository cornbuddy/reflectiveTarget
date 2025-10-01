"use strict";

const join = require("path").join;

const cookieParser = require("cookie-parser");
const morgan = require("morgan");
const express = require("express");
const app = express();

const router = require("./routes");

const path = join(__dirname, "..", "client");
app.use(express.static(path));
app.use(express.json());
app.use(morgan("tiny"));
app.use(cookieParser());

app.use(router);

const port = process.env.PORT || 8080;
app.listen(port);
