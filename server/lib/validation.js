/**
 * Reusable utils for xplex rig
 */

const validator = require('validator')

let sanitize = input => {
  for (let key of Object.keys(input)) {
    let i = input[key]
    if (typeof i === 'string') {
      i = validator.trim(i)
      i = validator.escape(i)
      input[key] = i
    }
  }
  return input
}

/**
 * Middleware to sanitize inputs
 */
let sanitizeAll = (req, res, next) => {
  if (req.body) req.body = sanitize(req.body)
  if (req.params) req.params = sanitize(req.params)
  if (req.query) req.query = sanitize(req.query)
  next()
}

/**
 * Middleware to ensure required fields are present
 */
let requiredFields = fields => {
  return (req, res, next) => {
    req.required = req.required || {}
    for (let f of fields) {
      let val = req.body[f] || req.query[f] || req.params[f]
      if (!val) {
        const _err = new Error(`Parameter "${f}" is required`)
        _err.status = 400
        return next(_err)
      } else {
        req.required[f] = val
      }
    }
    next()
  }
}

// Exports
module.exports = {
  sanitize,
  sanitizeAll,
  requiredFields
}
