const {
    beforeAll,
    afterAll,
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

    beforeAll(async () => {
        container = await new PostgreSqlContainer(IMAGE).start();
        client = new Client({
            connectionString: container.getConnectionUri(),
        });
        await client.connect();
        await initDatabase(client, DATABASE_NAME);
    });

    afterAll(async () => {
        await client.end();
        await container.stop();
    });

    test("should create database", async () => {
        const result = await client.query("SELECT datname FROM pg_database;");
        expect(result.rows).toContainEqual({ "datname": DATABASE_NAME });
    });

    test("should be idempotent", async () => {
        await initDatabase(client, DATABASE_NAME);

        const result = await client.query("SELECT datname FROM pg_database;");
        expect(result.rows).toContainEqual({ "datname": DATABASE_NAME });
    });
}, TIMEOUT);
