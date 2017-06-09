/**
 * API for provisioning xplex URLs
 */

const router = require("express").Router();
const mongoose = require("mongoose");

router.get("/", (req, res) => {
  res.status(200).json("Under construction");
});

module.exports = router;
