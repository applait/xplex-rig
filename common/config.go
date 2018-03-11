// Loads and parses JSON config for rig-server and merges them with environment
// variables

package common

import "io/ioutil"
import "encoding/json"

// ParseConfig parses a configuration file and loads into a config
func ParseConfig(confPath string) (JSONConfig, error) {
	file, err := ioutil.ReadFile(confPath)
	if err != nil {
		return Config, err
	}
	err = json.Unmarshal(file, &Config)
	if err != nil {
		return Config, err
	}
	return Config, nil
}

// CreateConfig generates a default JSON config file and writes to the given path
func CreateConfig(confPath string) (JSONConfig, error) {
	c := JSONConfig{
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
