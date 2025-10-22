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
        return user;
    }
}

module.exports = UserModel;
