/**
 * API for managing agents
 */

const router = require("express").Router();
const mongoose = require("mongoose");

const Agent = mongoose.model("Agent");
const utils = require("../../utils");
const agent_manager = require('../../agent_manager');

router.get("/", (req, res) => {
  res.status(200).json({
    msg: "Agents API",
    methods: [
      "POST /register",
      "GET /:hostname",
      "POST /:hostname/activate",
      "POST /:hostname/deactivate",
      "GET /:hostname/slots",
      "POST /:hostname/slots"
    ]
  });
});

let load_agent = (req, res, next) => {
  let opts = { criteria: { hostname: req.params.hostname },
               select: "hostname public_address capacity slots_used active" };

  Agent.load(opts, function (err, agent) {
    if (err) {
      res.status(500).json({ msg: "Error fetching agent. Try again later." });
      return;
    }
    if (agent == null) {
      res.status(404).json({ msg: "Agent not found." });
      return;
    }
    req.agent = agent;
    next();
  });
};

let fields = ["hostname", "public_address", "private_address", "capacity"];
let exists = (req, res, next) => {
  Agent.exists({ hostname: req.required.hostname }, function (err, exists) {
    if (err) {
      res.status(500).json({ msg: "Error creating new agent. Try again later." });
      return;
    }
    if (exists) {
      res.status(400).json({ msg: `Hostname already exists. Try another.` });
      return;
    }
    next();
  });
};

/**
 * Endpoint for agents to register themselves with rig's server
 */
router.post("/register", utils.sanitize_all, utils.required_fields(fields), exists, (req, res) => {
  let newagent = new Agent({
    hostname: req.required.hostname,
    public_address: req.required.public_address,
    private_address: req.required.private_address,
    capacity: req.required.capacity,
    slots: 0,
    active: true
  });

  try {
    newagent.save();
  } catch (e) {
    console.log(e);
    res.status(500).json({ msg: "Error validating new agent. Try again later." });
  }

  agent_manager.healthcheck_poll(newagent);

  res.status(200).json({
    msg: "Agent created",
    payload: {
      hostname: newagent.hostname,
      public_address: newagent.public_address,
      private_address: newagent.private_address
    }
  });
});

/**
 * Get information for agent with given hostname
 */
router.get("/:hostname", utils.sanitize_all, load_agent, (req, res) => {
  res.status(200).json({
    msg: "Agent information",
    payload: {
      hostname: req.agent.hostname,
      public_address: req.agent.public_address,
      private_address: req.agent.private_address,
      capacity: req.agent.capacity,
      slots_used: req.agent.slots_used,
      active: req.agent.active
    }
  });
});

/**
 * Mark specified agent as active
 */
router.post("/:hostname/activate", utils.sanitize_all, load_agent, (req, res) => {
  req.agent.active = true;
  req.agent.updated_at = new Date();
  try {
    req.agent.save();
    res.status(200).json({
      msg: "Agent activated",
      payload: {
        hostname: req.agent.hostname,
        active: req.agent.active
      }
    });
  } catch (err) {
      res.status(500).json({ msg: "Error activating agent. Try again later." });
  }
});

/**
 * Mark specified agent as inactive
 */
router.post("/:hostname/deactivate", utils.sanitize_all, load_agent, (req, res) => {
  req.agent.active = false;
  req.agent.updated_at = new Date();
  try {
    req.agent.save();
    res.status(200).json({
      msg: "Agent deactivated",
      payload: {
        hostname: req.agent.hostname,
        active: req.agent.active
      }
    });
  } catch (err) {
      res.status(500).json({ msg: "Error deactivating agent. Try again later." });
  }
});

/**
 * Get slots information for given agent
 */
router.get("/:hostname/slots", (req, res) => {
  res.status(501).json("Unimplemented");
});

/**
 * Update slots information for given agent
 *
 * @todo Implement function
 */
router.post("/:hostname/slots", (req, res) => {
  res.status(501).json("Unimplemented");
});

module.exports = router;
