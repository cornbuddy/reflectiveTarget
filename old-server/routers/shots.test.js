const {
    afterEach,
    beforeEach,
    expect,
    describe,
    test,
} = require("@jest/globals");
const supertest = require("supertest");

const { setupTestDb } = require("../test");
const { makeApp } = require("./index");

describe("http endpoints", () => {
    let app, client, container;

    beforeEach(async () => {
        ({ client, container } = await setupTestDb());
        app = await makeApp(client);
    });

    afterEach(async () => {
        await client.end();
        await container.stop();
    });

    describe("/api/shots", () => {
        test("get should succeed", async () => {
            const resp = await supertest(app).get("/api/shots");
            expect(resp.status).toEqual(200);
            expect(resp.headers["content-type"]).toMatch(/json/);
        });
    });
});
