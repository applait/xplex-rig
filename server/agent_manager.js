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

module.exports = {
  healthcheck_poll: healthcheck_poll
}
