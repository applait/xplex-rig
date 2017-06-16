/**
 * Users API
 */

const router = require("express").Router();
const mongoose = require("mongoose");
const User = mongoose.model("User");

router.get("/", (req, res) => {
  res.status(200).json({
    msg: "Users API",
    methods: [
      "POST /new",
      "POST /:id",
      "POST /auth",
      "GET /validate_token"
    ]
  });
});

/**
 * Create new user
 */
router.post("/new", (req, res) => {
  // TODO: handle validations

  let newuser = new User(req.body);
  newuser.provider = "local";

  try {
    newuser.save();
  } catch (e) {
    console.error("[API] Error creating new user", e);
    res.status(500).json({
      msg: "Error creating new user",
      status: 500
    });
    return;
  }

  res.status(200).json({
    msg: "User created"
  });

});

/**
 * Update user information
 *
 * @todo Implement function
 */
router.post("/:id", (req, res) => {
  res.status(501).json("Unimplemented");
});

/**
 * Attempt authenticating users and generate auth token if successful
 *
 * @todo Implement function
 */
router.post("/auth", (req, res) => {
  res.status(501).json("Unimplemented");
});

/**
 * Validate auth token for users
 *
 * @todo Implement function
 */
router.get("/validate_token", (req, res) => {
  res.status(501).json("Unimplemented");
});

module.exports = router;
