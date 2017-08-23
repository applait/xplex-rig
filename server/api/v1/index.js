/* global requireRelative */
/**
 * v1 API routes for xplex-internal
 */

const router = require('express').Router()
const { sanitizeAll } = requireRelative('lib/helper')

router.get('/', (req, res) => {
  res.status(200).json({
    version: 'v1',
    methods: [
      'GET /streams',
      'GET /agents',
      'GET /users'
    ]
  })
})

// Sanitize all requests
router.use(sanitizeAll)

router.use('/streams', require('./streams'))
router.use('/agents', require('./agents'))
router.use('/users', require('./users'))

module.exports = router
