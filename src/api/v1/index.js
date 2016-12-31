/**
 * v1 API routes for xplex-internal
 */

var router = require("express").Router();

router.get("/", (req, res) => {
  res.status(200).json({
    version: "v1",
    methods: [
      "/user",
      "/provision"
    ]
  });
});

module.exports = router;
