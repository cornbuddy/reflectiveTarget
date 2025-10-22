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
        app = makeApp(client);
    });

    afterEach(async () => {
        await client.end();
        await container.stop();
    });

    test.each([
        { url: "/", text: "Reflective target" },
        { url: "/login", text: "Login" },
        { url: "/signup", text: "Signup" },
    ])("view $url", async ({ url, text }) => {
        const resp = await supertest(app).get(url);
        expect(resp.status).toEqual(200);
        expect(resp.headers["content-type"]).toMatch(/html/);
        expect(resp.text).toContain(text);
    });

    describe("/api/health", () => {
        test("get should fail when db is broken", async () => {
            await client.end();
            const resp = await supertest(app).get("/api/health");
            expect(resp.status).toEqual(503);
            expect(resp.headers["content-type"]).toMatch(/json/);
            expect(resp.body.connected).toBeFalsy();
        });

        test("get should succeed when connected to db", async () => {
            const resp = await supertest(app).get("/api/health");
            expect(resp.status).toEqual(200);
            expect(resp.headers["content-type"]).toMatch(/json/);
            expect(resp.body.connected).toBeTruthy();
        });
    });

    describe("/api/shots", () => {
        test("get should succeed", async () => {
            const resp = await supertest(app).get("/api/shots");
            expect(resp.status).toEqual(200);
            expect(resp.headers["content-type"]).toMatch(/json/);
        });
    });
});
