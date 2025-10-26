const { pbkdf2Sync, randomBytes } = require("crypto");

const ITERATIONS = 64;
const KEYLEN = 32;
const SALT_SIZE = 32;
const ALGORITHM = "sha512";


function hashPassword(password, salt) {
    const hash = pbkdf2Sync(password, salt, ITERATIONS, KEYLEN, ALGORITHM);
    return hash.toString("hex");
}

function verifyPassword(password, salt, wantHash) {
    const gotHash = hashPassword(password, salt);
    return wantHash === gotHash;
}

function makeSalt() {
    return randomBytes(SALT_SIZE).toString("hex");
}

module.exports = { hashPassword, verifyPassword, makeSalt };
