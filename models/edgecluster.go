package models

// EdgeCluster stores DNS and location information for edge clusters
type EdgeCluster struct {
	ID        int
	PublicDNS string `sql:",notnull,unique"`
	Region    string `sql:",notnull"`
	Provider  string
	IsActive  bool `sql:",notnull,default:false"`

	AgentNodes []*AgentNode // EdgeCluster hasMany AgentNode
}

// AgentNode stores hostname and information
type AgentNode struct {
	ID       int
	Hostname string `sql:",notnull,unique"`
	IsActive bool   `sql:",notnull,default:false"`

	EdgeClusterID int // AgentNode belongsTo EdgeCluster
}
