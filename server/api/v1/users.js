/**
 * Users API
 */

const router = require('express').Router()

router.get('/', (req, res) => {
  res.status(200).json({
    msg: 'Users API',
    methods: [
      'POST /new',
      'POST /:id',
      'POST /auth',
      'GET /validate_token'
    ]
  })
})

/**
 * Create new user
 */
router.post('/new', (req, res) => {
  // TODO: handle validations
  res.status(501).json('Unimplemented')
})

/**
 * Update user information
 *
 * @todo Implement function
 */
router.post('/:id', (req, res) => {
  res.status(501).json('Unimplemented')
})

/**
 * Attempt authenticating users and generate auth token if successful
 *
 * @todo Implement function
 */
router.post('/auth', (req, res) => {
  res.status(501).json('Unimplemented')
})

/**
 * Validate auth token for users
 *
 * @todo Implement function
 */
router.get('/validate_token', (req, res) => {
  res.status(501).json('Unimplemented')
})

module.exports = router
