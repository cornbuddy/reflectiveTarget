const { Router } = require("express");

function authzRouter(userModel) {
    const router = Router();
    router.get("/login", (_, res) => res.render("login"));
    router.get("/signup", (_, res) => res.render("signup"));
    router.post("/signup", async (req, res) => {
        const user = req.body;
        const found = await userModel.find(user.username);
        if (found) {
            return res.status(403).write("exists");
        }

        userModel.save(user)
            .then(() => res.status(201).write("created"))
            .catch((err) => {
                console.error(err);
                res.status(503).write("error");
            });
    });

    return router;
}

module.exports = authzRouter;
