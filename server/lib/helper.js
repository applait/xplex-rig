/**
 * Reusable utils for xplex rig
 */

const validator = require('validator')
const crypto = require('crypto')

const sanitize = input => {
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

const isUUIDv4 = input => validator.isUUID(input, 4)

/**
 * Middleware to sanitize inputs
 */
const sanitizeAll = (req, res, next) => {
  if (req.body) req.body = sanitize(req.body)
  if (req.params) req.params = sanitize(req.params)
  if (req.query) req.query = sanitize(req.query)
  next()
}

/**
 * Middleware to ensure required fields are present
 */
const requiredFields = fields => {
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

const sha512 = function (input, salt) {
  var hash = crypto.createHmac('sha512', salt)
  hash.update(input)
  return hash.digest('hex')
}

const sha1 = function (input, salt) {
  var hash = crypto.createHmac('sha1', salt)
  hash.update(input)
  return hash.digest('hex')
}

const generateSalt = function () {
  return crypto.randomBytes(Math.ceil(32 / 2))
    .toString('hex')
    .slice(0, 32)
}

// Exports
module.exports = {
  isUUIDv4,
  sanitize,
  sanitizeAll,
  requiredFields,
  generateSalt,
  sha512,
  sha1
}
