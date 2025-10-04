"use strict";

const cookieParser = require("cookie-parser");
const morgan = require("morgan");
const express = require("express");
const app = express();

const router = require("./router");

app.use(express.json());
app.use(morgan("tiny"));
app.use(cookieParser());

app.use(router);

const port = process.env.PORT || 8080;
app.listen(port);
