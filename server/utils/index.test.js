const {
    expect,
    describe,
    test,
} = require("@jest/globals");

const { hashPassword, verifyPassword, makeSalt } = require("./index");

const table = [{
    password: "kek",
    salt: "kek",
    hashed: "cb9ac34dcbf56f6d81b7fe25f15a2bb4bf71e587cb16b034d805a8d9a15dd266",
}, {
    password: "not kek",
    salt: "not kek",
    hashed: "4c0fa623535cfe0b38c0278c8e94803f4c73962affedfe7265cca4a746749b8d",
}];

describe("hashPassword", () => {
    test.each(table)("hash `$password`", ({ password, salt, hashed }) => {
        const got = hashPassword(password, salt);
        expect(got).toEqual(hashed);
    });
});

describe("verifyPassword", () => {
    test.each(table)("verify `$password`", ({ password, salt, hashed }) => {
        const got = verifyPassword(password, salt, hashed);
        expect(got).toBeTruthy();
    });

    test.each([{
        password: "kek",
        salt: "kek",
        hashed: "kek",
    }])("not verify `$password`", ({ password, salt, hashed }) => {
        const got = verifyPassword(password, salt, hashed);
        expect(got).toBeFalsy();
    });
});

describe("makeSalt", () => {
    test("should produce string with fixed width", () => {
        const got = makeSalt();
        expect(got).toHaveLength(64);
    });

    test("should produce unique strings", () => {
        const got1 = makeSalt();
        const got2 = makeSalt();
        expect(got1).not.toEqual(got2);
    });
});
