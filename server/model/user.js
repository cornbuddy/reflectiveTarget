const { hashPassword, makeSalt } = require("../utils");

class UserModel {
    constructor(client) {
        this.client = client;
    }

    async find(username) {
        const q = "SELECT id, username FROM users WHERE username = $1";
        const res = await this.client.query(q, [username]);
        return res.rows[0];
    }

    async save(user) {
        const q = [
            "INSERT INTO users(username, hashed_password, salt)",
            "VALUES($1, $2, $3)",
            "RETURNING *",
        ].join("\n");
        const salt = makeSalt();
        const hashedPassword = hashPassword(user.password, salt);
        const params = [user.username, hashedPassword, salt];
        const res = await this.client.query(q, params);
        return res.rows[0];
    }
}

module.exports = UserModel;
