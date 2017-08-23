/* global requireRelative */
/**
 * Database model for User
 */

const { generateSalt, sha512 } = requireRelative('lib/helper')

module.exports = function (sequelize, DataTypes) {
  var User = sequelize.define('User', {
    id: {
      type: DataTypes.UUID,
      defaultValue: DataTypes.UUIDV4,
      primaryKey: true,
      validate: {
        isUUID: 4
      }
    },
    username: {
      type: DataTypes.STRING,
      allowNull: false,
      unique: true
    },
    email: {
      type: DataTypes.STRING,
      validate: {
        isEmail: true
      },
      allowNull: false,
      unique: true
    },
    password: {
      type: DataTypes.STRING,
      allowNull: false,
      set (val) {
        const salt = generateSalt()
        const hashed = sha512(val, salt)
        this.setDataValue('password', hashed)
        this.setDataValue('salt', salt)
      }
    },
    isActive: {
      type: DataTypes.BOOLEAN,
      allowNull: false,
      defaultValue: false
    },
    salt: {
      type: DataTypes.STRING,
      allowNull: false
    }
  })

  User.associate = function (models) {
    User.MultiStreams = User.hasMany(models.MultiStream)
  }

  User.prototype.checkPassword = function (input) {
    return this.getDataValue('password') === sha512(input, this.getDataValue('salt'))
  }

  User.prototype.toJSON = function () {
    var values = Object.assign({}, this.get({ plain: true }))

    delete values.password
    delete values.salt
    return values
  }

  return User
}
