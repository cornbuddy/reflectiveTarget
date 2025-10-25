const { Pool } = require("pg");

const { makeApp } = require("./routers");

async function retry(fn, retries = 5, delay = 5000) {
    try {
        console.log(`starting application, retry #${retries}`);
        return await fn();
    } catch (error) {
        if (retries <= 0) {
            throw error;
        }
        await new Promise(resolve => setTimeout(resolve, delay));
        return retry(fn, retries - 1, delay);
    }
}

const PORT = process.env.PORT || 8080;

const client = new Pool();
client.on("error", console.error);

retry(() => makeApp(client))
    .then((app) => {
        console.log("db is init, starting http server");
        app.listen(PORT, () => console.log(`server started on port ${PORT}`));
    })
    .catch(console.error);
