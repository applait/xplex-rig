/**
 * Database wrapper
 */

const db = require("mongoose");
const readdir = require("fs").readdirSync;
const pathjoin = require("path").join;

let connect = (dbconfig) => {
  console.log(`[DB] Connecting to DB at ${dbconfig.mongo_url}`);

  // Load models from `models/*_model.js`
  let modeldir = pathjoin(__dirname, "models");
  readdir(modeldir)
    .filter(file => ~file.search(/^[^\.].*\_model\.js$/))
    .forEach(file => {
      console.log(`[DB] Loading model: ${file}`);
      require(pathjoin(modeldir, file));
    });

  // Create connect function
  let dbconn = () => {
    let conn = db.connect(dbconfig.mongo_url, { server: { socketOptions: { keepAlive: 1 }}});
    let timeout;

    conn.connection.removeAllListeners();

    conn.connection
      .once("error", () => {
        console.error(`[DB] Unable to connect to DB at ${dbconfig.mongo_url}`);
      })
      .once("disconnected", () => {
        console.warn(`[DB] Disconnected from DB at ${dbconfig.mongo_url}`);
        if (timeout) {
          clearTimeout(timeout);
        }
        timeout = setTimeout(() => {
          dbconn();
        }, 5000);
      })
      .once("open", () => {
        if (timeout) {
          clearTimeout(timeout);
        }
        console.log(`[DB] Connected to DB at ${dbconfig.mongo_url}`);
      });
  };

  dbconn();
};

module.exports = {
  connect: connect
};
