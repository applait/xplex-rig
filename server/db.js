/**
 * Database wrapper. Set up connections and load models
 *
 * Call `.connect()` once and use the module for the rest of your life
 */

const debug = require('util').debuglog('DB')
const Sequelize = require('sequelize')
const { readdirSync } = require('fs')
const { join } = require('path')

let models = {}
let sequelize

/**
 * Method to call to initiate DB connection
 *
 * @param {string} postgresURL - Full URL to PostgreSQL instance
 * @return {Sequelize} Connected Sequelize instance
 */
let connect = postgresURL => {
  debug(`[DB] Connecting to DB at ${postgresURL}`)
  sequelize = new Sequelize(postgresURL, {
    dialect: 'postgres',
    logging: i => debug(i)
  })

  readdirSync(join(__dirname, 'models'))
    .filter(function (file) {
      return (file.indexOf('.') !== 0) && (file !== 'index.js')
    })
    .forEach(function (file) {
      const model = sequelize.import(join(__dirname, 'models', file))
      models[model.name] = model
    })

  Object.keys(models).forEach(function (modelName) {
    if ('associate' in models[modelName]) {
      models[modelName].associate(models)
    }
  })
}

module.exports = {
  connect: connect,
  models: models,
  sequelize: sequelize
}
