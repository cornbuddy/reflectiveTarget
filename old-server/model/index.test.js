const {
    afterEach,
    beforeEach,
    expect,
    describe,
    test,
} = require("@jest/globals");

const { setupTestDb } = require("../test");
const { initDatabase } = require("./index");

describe("initDatabase", () => {
    let client, container;

    beforeEach(async () => {
        ({ client, container } = await setupTestDb());
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
