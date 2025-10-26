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
});
