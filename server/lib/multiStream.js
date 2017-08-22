/* global requireRelative */
/**
 * Multi-stream config and dispatch handler
 */

const { models } = requireRelative('db')
const debug = require('debug')('lib:multistream')
const services = require('./multiStreamServices')

/**
 *
 * @param {string} userID - User's ID to add multi-stream config for
 * @param {string} service - A service name as defined in `multiStreamServices` module
 * @param {string} key - Secret key for the service for the user
 * @param {string} [server='default'] - A specific server to route traffic to for given `service`
 */
function addConfig (userID, service, key, server = 'default') {
  return models.MultiStreamConfig.create({
    service: service,
    server: server,
    key: key,
    user: userID,
    isActive: true
  })
    .then(msconfig => {
      debug(`Created new multistream config. User: ${userID}, service: ${service}`)
      return msconfig
    })
    .catch(err => {
      debug(`Error creating multistream config. User: ${userID}, service: ${service}`)
      return err
    })
}

/**
 *
 * @param {string} service - A service name as defined in `multiStreamServices` module
 * @param {string} key - Stream key provided by the service
 * @param {string} [server='default'] - Optional server name to stream to for given service
 */
function createRTMPURL (service, key, server = 'default') {
  if (!services[service]) {
    return null
  }
  const s = services[service]
  const serverURL = s.servers.get(server)
  if (!serverURL) {
    return null
  }
  return `${serverURL}/${key}`
}

module.exports = {
  addConfig,
  createRTMPURL
}
