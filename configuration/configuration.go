package configuration

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// Configuration the net-server configuration struct
type Configuration struct {
	Server struct {
		Module string `yaml:"module"`
		Listen struct {
			Unix struct {
				Enabled  bool   `yaml:"enabled"`
				Protocol string `yaml:"protocol"`
				Socket   string `yaml:"socket"`
			} `yaml:"unix"`

			TCP struct {
				Enabled  bool   `yaml:"enabled"`
				Port     uint16 `yaml:"port"`
				Protocol string `yaml:"protocol"`
			} `yaml:"tcp"`

			HTTP struct {
				Enabled  bool   `yaml:"enabled"`
				Port     uint16 `yaml:"port"`
				Protocol string `yaml:"protocol"`
			} `yaml:"http"`
		} `yaml:"listen"`
	} `yaml:"server"`
}

// DefaultConfiguration get the default net server's configuration
func DefaultConfiguration() (configuration Configuration) {
	configuration.Server.Listen.HTTP.Enabled = true
	configuration.Server.Listen.HTTP.Port = 8080
	configuration.Server.Listen.HTTP.Protocol = "json"

	configuration.Server.Listen.TCP.Enabled = true
	configuration.Server.Listen.TCP.Port = 15252
	configuration.Server.Listen.TCP.Protocol = "proto"

	configuration.Server.Listen.Unix.Enabled = true
	configuration.Server.Listen.Unix.Protocol = "proto"
	configuration.Server.Listen.Unix.Socket = "/var/run/net-server.socket"

	return configuration
}

// ReadConf read the net server's configuration from the provided path
func ReadConf(path string) (configuration Configuration, mapSlice yaml.MapSlice) {
	bytes, err := ioutil.ReadFile(path)
	if path == "" {
		// try with default locations
		paths := []string{
			"configuration.yml",                       // local
			"dist/configuration.yml",                  // local
			"/etc/net-server/configuration.yml",       // etc
			"/etc/local/net-server/configuration.yml", // local etc
		}

		for _, path = range paths {
			bytes, err = ioutil.ReadFile(path)

			if err == nil {
				break
			}
		}
	}

	if err != nil {
		log.WithField("Configuration file", path).Warn("Failed to read configuration file")
		configuration = DefaultConfiguration()
	} else {
		log.WithField("Configuration file", path).Info("Reading configuration file")

		err = yaml.Unmarshal(bytes, &configuration)
		if err == nil {
			err = yaml.Unmarshal(bytes, &mapSlice)
		}
		if err != nil {
			log.Warn("Failed to read configuration file, reading default")
			configuration = DefaultConfiguration()
		}
	}

	return configuration, mapSlice
}
