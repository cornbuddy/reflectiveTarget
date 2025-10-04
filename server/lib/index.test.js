"use strict";

const {
    expect,
    describe,
    test,
} = require("@jest/globals");
const supertest = require("supertest");

const { makeApp } = require("./index");

describe("router", () => {
    test("get /shots should succeed", async () => {
        const app = makeApp();
        const resp = await supertest(app).get("/shots");
        expect(resp.status).toEqual(200);
        expect(resp.headers["content-type"]).toMatch(/json/);
    });
});
