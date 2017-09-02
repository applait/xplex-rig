/**
 * Model for storing information on edge clusters
 */

module.exports = function (sequelize, DataTypes) {
  var EdgeCluster = sequelize.define('EdgeCluster', {
    publicDNS: {
      type: DataTypes.STRING,
      allowNull: false,
      unique: true
    },
    region: {
      type: DataTypes.STRING,
      allowNull: false
    },
    provider: {
      type: DataTypes.STRING,
      allowNull: false
    },
    isActive: {
      type: DataTypes.BOOLEAN,
      allowNull: false,
      defaultValue: false
    }
  })

  EdgeCluster.associate = function (models) {
    EdgeCluster.hasMany(models.Agent)
  }

  return EdgeCluster
}
