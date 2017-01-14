/**
 * v1 API routes for xplex-internal
 */

var router = require("express").Router();

router.get("/", (req, res) => {
  res.status(200).json({
    version: "v1",
    methods: [
      "GET /agents",
      "GET /provision",
      "GET /users"
    ]
  });
});

router.use("/agents", require("./agents"));
router.use("/provision", require("./provision"));
router.use("/users", require("./users"));

module.exports = router;
