/**
 * API for provisioning xplex URLs
 */

const router = require("express").Router();
const mongoose = require("mongoose");


router.get("/", (req, res) => {
  res.status(501).json("Unimplemented");
});

/**
 * Provision streaming URL for given userid.
 *
 * Look up if given user can have new URL. If yes, pick an agent, if any
 * available to provision a URL.
 *
 * @todo Implement function
 */
router.post("/:userid", (req, res) => {
  res.status(501).json("Unimplemented");
});

/**
 * Get current streaming URLs for given user
 *
 * @todo Implement function
 */
router.get("/:userid", (req, res) => {
  res.status(501).json("Unimplemented");
});

/**
 * Refresh all stream keys for given user
 *
 * @todo Implement function
 */
router.post("/:userid/refresh_keys", (req, res) => {
  res.status(501).json("Unimplemented");
});

/**
 * Refresh stream key for specified streaming URL for given user
 *
 * @todo Implement function
 */
router.post("/:userid/refresh_keys/:stream_key", (req, res) => {
  res.status(501).json("Unimplemented");
});

module.exports = router;
