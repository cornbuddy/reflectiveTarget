const { Router } = require("express");

function authzRouter(userModel) {
    const router = Router();
    router.get("/login", (_, res) => res.render("login"));
    router.get("/signup", (_, res) => res.render("signup"));
    router.post("/signup", async (req, res) => {
        const username = req.body.username;
        const user = await userModel.find(username);
        if (user) {
            return res.status(403).write("exists");
        }

        await userModel.save(username, req.body.password)
            .then(() => res.status(201).write("created"))
            .catch(() => res.status(503));
    });
    return router;
}

module.exports = authzRouter;
