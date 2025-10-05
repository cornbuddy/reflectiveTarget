const { Router } = require("express");

function router(client) {
    const route = Router();
    route.get("/", async (_, res) => {
        await client.query("SELECT 1")
            .then(() => res.status(200).json({ connected: true }))
            .catch(() => res.status(503).json({ connected: false }));
    });

    return route;
}

module.exports = router;
