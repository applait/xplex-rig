/* global requireRelative */
/**
 * User library for rig-server. Includes methods for creating and editing users
 */

const { models } = requireRelative('db')
const debug = require('debug')('lib:users')

let Users = {

  create: function (userdata, notify = true, confirm = true) {
    return models.User.create({
      username: userdata.username,
      email: userdata.email,
      password: userdata.password,
      // Currently mark user as active by default
      // TODO: Change this behaviour once notifications are enabled
      isActive: true
    })
      .then(user => {
        debug(`Created user. ID: ${user.get('id')}, Username: ${user.get('username')}`)
        return user
      })
  },

  update: function (userId, userdata) {
    return models.User.update(userdata, { where: { id: userId } })
      .then(count => count > 0)
  },

  authenticate: function (username, password) {
    return models.User.findOne({ where: { username: username } })
      .then(user => {
        if (user === null || !user.checkPassword(password)) {
          const _err = new Error('Invalid username or password')
          _err.status = 401
          return Promise.reject(_err)
        }
        return user.toJSON()
      })
  },

  changePassword: function (userId, newPassword) {
    models.User.findById(userId)
      .then(user => {
        if (user === null) {
          const _err = new Error('Invalid user ID')
          _err.status = 404
          return Promise.reject(_err)
        }
        return this.update(userId, { password: newPassword })
      })
  }
}

module.exports = Users
