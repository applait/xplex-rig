/**
 * Agent for xplex-rig's server
 */

const express = require("express");
const bodyparser = require("body-parser");
const os = require('os');

// Load config
let config;

try {
  config = require("../config");
} catch (e) {
  console.error("Unable to load config.", e);
  process.exit(1);
}

function get_addresses() {
  let addresses = {
    private_address: config.agent.private_address || null,
    public_address: config.agent.public_address || null
  };
  if (addresses.private_address != null && addresses.public_address != null) {
    return addresses;
  }
  let interfaces = os.networkInterfaces();
  let priv_addr_re = /^(192\.168|10\.|172\.(1[6-9]|2[0-9]|3[0-1])\.)/;
  if (addresses.private_address == null) {
    for (const item in interfaces) {
      let do_break = false;
      for (const i of interfaces[item]) {
        if (i.internal === false && i.family === 'IPv4' && priv_addr_re.test(i.address)) {
          addresses.private_address = i.address;
          do_break = true;
          break;
        }
        if (do_break) break;
      }
    }
  }
  if (addresses.public_address == null) {
    for (const item in interfaces) {
      let do_break = false;
      for (const i of interfaces[item]) {
        if (i.internal === false && i.family === 'IPv4' && !priv_addr_re.test(i.address)) {
          addresses.public_address = i.address;
          do_break = true;
          break;
        }
        if (do_break) break;
      }
    }
  }
  return addresses;
}

global.addresses = get_addresses();

// Instantiate app
let app = express();

// Enable body-parser
app.use(bodyparser.json());
app.use(bodyparser.urlencoded({ extended: true }));

// Security measures
app.disable("x-powered-by");

// Mount middlewares
app.use((req, res, next) => {
  req.config = config;
  next();
});

app.get("/", (req, res) => {
  res.status(200).json({
    msg: "xplex agent API",
    methods: [
      "status"
    ]
  });
});

app.get("/status", (req, res) => {
  res.status(200).json({
    msg: "OK",
    payload: {
      addresses: addresses,
      timestamp: new Date()
    }
  });
});

// Start server
app.listen(config.agent.port, addresses.private_address, () => {
  console.log(`[Agent] xplex agent API listening on ${addresses.private_address}:${config.agent.port}`);
});
