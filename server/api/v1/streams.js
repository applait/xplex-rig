/* global requireRelative */
/**
 * API for creating, updating and retrieving streams information
 */

const router = require('express').Router()
const { requiredFields } = requireRelative('lib/helper')
const multiStream = requireRelative('lib/multiStream')
const jwt = requireRelative('lib/jwt')

router.get('/', (req, res) => {
  res.status(200).json({
    msg: 'Streams API',
    methods: [
      'POST /create',
      'POST /config/new',
      'GET /config'
    ]
  })
})

/**
 * Create new multi-stream entry for given user
 */
router.post('/create', jwt.verifyUser, (req, res, next) => {
  multiStream.createMultiStream(req.user.id)
    .then(ms => {
      res.status(200).json({
        msg: 'Stream created',
        payload: {
          streamID: ms.id,
          streamKey: ms.key,
          active: ms.active,
          user: req.user.id
        }
      })
    })
    .catch(next)
})

/**
 * Retrieve multi-stream configs for given `streamKey`
 */
router.get('/config', requiredFields(['streamKey']), jwt.verifyUser, (req, res, next) => {
  multiStream.getMultiStreamConfigs(req.required.streamKey, true, false)
    .then(c => {
      res.status(200).json({
        msg: 'Stream config',
        payload: c
      })
    })
    .catch(next)
})

/**
 * Add a service-provider configuration entry for a given multi-stream URL
 */
router.post('/config/new', requiredFields(['streamID', 'service', 'key', 'server']), jwt.verifyUser, (req, res, next) => {
  multiStream.addMultiStreamConfig(
    req.user.id,
    req.required.streamID,
    req.required.service,
    req.required.key,
    req.required.server)
    .then(msconfig => {
      res.status(200).json({
        msg: 'Stream config added',
        payload: {
          streamID: req.required.streamID,
          service: req.required.service,
          server: req.required.server
        }
      })
    })
    .catch(next)
})

module.exports = router
