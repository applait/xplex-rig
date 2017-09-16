// Loads and parses JSON config for rig-server and merges them with environment
// variables

package main

import "io/ioutil"
import "encoding/json"

// Config holds the structure for configuration
type Config struct {
	Server struct {
		Port        int    `json:"port"`
		JWTSecret   string `json:"jwtsecret"`
		PostgresURL string `json:"postgresUrl"`
	} `json:"server"`
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
