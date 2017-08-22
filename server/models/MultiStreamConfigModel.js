/**
 * Model for storing user-defined output configurations for `MultiStream`
 */

module.exports = function (sequelize, DataTypes) {
  var MultiStreamConfig = sequelize.define('MultiStreamConfig', {
    service: {
      type: DataTypes.STRING,
      allowNull: false
    },
    key: {
      type: DataTypes.STRING,
      allowNull: false
    },
    server: {
      type: DataTypes.STRING,
      defaultValue: 'default',
      allowNull: true
    },
    isActive: {
      type: DataTypes.BOOLEAN,
      allowNull: false,
      defaultValue: true
    }
  })

  MultiStreamConfig.associate = function (models) {
    MultiStreamConfig.belongsTo(models.MultiStream)
  }

  return MultiStreamConfig
}
