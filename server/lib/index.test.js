"use strict";

const {
    expect,
    describe,
    test,
} = require("@jest/globals");
const supertest = require("supertest");

const { makeApp } = require("./index");

describe("endpoints", () => {
    const app = makeApp();

    test("get /health should succeed", async () => {
        const resp = await supertest(app).get("/health");
        expect(resp.status).toEqual(200);
        expect(resp.headers["content-type"]).toMatch(/json/);
    });

    test("get /shots should succeed", async () => {
        const resp = await supertest(app).get("/shots");
        expect(resp.status).toEqual(200);
        expect(resp.headers["content-type"]).toMatch(/json/);
    });
});
