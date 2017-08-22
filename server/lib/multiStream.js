/* global requireRelative */
/**
 * Multi-stream config and dispatch handler
 */

const { models } = requireRelative('db')
const debug = require('debug')('lib:multistream')
const uuid = require('uuid')
const services = require('./multiStreamServices')

/**
 * Create a new multi-stream entry for given user
 *
 * @param {string} userID - Create
 * @return {Promise} Resolves to new `MultiStream` instance if successful
 */
function createMultiStream (userID) {
  return models.MultiStream.create({
    key: uuid.v4(),
    User: userID,
    isActive: true
  })
    .then(ms => {
      debug(`Created new multi stream for user: ${userID}, streamID: ${ms.id}`)
      return ms
    })
    .catch(err => {
      debug(`Error creating new multistream for user ${userID}`)
      return err
    })
}

/**
 * Add multi-streaming config for specific multi-stream
 *
 * @param {string} streamID - MultiStream ID to add stream config for
 * @param {string} service - A service name as defined in `multiStreamServices` module
 * @param {string} key - Secret key for the service for the user
 * @param {string} [server='default'] - A specific server to route traffic to for given `service`
 * @return {Promise} Resolves to `MultiStreamConfig` instance if successful
 */
function addMultiStreamConfig (streamID, service, key, server = 'default') {
  return models.MultiStreamConfig.create({
    service: service,
    server: server,
    key: key,
    isActive: true,
    MultiStream: streamID
  })
    .then(msconfig => {
      debug(`Created new multistream config. Stream: ${streamID}, service: ${service}`)
      return msconfig
    })
    .catch(err => {
      debug(`Error creating multistream config. Stream: ${streamID}, service: ${service}`)
      return err
    })
}

/**
 * Get current configurations for given stream key
 *
 * @param {string} streamKey - xplex multi-stream key
 * @param {boolean} [isActive=true] - If `true` return only active keys
 * @param {boolean} [isStreaming=false] - If `false` return only streams that are not marked to be currently streaming
 */
function getMultiStreamConfigs (streamKey, isActive = true, isStreaming = false) {
  return models.MultiStream.find({
    where: { key: streamKey, isActive: isActive, isStreaming: isStreaming },
    include: [{
      model: models.MultiStreamConfig,
      where: { isActive: isActive }
    }]
  })
    .then(ms => {
      if (ms === null) {
        return Promise.reject(new Error('No valid stream exists for given stream key'))
      }
      const configs = ms.getMultiStreamConfigs({ plain: true })
      const output = {
        streamID: ms.id,
        streamKey: ms.key,
        destinations: []
      }
      for (let c of configs) {
        const serviceURL = createRTMPURL(c.service, c.key, c.server)
        if (serviceURL === null) {
          continue
        }
        output.destinations.push({
          service: c.service,
          url: serviceURL
        })
      }
      return Promise.resolve(output)
    })
    .catch(err => {
      debug(`Error retrieving stream configs for key: ${streamKey}. Reason: ${err.message}`)
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
  createMultiStream,
  addMultiStreamConfig,
  getMultiStreamConfigs,
  createRTMPURL
}
