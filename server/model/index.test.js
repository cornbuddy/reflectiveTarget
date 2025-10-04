const {
    beforeEach,
    afterEach,
    expect,
    describe,
    test,
} = require("@jest/globals");
const { PostgreSqlContainer } = require("@testcontainers/postgresql");
const { Client } = require("pg");

const { initDatabase } = require("./index");

const IMAGE = "postgres:18-alpine";

describe("initDatabase", () => {
    let client, container;

    beforeEach(async () => {
        container = await new PostgreSqlContainer(IMAGE).start();
        client = new Client({ connectionString: container.getConnectionUri() });
        await client.connect();
        await initDatabase(client);
    });

    afterEach(async () => {
        await client.end();
        await container.stop();
    });

    test("should create tables", async () => {
        const tables = ["shots", "questions", "targets", "users"];
        for (const table of tables) {
            const query = `select * from ${table}`;
            const result = await client.query(query);
            expect(result.rows).toHaveLength(0);
        }
    });
});
