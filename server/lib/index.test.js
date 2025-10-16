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

describe("/health", () => {
    let app, client, container;

    beforeEach(async () => {
        ({ client, container } = await setupTestDb());
        app = makeApp(client);
    });

    afterEach(async () => {
        await client.end();
        await container.stop();
    });

    test("get should fail if connection is not established", async () => {
        await client.end();
        const resp = await supertest(app).get("/health");
        expect(resp.status).toEqual(503);
        expect(resp.headers["content-type"]).toMatch(/json/);
        expect(resp.body.connected).toBeFalsy();
    });

    test("get should succeed when connected to db", async () => {
        const resp = await supertest(app).get("/health");
        expect(resp.status).toEqual(200);
        expect(resp.headers["content-type"]).toMatch(/json/);
        expect(resp.body.connected).toBeTruthy();
    });
});

describe("/shots", () => {

    test("get should succeed", async () => {
        const app = makeApp();
        const resp = await supertest(app).get("/shots");
        expect(resp.status).toEqual(200);
        expect(resp.headers["content-type"]).toMatch(/json/);
    });
});
