/* global requireRelative */
/**
 * Multi-stream config and dispatch handler
 */

const { models } = requireRelative('db')
const debug = require('debug')('lib:multistream')
const services = require('./multiStreamServices')
const { sha1 } = requireRelative('lib/helper')

/**
 * Create a new multi-stream entry for given user
 *
 * @param {string} userID - Create
 * @return {Promise} Resolves to new `MultiStream` instance if successful
 */
function createMultiStream (userID) {
  return models.User.find({
    where: { id: userID, isActive: true },
    attributes: ['salt']
  })
    .then(user => {
      if (user === null) {
        const _err = new Error('Invalid user id')
        _err.status = 401
        return Promise.reject(_err)
      }
      return models.MultiStream.create({
        key: sha1(`xplex://${userID}@${Date.now()}`, user.salt),
        UserId: userID,
        isActive: true
      })
    })
    .then(ms => {
      debug(`Created new multi stream for user: ${userID}, streamID: ${ms.id}`)
      return ms
    })
    .catch(err => {
      debug(`Error creating new multistream for user ${userID}`)
      return Promise.reject(err)
    })
}

/**
 * Update streaming key of given multi-stream entry
 */
function updateMultiStreamKey (userID, streamID) {
  return models.MultiStream.find({
    where: {
      id: streamID,
      UserId: userID
    },
    include: [{
      model: models.User,
      attributes: ['salt']
    }]
  })
    .then(ms => {
      if (ms === null) {
        const _err = new Error('Invalid user ID or stream ID specified')
        _err.status = 403
        return Promise.reject(_err)
      }
      ms.key = sha1(`xplex://${ms.User.id}@${Date.now()}`, ms.User.salt)
      return ms.save({ fields: ['key'] })
        .then(() => Promise.resolve(ms.key))
    })
    .catch(err => {
      debug(`Error updating streaming key for multistream ID ${streamID}`)
      return Promise.reject(err)
    })
}

/**
 * Add multi-streaming config for specific multi-stream
 *
 * @param {string} userID - UserID to check if stream belongs to right user
 * @param {string} streamID - MultiStream ID to add stream config for
 * @param {string} service - A service name as defined in `multiStreamServices` module
 * @param {string} key - Secret key for the service for the user
 * @param {string} [server='default'] - A specific server to route traffic to for given `service`
 * @return {Promise} Resolves to `MultiStreamConfig` instance if successful
 */
function addMultiStreamConfig (userID, streamID, service, key, server = 'default') {
  return models.MultiStream.find({
    where: {
      id: streamID,
      UserId: userID,
      isActive: true
    }
  })
    .then(ms => {
      if (ms === null) {
        const _err = new Error('Invalid user ID or stream ID specified')
        _err.status = 403
        return Promise.reject(_err)
      }
      return models.MultiStreamConfig.create({
        service: service,
        server: server,
        key: key,
        isActive: true,
        MultiStreamId: streamID
      })
    })
    .then(msconfig => {
      debug(`Created new multistream config. Stream: ${streamID}, service: ${service}`)
      return Promise.resolve(msconfig)
    })
    .catch(err => {
      debug(err)
      debug(`Error creating multistream config. Stream: ${streamID}, service: ${service}`)
      return Promise.reject(err)
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
    attributes: ['id', 'key', 'isActive'],
    include: [{
      model: models.MultiStreamConfig,
      attributes: ['id', 'service', 'key', 'server'],
      required: false,
      where: { isActive: isActive }
    }]
  })
    .then(msInstance => {
      if (msInstance === null) {
        const _err = new Error('No valid stream exists for given stream key')
        _err.status = 404
        return Promise.reject(_err)
      }
      const ms = msInstance.get({ plain: true })
      const output = {
        streamID: ms.id,
        streamKey: ms.key,
        destinations: []
      }
      for (let c of ms.MultiStreamConfigs) {
        const serviceURL = getRTMPURL(c.service, c.key, c.server)
        if (serviceURL === null) {
          continue
        }
        output.destinations.push({
          configID: c.id,
          service: c.service,
          url: serviceURL
        })
      }
      return Promise.resolve(output)
    })
    .catch(err => {
      debug(`Error retrieving stream configs for key: ${streamKey}. Reason: ${err.message}`)
      return Promise.reject(err)
    })
}

/**
 *
 * @param {string} service - A service name as defined in `multiStreamServices` module
 * @param {string} key - Stream key provided by the service
 * @param {string} [server='default'] - Optional server name to stream to for given service
 */
function getRTMPURL (service, key, server = 'default') {
  if (!services[service]) {
    return null
  }
  const s = services[service]
  const serverInfo = s.servers.get(server)
  if (!serverInfo) {
    return null
  }
  return `${serverInfo.url}/${key}`
}

module.exports = {
  createMultiStream,
  updateMultiStreamKey,
  addMultiStreamConfig,
  getMultiStreamConfigs,
  getRTMPURL
}
