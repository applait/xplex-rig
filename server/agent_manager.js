/**
 * Agent manager - Creates and manages agents, including healthchecks
 */

const request = require('request');
const {promisify} = require('util');

function make_agent_request(agent, options) {
  let r = promisify(request);
  let opts = {
    uri: `http://${agent.private_address}/${options.action}`,
    method: options.method || 'GET',
    headers: options.headers || {},
    body: options.body || {},
    strictSSL: false,
    json: true,
    qs: options.qs || {}
  };
  return r(opts);
}

async function healthcheck(agent) {
  make_agent_request(agent, { action: 'status' })
    .then(function (res) {
      if (res.statusCode == 200) {
        agent.active = true;
      } else {
        agent.active = false;
      }
      agent.updated_at = new Date();
      agent.save();
    }).catch(function (err) {
      agent.active = false;
      agent.updated_at = new Date();
      agent.save();
    });
}

async function healthcheck_poll(agent) {
  setInterval(function () {
    healthcheck(agent);
  }, 3000);
}

async function create_agent_digital_ocean (DO, hostname="live", region="sgp1") {
  return DO.droplets_create({
    "name": `${hostname}-${region}.xplex.me`,
    "region": region,
    "size": "512mb",
    "image": "fedora-25-x64",
    "ssh_keys": null, // @TODO Add SSH keys
    "backups": false,
    "ipv6": true,
    "user_data": null,
    "private_networking": true,
    "volumes": null,
    "tags": [
      "xplex-agent",
      "xplex"
    ]
  })
  .then(function(do_res, do_body) {
    console.log("[Digital Ocean] Created droplet", do_body);
    return do_body;
  }, function (err) {
     console.log("[Digital Ocean] Error creating droplet", err);
     return err;
  });
}

module.exports = {
  healthcheck_poll: healthcheck_poll,
  create_agent_digital_ocean: create_agent_digital_ocean
}
