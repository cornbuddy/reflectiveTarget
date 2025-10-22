const {
    afterEach,
    beforeEach,
    expect,
    describe,
    test,
} = require("@jest/globals");

const { setupTestDb } = require("../test");
const UserModel = require("./user");
const { initDatabase } = require("./index");

describe("UserModel", () => {
    let client, container, model;

    beforeEach(async () => {
        ({ client, container } = await setupTestDb());
        await initDatabase(client);
        model = new UserModel(client);
    });

    afterEach(async () => {
        await client.end();
        await container.stop();
    });

    test("constructor should set client as property", () => {
        expect(model.client).toBeTruthy();
    });

    test(".find should return user if user exists", async () => {
        const q = "INSERT INTO users(username, password) VALUES($1, $2)";
        const username = "kek";
        const params = [username, "kek"];
        await client.query(q, params);
        const res = await model.find(username);
        expect(res.username).toEqual(username);
        expect(res.id).toEqual(1);
    });

    test(".find should return undefined if user doesn't exist", async () => {
        const res = await model.find("kek");
        expect(res).toBeUndefined();
    });

    test.skip(".save should hash password", async () => {
        const userObj = {
            username: "kek",
            password: "kek",
        };
        const user = await model.save(userObj);
        const q = "SELECT username, password FROM users WHERE username = $1";
        const values = [user.username];
        const res = await client.query(q, values);
        const gotUser = res.rows[0];
        expect(gotUser.id).toEqual(1);
        expect(gotUser.username).toEqual(userObj.username);
    });
});
