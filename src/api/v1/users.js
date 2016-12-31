/**
 * Users API
 */

const router = require("express").Router();
const mongoose = require("mongoose");
const User = mongoose.model("User");

router.get("/", (req, res) => {
  res.status(200).json({
    methods: [
      "POST /new"
    ]
  });
});

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

module.exports = router;
