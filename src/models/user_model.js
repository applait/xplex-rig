/**
 * User Model for DB
 */

const mongoose = require("mongoose");
const Schema = mongoose.Schema;

// Create User Schema
const UserSchema = new Schema({
  name: { type: String, default: "" },
  email: { type: String, default: "" },
  username: { type: String, default: "" },
  tokens: [{ token: String, client: String, lastUpdated: { type: Date, default: Date.now } }]
});

// Create the User model
mongoose.model("User", UserSchema);
