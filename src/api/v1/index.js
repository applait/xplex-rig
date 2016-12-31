/**
 * v1 API routes for xplex-internal
 */

var router = require("express").Router();

router.get("/", (req, res) => {
  res.status(200).json({
    version: "v1",
    methods: [
      "GET /users",
      "GET /provision"
    ]
  });
});

router.use("/users", require("./users.js"));
router.use("/provision", require("./provision.js"));

module.exports = router;
