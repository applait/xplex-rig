/**
 * Reusable utils for xplex rig
 */

const validator = require("validator");

let sanitize = input => {
  for (let key of Object.keys(input)) {
    let i = input[key];
    if (typeof i === "string") {
      i = validator.trim(i);
      i = validator.escape(i);
      input[key] = i;
    }
  }
  return input;
};

/**
 * Middleware to sanitize inputs
 */
let sanitize_all = (req, res, next) => {
  if (req.body) req.body = sanitize(req.body);
  if (req.params) req.params = sanitize(req.params);
  if (req.query) req.query = sanitize(req.query);
  next();
};

/**
 * Middleware to ensure required fields are present
 */
let required_fields = fields => {
  return (req, res, next) => {
    req.required = req.required || {};
    for (let f of fields) {
      let val = req.body[f] || req.query[f] || req.params[f];
      if (!val) {
        res.status(400).json({
          msg: `Parameter "${f}" is required`
        });
      } else {
        req.required[f] = val;
      }
    }
    next();
  };
};

// Exports
module.exports = {
  sanitize: sanitize,
  sanitize_all: sanitize_all,
  required_fields: required_fields
};
