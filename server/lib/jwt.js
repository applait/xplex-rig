/**
 * Library for handling JWT based authentications
 */

const jwt = require('jsonwebtoken')
const debug = require('debug')('lib:jwt')
const { isUUIDv4 } = require('./helper')

/**
 * Extract token from header
 *
 * @param {object} req - ExpressJS `request` object. Token is extracted from `Authorization` header with the scheme
 * `Bearer`
 * @return {string|null} Returns string containing token if found, else returns `null`
 */
function extractAuthToken (req) {
  if (req.headers.authorization && req.headers.authorization.split(' ')[0] === 'Bearer') {
    return req.headers.authorization.split(' ')[1]
  }
  return null
}

/**
 * Verify JWT by extracting token from request header
 *
 * @param {object} req - ExpressJS `request` object. Token is extracted from `Authorization` header with the scheme
 * `Bearer`
 * @param {string} secret - Secret used to sign the token
 * @param {object} jwtopts - Options to pass to `jsonwebtoken.verify`
 * @return {Promise} Return promise that resolves to decoded and verified payload
 */
function jwtVerify (req, secret, jwtopts) {
  return new Promise(function (resolve, reject) {
    const token = extractAuthToken(req)
    if (!token) {
      const _err = new Error('Bearer token not specified in Authorization header')
      _err.status = 400
      return reject(_err)
    }
    jwt.verify(token, secret, jwtopts, function (err, payload) {
      if (err) {
        debug('Error decoding token', err.message)
      }
      if (err || !payload.iss || !payload.ist || ['user', 'client'].indexOf(payload.ist) < 0) {
        const _err = new Error('Invalid token')
        _err.status = 400
        return reject(_err)
      }
      return resolve(payload)
    })
  })
};

/**
 * ExpressJS middleware to verify user JWT claims
 *
 * Token secret is retrieved from `jwtsecret` property in `app` settings.
 *
 * If verfied successfully, `req.user.id` contains decoded user ID and `req.user.type` is set to `user`. Plus, `req.jwt`
 * is set to decoded value of the token's `payload`.
 *
 * @param {object} req - ExpressJS `request` object
 * @param {object} res - ExpressJS `response` object
 * @param {function} next - ExpressJS middleware `next` function
 */
function verifyUser (req, res, next) {
  const jwtopts = {
    audience: ['rig.xplex.me'],
    algorithms: ['HS512']
  }
  jwtVerify(req, req.app.get('jwtsecret'), jwtopts)
    .then(payload => {
      if (payload.ist !== 'user' || !isUUIDv4(payload.iss)) {
        const _err = new Error('Invalid user authentication token')
        _err.status = 401
        return Promise.reject(_err)
      }
      req.user = req.user || {}
      req.user.id = payload.iss
      req.user.type = payload.ist
      req.jwt = payload
      next()
    })
    .catch(next)
}

/**
 * Create a JWT for users
 *
 * @param {string|number} userId - User ID for which to generate token
 * @param {string} secret - Secret used to sign the token
 * @return {Promise} Promise that resolves to generated JWT
 */
function createUserJWT (userId, secret) {
  return new Promise(function (resolve, reject) {
    const tokenInput = {
      iss: userId,
      ist: 'user'
    }
    const jwtopts = {
      algorithm: 'HS512',
      expiresIn: '28d',
      audience: 'rig.xplex.me'
    }
    jwt.sign(tokenInput, secret, jwtopts, function (err, token) {
      if (err) {
        debug('Error creating user token', err.message)
        const _err = new Error('Error generating user token')
        _err.status = 500
        return reject(err)
      }
      resolve(token)
    })
  })
}

/**
 * Create a JWT for invites
 *
 * @param {string|number} userId - User ID who generates the invite
 * @param {string} email - Email for which to generate token
 * @param {string} secret - Secret used to sign the token
 * @return {Promise} Promise that resolves to generated JWT
 */
function createInviteJWT (userId, email, secret) {
  return new Promise(function (resolve, reject) {
    const tokenInput = {
      iss: userId,
      sub: email,
      ist: 'invite'
    }
    const jwtopts = {
      algorithm: 'HS512',
      expiresIn: '14d',
      audience: 'rig.xplex.me'
    }
    jwt.sign(tokenInput, secret, jwtopts, function (err, token) {
      if (err) {
        debug('Error creating invite token', err.message)
        const _err = new Error('Error generating invite token')
        _err.status = 500
        return reject(err)
      }
      resolve(token)
    })
  })
}

/**
 * Verifies JWT signature for invites
 *
 * @param {string} email - Email to be verified
 * @param {string} token - Token whose signature to be verified
 * @return {Promise} Promise that resolves to generated JWT
 */
function verifyInvite (email, tokenInput, secret) {
  return new Promise(function (resolve, reject) {
    const jwtopts = {
      algorithms: 'HS512',
      audience: 'rig.xplex.me',
      sub: email // Verifies if subject matches
    }
    jwt.verify(tokenInput, secret, jwtopts, function (err, token) {
      if (err) {
        debug('Error verifying invite token', err.message)
        const _err = new Error('Error verifying invite token')
        _err.status = 500
        return reject(err)
      }
      resolve(email)
    })
  })
}

module.exports = {
  verifyUser,
  createUserJWT,
  createInviteJWT,
  verifyInvite
}
