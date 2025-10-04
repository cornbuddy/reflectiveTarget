"use strict";

const { makeApp } = require("./lib");

const app = makeApp();
const port = process.env.PORT || 8080;
app.listen(port);
