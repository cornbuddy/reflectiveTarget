const { Pool } = require("pg");

const { makeApp } = require("./lib");

const client = new Pool();
const app = makeApp(client);
const port = process.env.PORT || 8080;
app.listen(port);
