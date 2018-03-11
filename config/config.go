// Loads and parses JSON config for rig-server and merges them with environment
// variables

package config

import "io/ioutil"
import "encoding/json"

// JWTKeys defines the different keys used by rig to sign and verify JWTs for different use cases. These keys need to be
// shared across all rig instances.
type JWTKeys struct {
	Users  string `json:"users"`
	Agents string `json:"agents"`
	Admins string `json:"admins"`
}

// Config holds the structure for configuration
type Config struct {
	Port        int     `json:"port"`
	JWTKeys     JWTKeys `json:"jwtKeys"`
	PostgresURL string  `json:"postgresUrl"`
}

// ParseConfig parses a configuration file and loads into a config
func ParseConfig(confPath string) (Config, error) {
	var config Config
	file, err := ioutil.ReadFile(confPath)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

// CreateConfig generates a default JSON config file and writes to the given path
func CreateConfig(confPath string) (Config, error) {
	c := Config{
		Port: 8081,
		JWTKeys: JWTKeys{
			Users:  "keyhere",
			Agents: "keyhere",
			Admins: "keyhere",
		},
		PostgresURL: "postgres://user:pass@host/db",
	}
	j, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return c, err
	}
	err = ioutil.WriteFile(confPath, j, 0644)
	if err != nil {
		return c, err
	}
	return c, nil
}
