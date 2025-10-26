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

describe("/signup", () => {
    let app, client, container;

    beforeEach(async () => {
        ({ client, container } = await setupTestDb());
        app = await makeApp(client);
    });

    afterEach(async () => {
        await client.end();
        await container.stop();
    });

    test("should create new user if not exists", async () => {
        const resp = await createUser(app, "kek", "kek");
        expect(resp.status).toEqual(201);
        expect(resp.text).toContain("created");
    });

    test("should respond with error when db is broken", async () => {
        await client.end();
        const resp = await createUser(app, "kek", "kek");
        expect(resp.status).toEqual(503);
        expect(resp.text).toContain("error");
    });

    test("should reject when user already exists", async () => {
        await createUser(app, "kek", "kek").expect(201);
        const resp = createUser(app, "kek", "kek");
        expect(resp.status).toEqual(403);
        expect(resp.text).toContain("exists");
    });
});

async function createUser(app, username, password) {
    return supertest(app).post("/signup")
        .field("username", username)
        .field("password", password);
}
