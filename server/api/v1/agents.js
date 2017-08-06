/**
 * API for managing agents
 */

const router = require('express').Router()

router.get('/', (req, res) => {
  res.status(200).json({
    msg: 'Agents API',
    methods: [
      'POST /register',
      'GET /:hostname',
      'GET /:hostname/slots'
    ]
  })
})

/**
 * Endpoint for agents to register themselves with rig's server
 */
router.post('/register', (req, res) => {
  res.status(501).json('Unimplemented')
})

/**
 * Get information for agent with given hostname
 */
router.get('/:hostname', (req, res) => {
  res.status(501).json('Unimplemented')
})

/**
 * Get slots information for given agent
 */
router.get('/:hostname/slots', (req, res) => {
  res.status(501).json('Unimplemented')
})

module.exports = router
