/**
 * Database wrapper. Set up connections and load models
 *
 * Call `.connect()` once and use the module for the rest of your life
 */

const debug = require('util').debuglog('DB')
const Sequelize = require('sequelize')
const { readdirSync } = require('fs')
const { join } = require('path')

let DB = {}

/**
 * Method to call to initiate DB connection
 *
 * @param {string} postgresURL - Full URL to PostgreSQL instance
 * @return {Sequelize} Connected Sequelize instance
 */
let connect = postgresURL => {
  debug(`[DB] Connecting to DB at ${postgresURL}`)
  let s = new Sequelize(postgresURL, {
    dialect: 'postgres',
    logging: i => debug(i)
  })

  readdirSync(join(__dirname, 'models'))
    .filter(function (file) {
      return (file.indexOf('.') !== 0) && (file !== 'index.js')
    })
    .forEach(function (file) {
      const model = s.import(join(__dirname, 'models', file))
      DB[model.name] = model
    })

  Object.keys(DB).forEach(function (modelName) {
    if ('associate' in DB[modelName]) {
      DB[modelName].associate(DB)
    }
  })

  DB.sequelize = s
  DB.Sequelize = Sequelize
}

module.exports = {
  connect: connect,
  DB: DB
}
