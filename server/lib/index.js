const path = require("path");

const cookieParser = require("cookie-parser");
const morgan = require("morgan");
const express = require("express");
const expressLayouts = require("express-ejs-layouts");

const { healthRouter, shotsRouter } = require("../routers");

function makeApp(client) {
    const app = express();
    const views = path.join(__dirname, "..", "views");
    app.use(expressLayouts);
    app.set("layout", path.join(views, "_layout.ejs"));
    app.set("view engine", "ejs");
    app.set("views", views);

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
