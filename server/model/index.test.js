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

const DATABASE_NAME = "reflective_target";
const IMAGE = "postgres:17-alpine";
const TIMEOUT = 60 * 1000;

describe("initDatabase", () => {
    var client, container;

    beforeEach(async () => {
        container = await new PostgreSqlContainer(IMAGE)
            .withDatabase(DATABASE_NAME)
            .start();
        client = new Client({ connectionString: container.getConnectionUri() });
        await client.connect();
        await initDatabase(client);
    }, TIMEOUT);

    afterEach(async () => {
        await client.end();
        await container.stop();
    }, TIMEOUT);

    test("should create tables", async () => {
        const tables = ["shots", "questions", "targets", "users"];
        for (const table of tables) {
            const query = `select * from ${table}`;
            const result = await client.query(query);
            expect(result.rows).toHaveLength(0);
        }
    });
}, TIMEOUT);
