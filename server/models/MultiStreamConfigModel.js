/**
 * Model for storing multi-streaming user configurations
 */

module.exports = function (sequelize, DataTypes) {
  var MultiStreamConfig = sequelize.define('MultiStreamConfig', {
    service: {
      type: DataTypes.STRING,
      allowNull: false
    },
    server: {
      type: DataTypes.STRING,
      allowNull: true
    },
    key: {
      type: DataTypes.STRING,
      allowNull: false
    },
    isActive: {
      type: DataTypes.BOOLEAN,
      allowNull: false,
      defaultValue: true
    }
  })

  MultiStreamConfig.associate = function (models) {
    MultiStreamConfig.User = MultiStreamConfig.belongsTo(models.User)
  }

  return MultiStreamConfig
}
