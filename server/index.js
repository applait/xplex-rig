/**
 * Server for xplex Internal API (nicknamed, rig)
 */

const express = require('express')
const bodyparser = require('body-parser')

// Load config
let config

try {
  config = require('../config')
} catch (e) {
  console.error('Unable to load config.', e)
  process.exit(1)
}

// Instantiate app
let app = express()

// Set JWT secret
app.set('jwtsecret', process.env.JWTSECRET || config.server.jwtsecret)

// Enable body-parser
app.use(bodyparser.json())
app.use(bodyparser.urlencoded({ extended: true }))

// Security measures
app.disable('x-powered-by')

// Connect to DB
require('./db').connect(config.server.postgres_url)

// Mount middlewares
app.use((req, res, next) => {
  req.config = config
  next()
})

// Mount API routes
app.use(['/latest', '/v1'], require('./api/v1'))

app.get('/', (req, res) => {
  res.status(200).json({
    msg: 'xplex Internal API',
    versions: [
      'v1'
    ]
  })
})

// catch 404 and forward to error handler
app.use(function (req, res, next) {
  var err = new Error('Not Found')
  err.status = 404
  next(err)
})

// error handler
// no stacktraces leaked to user unless in development environment
app.use(function (err, req, res, next) {
  err.status = err.status || 500
  res.status(err.status).json({
    msg: err.message,
    status: err.status
  })
})

// Start server
app.listen(config.server.port, config.server.host, () => {
  console.log(`[Server] xplex Internal API listening on ${config.server.host}:${config.server.port}`)
})
