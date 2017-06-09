/**
 * v1 API routes for xplex-internal
 */

let router = require("express").Router();
const DO = require("../../providers/digital_ocean");

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

router.use((req, res, next) => {
  req.providers = {
    DO: new DO(req.config.server.digital_ocean.token)
  };
  next();
});

router.use("/agents", require("./agents"));
router.use("/provision", require("./provision"));
router.use("/users", require("./users"));

module.exports = router;
