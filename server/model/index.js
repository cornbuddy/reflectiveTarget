const fs = require("node:fs/promises");

async function initDatabase(client) {
    const script = await fs.readFile(
        `${__dirname}/tables.sql`,
        { encoding: "utf8" },
    );
    await client.query(script);
}

module.exports = { initDatabase };
