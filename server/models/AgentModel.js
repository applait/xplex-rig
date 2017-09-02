/**
 * Model for storing agent information
 */

module.exports = function (sequelize, DataTypes) {
  var Agent = sequelize.define('Agent', {
    hostname: {
      type: DataTypes.STRING,
      primary: true
    },
    isActive: {
      type: DataTypes.BOOLEAN,
      allowNull: false,
      defaultValue: false
    }
  })

  Agent.associate = function (models) {
    Agent.belongsTo(models.EdgeCluster)
  }

  return Agent
}
