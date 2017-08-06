/**
 * Database wrapper
 */

const debug = require('util').debuglog('DB')

let connect = (postgresURL) => {
  debug(`[DB] Connecting to DB at ${postgresURL}`)
}

module.exports = {
  connect: connect
}
