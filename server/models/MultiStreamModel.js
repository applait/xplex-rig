/**
 * Model for storing multi-stream information
 */

module.exports = function (sequelize, DataTypes) {
  var MultiStream = sequelize.define('MultiStream', {
    key: {
      type: DataTypes.STRING,
      allowNull: false,
      unique: true
    },
    isActive: {
      type: DataTypes.BOOLEAN,
      allowNull: false,
      defaultValue: true
    },
    isStreaming: {
      type: DataTypes.BOOLEAN,
      allowNull: false,
      defaultValue: false
    }
  })

  MultiStream.associate = function (models) {
    MultiStream.hasMany(models.MultiStreamConfig)
    MultiStream.belongsTo(models.User)
  }

  return MultiStream
}
