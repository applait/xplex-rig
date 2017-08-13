/* global requireRelative */
/**
 * v1 API routes for xplex-internal
 */

const router = require('express').Router()
const { sanitizeAll } = requireRelative('lib/validation')

router.get('/', (req, res) => {
  res.status(200).json({
    version: 'v1',
    methods: [
      'GET /agents',
      'GET /users'
    ]
  })
})

// Sanitize all requests
router.use(sanitizeAll)

router.use('/agents', require('./agents'))
router.use('/users', require('./users'))

module.exports = router
