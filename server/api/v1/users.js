/* global requireRelative */
/**
 * Users API
 */

const router = require('express').Router()
const debug = require('debug')('api:v1:users')
const usersLib = requireRelative('lib/users')
const { requiredFields } = requireRelative('lib/validation')
const jwt = requireRelative('lib/jwt')

router.get('/', (req, res) => {
  res.status(200).json({
    msg: 'Users API',
    methods: [
      'POST /',
      'POST /password',
      'POST /auth'
    ]
  })
})

/**
 * Create new user
 */
router.post('/', requiredFields(['username', 'password', 'email']), (req, res, next) => {
  usersLib.create({
    username: req.required.username,
    email: req.required.email,
    password: req.required.password
  })
    .then(user => {
      return jwt.createUserJWT(user.get('id'), req.app.get('jwtsecret'))
    })
    .then(token => {
      res.status(200).json({
        msg: 'User created',
        status: 200,
        payload: {
          username: req.required.username,
          email: req.required.email,
          token: token
        }
      })
    })
    .catch(err => {
      debug('Unable to create user', err.message)
      let errmsg = err.message
      let errstatus = err.status || 500
      if (err.name && err.name === 'SequelizeValidationError') {
        errmsg = 'Invalid input'
        errstatus = 401
      }
      const _err = new Error(errmsg)
      _err.status = errstatus
      next(_err)
    })
})

/**
 * Update user password
 *
 * @todo Implement function
 */
router.post('/password', (req, res) => {
  res.status(501).json('Unimplemented')
})

/**
 * Attempt authenticating users and generate auth token if successful
 */
router.post('/auth', requiredFields(['username', 'password']), (req, res, next) => {
  usersLib.authenticate(req.required.username, req.required.password)
    .then(user => {
      if (user.isActive) {
        return jwt.createUserJWT(user.id, req.app.get('jwtsecret'))
      }
      const _err = new Error('User is not active')
      _err.status = 401
      return Promise.reject(_err)
    })
    .then(token => {
      res.status(200).json({
        msg: 'Authentication successful',
        status: 200,
        payload: {
          token: token
        }
      })
    })
    .catch(next)
})

module.exports = router
