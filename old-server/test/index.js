const { PostgreSqlContainer } = require("@testcontainers/postgresql");
const { Client } = require("pg");

async function setupTestDb() {
    const image = "postgres:18-alpine";
    const container = await new PostgreSqlContainer(image).start();
    const client = new Client({
        connectionString: container.getConnectionUri(),
    });

    await client.connect();
    return { container, client };
}

module.exports = { setupTestDb };
