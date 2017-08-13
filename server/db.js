/**
 * Database wrapper. Set up connections and load models
 *
 * Call `.connect()` once and use the module for the rest of your life
 */

const debug = require('debug')('db')
const Sequelize = require('sequelize')
const { readdirSync } = require('fs')
const { join } = require('path')

let db = {
  models: {},
  sequelize: null,

  /**
   * Method to call to initiate DB connection
   *
   * @param {string} postgresURL - Full URL to PostgreSQL instance
   * @return {Sequelize} Connected Sequelize instance
   */
  connect: function (postgresURL) {
    debug(`[DB] Connecting to DB at ${postgresURL}`)
    return new Promise(function (resolve, reject) {
      this.sequelize = new Sequelize(postgresURL, {
        dialect: 'postgres',
        logging: i => debug(i)
      })

      readdirSync(join(__dirname, 'models'))
        .filter(function (file) {
          return (file.indexOf('.') !== 0) && (file !== 'index.js')
        })
        .forEach(function (file) {
          const model = this.sequelize.import(join(__dirname, 'models', file))
          this.models[model.name] = model
        }.bind(this))

      Object.keys(this.models).forEach(function (modelName) {
        if ('associate' in this.models[modelName]) {
          this.models[modelName].associate(this.models)
        }
      }.bind(this))
      resolve()
    }.bind(this))
  }
}

module.exports = db
