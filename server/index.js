/**
 * Server for xplex Internal API (nicknamed, rig)
 */

const express = require("express");
const bodyparser = require("body-parser");

// Load config
let config;

try {
  config = require("../config");
} catch (e) {
  console.error("Unable to load config.", e);
  process.exit(1);
}

// Instantiate app
let app = express();

// Enable body-parser
app.use(bodyparser.json());
app.use(bodyparser.urlencoded({ extended: true }));

// Security measures
app.disable("x-powered-by");

// Connect to DB
require("./db").connect(config.server.db);

// Mount middlewares
app.use((req, res, next) => {
  req.config = config;
  next();
});

// Mount API routes
app.use(["/latest", "/v1"], require("./api/v1"));

app.get("/", (req, res) => {
  res.status(200).json({
    msg: "xplex Internal API",
    versions: [
      "v1"
    ]
  });
});

// Start server
app.listen(config.server.port, config.server.host, () => {
  console.log(`[Server] xplex Internal API listening on ${config.server.host}:${config.server.port}`);
});
