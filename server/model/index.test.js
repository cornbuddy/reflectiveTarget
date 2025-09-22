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
    });

    afterAll(async () => {
        await client.end();
        await container.stop();
    });

    test("db should be connectable", async () => {
        const result = await client.query("SELECT 1");
        expect(result.rows[0]).toEqual({ "?column?": 1 });
    });

    test("should create database", async () => {
        expect(initDatabase(client)).toBe(undefined);
    });

}, TIMEOUT);
