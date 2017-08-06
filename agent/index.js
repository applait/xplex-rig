/**
 * Agent for xplex-rig's server
 */

const express = require('express')
const bodyparser = require('body-parser')
const os = require('os')

// Load config
let config

try {
  config = require('../config')
} catch (e) {
  console.error('Unable to load config.', e)
  process.exit(1)
}

function getAddresses () {
  let addresses = {
    private_address: config.agent.private_address || null,
    public_address: config.agent.public_address || null
  }
  if (addresses.private_address != null && addresses.public_address != null) {
    return addresses
  }
  let interfaces = os.networkInterfaces()
  let privAddrRegex = /^(192\.168|10\.|172\.(1[6-9]|2[0-9]|3[0-1])\.)/
  if (addresses.private_address == null) {
    for (const item in interfaces) {
      let doBreak = false
      for (const i of interfaces[item]) {
        if (i.internal === false && i.family === 'IPv4' && privAddrRegex.test(i.address)) {
          addresses.private_address = i.address
          doBreak = true
          break
        }
        if (doBreak) break
      }
    }
  }
  if (addresses.public_address == null) {
    for (const item in interfaces) {
      let doBreak = false
      for (const i of interfaces[item]) {
        if (i.internal === false && i.family === 'IPv4' && !privAddrRegex.test(i.address)) {
          addresses.public_address = i.address
          doBreak = true
          break
        }
        if (doBreak) break
      }
    }
  }
  return addresses
}

global.addresses = getAddresses()

// Instantiate app
let app = express()

// Enable body-parser
app.use(bodyparser.json())
app.use(bodyparser.urlencoded({ extended: true }))

// Security measures
app.disable('x-powered-by')

// Mount middlewares
app.use((req, res, next) => {
  req.config = config
  next()
})

app.get('/heartbeat', (req, res) => {
  res.status(200).send()
})

app.get('/status', (req, res) => {
  res.status(200).json({
    msg: 'OK',
    payload: {
      addresses: global.addresses,
      timestamp: new Date()
    }
  })
})

app.get('/', (req, res) => {
  res.status(200).json({
    msg: 'xplex agent API',
    methods: [
      'status'
    ]
  })
})

// Start server
app.listen(config.agent.port, global.addresses.private_address, () => {
  console.log(`[Agent] xplex agent API listening on ${global.addresses.private_address}:${config.agent.port}`)
})
