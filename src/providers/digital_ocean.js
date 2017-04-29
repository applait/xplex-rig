const request = require("request");

/**
 * Wrapper for Digital Ocean REST API v2
 *
 * @param {string} api_token - Digital Ocean OAuth token. See:
 * https://developers.digitalocean.com/documentation/v2/#authentication
 */
function DigitalOcean (api_token) {
  this.api_token = api_token;
  this.api_url = "https://api.digitalocean.com/v2/";
  this.request_headers = {
    authorization: `Bearer ${api_token}`,
    content_type: "application/json"
  };
}


/**
 * Send HTTP requests to Digital Ocean REST API and return a promise
 */
DigitalOcean.prototype.make_request = function (options) {
  let callback;

  let promise = new Promise((resolve, reject) => {
    callback = (err, res, body) => {
      if (err) {
        reject(err);
      } else {
        resolve({ res, body });
      }
    }
  });

  let opts = {
    uri: `${this.api_url}${options.action}`,
    method: options.method || "GET",
    headers: options.headers || this.request_headers,
    body: options.body || {},
    strictSSL: true,
    json: true,
    qs: options.qs || {}
  };

  request(opts, (err, res, body) => {
    if (err) {
      callback(err);
    } else if (!err && !/^[2][0-9][0-9]$/.test(res.statusCode)) {
      callback({ res, body });
    } else {
      callback(null, res, body);
    }
  });

  return promise;
};


/**
 * Get Digital Ocean user account information
 */
DigitalOcean.prototype.account = function () {
  return this.make_request({ action: "account" });
};


/**
 * Get information about specific droplet
 *
 * @param {string|number} id - The droplet ID
 */
DigitalOcean.prototype.droplet = function (id, sub_action="", page=1, limit=200) {
  return this.make_request({
    action: `droplets/${id}${sub_action ? "/" + sub_action : ""}`,
    qs: {
      page: page,
      per_page: limit
    }
  });
};


/**
 * Get information about multiple droplets
 */
DigitalOcean.prototype.droplets = function (tag_name=null, page=1, limit=200) {
  return this.make_request({
    action: "droplets",
    qs: {
      page: page,
      per_page: limit,
      tag_name: tag_name
    }
  });
};


/**
 * Create one or multiple new droplets
 *
 * @param {object} droplet_options - Request body attributes set according to
 * the specification described at
 * https://developers.digitalocean.com/documentation/v2/#create-a-new-droplet
 */
DigitalOcean.prototype.droplets_create = function (droplet_options) {
  return this.make_request({
    action: "droplets",
    method: "POST",
    body: droplet_options
  });
};

/**
 * Perform action on a single droplet
 *
 * For available actions and their options, see
 * https://developers.digitalocean.com/documentation/v2/#droplet-actions.
 *
 * @param {string|number} id - The droplet ID
 * @param {object} action_options - The object to be passed as request body
 */
DigitalOcean.prototype.droplet_action = function (id, action_options) {
  return this.make_request({
    action: `droplets/${id}/actions`,
    method: "POST",
    body: action_options
  });
};


/**
 * Delete a single droplet
 *
 * @param {string|number} id - The droplet ID
 */
DigitalOcean.prototype.droplet_delete = function (id) {
  return this.make_request({
    action: `droplets/${id}`,
    method: "DELETE"
  });
};


/**
 * Delete all droplets associated with a tag
 */
DigitalOcean.prototype.droplets_delete_by_tag = function (tag) {
  return this.make_request({
    action: `droplets`,
    method: "DELETE",
    qs: {
      tag_name: tag
    }
  });
};


module.exports = DigitalOcean;
