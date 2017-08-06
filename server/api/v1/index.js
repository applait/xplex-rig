/**
 * v1 API routes for xplex-internal
 */

let router = require('express').Router()

router.get('/', (req, res) => {
  res.status(200).json({
    version: 'v1',
    methods: [
      'GET /agents',
      'GET /users'
    ]
  })
})

router.use('/agents', require('./agents'))
router.use('/users', require('./users'))

module.exports = router
