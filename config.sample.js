/**
 * Config sample
 */

module.exports = {

  /**
   * Server specific config
   */
  server: {
    host: "0.0.0.0",
    port: "8081"
  },

  /**
   * DB config
   */
  db: {
    mongo_url: "mongodb://127.0.0.1:27017/xplex"
  },

  /**
   * Digital Ocean config
   *
   * See: https://developers.digitalocean.com/documentation/v2/#authentication
   */
  digital_ocean: {
    token: "xxxx"
  }
};
