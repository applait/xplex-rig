// Loads and parses JSON config for rig-server and merges them with environment
// variables

package config

import "io/ioutil"
import "encoding/json"

type serverConfig struct {
	Port        int    `json:"port"`
	JWTSecret   string `json:"jwtsecret"`
	PostgresURL string `json:"postgresUrl"`
}

// Config holds the structure for configuration
type Config struct {
	Server serverConfig `json:"server"`
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
		serverConfig{
			Port:        8081,
			JWTSecret:   "replacemewithanicelongstringthatyouwillnotsharewithothers",
			PostgresURL: "postgres://user:pass@host/db",
		},
	}
	j, err := json.Marshal(c)
	if err != nil {
		return c, err
	}
	err = ioutil.WriteFile(confPath, j, 0644)
	if err != nil {
		return c, err
	}
	return c, nil
}
