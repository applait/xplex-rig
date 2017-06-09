/**
 * Agent Model for DB
 */

const mongoose = require("mongoose");
const Schema = mongoose.Schema;

// Create Agent Schema
const AgentSchema = new Schema({
  hostname: { type: String, required: true, index: true },
  capacity: { type: Number, default: 0 },
  slots_used: { type: Number, default: 0 },
  public_address: { type: String },
  private_address: { type: String },
  slots_available: { type: Boolean, default: false, index: true },
  active: { type: Boolean, default: false, index: true },
  updated_at: Date
});

AgentSchema.pre("save", function (next) {
  this.slots_available = this.capacity > this.slots_used;
  this.updated_at = new Date();
  next();
});

AgentSchema.statics.exists = function (criteria, done) {
  this.count(criteria, function (err, count) {
    if (err) {
      done(err);
      return;
    }
    if (count > 0) {
      done(null, true);
    } else {
      done(null, false);
    }
  });
};

AgentSchema.statics.available_agents = function (done) {
  return this.load({ criteria: { slots_available: true, active: true } }, done);
};

AgentSchema.statics.load = function (options, done) {
  options = options || {};
  options.select = options.select || "name capacity slots_used public_address private_address";
  return this.findOne(options.criteria)
    .select(options.select)
    .populate("tokens", "public name ot allowed_hosts")
    .exec(done);
};

// Create the Agent model
mongoose.model("Agent", AgentSchema);
