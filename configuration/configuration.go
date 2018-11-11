package configuration

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type Configuration struct {
	Server struct {
		Module string `yaml:"module"`
		Listen struct {
			Unix struct {
				Enabled  bool   `yaml:"enabled"`
				Protocol string `yaml:"protocol"`
				Socket   string `yaml:"socket"`
			} `yaml:"unix"`
			Tcp struct {
				Enabled  bool   `yaml:"enabled"`
				Protocol string `yaml:"protocol"`
				Port     uint16 `yaml:"port"`
			} `yaml:"tcp"`
			Http struct {
				Enabled  bool   `yaml:"enabled"`
				Protocol string `yaml:"protocol"`
				Port     uint16 `yaml:"port"`
			} `yaml:"http"`
		} `yaml:"listen"`
	} `yaml:"server"`
}

func DefaultConfiguration() (configuration Configuration) {
	configuration.Server.Listen.Http.Enabled = true
	configuration.Server.Listen.Http.Port = 8080
	configuration.Server.Listen.Http.Protocol = "json"

	configuration.Server.Listen.Tcp.Enabled = true
	configuration.Server.Listen.Tcp.Port = 15252
	configuration.Server.Listen.Tcp.Protocol = "proto"

	configuration.Server.Listen.Unix.Enabled = true
	configuration.Server.Listen.Unix.Protocol = "proto"
	configuration.Server.Listen.Unix.Socket = "/var/run/net-server.socket"

	return configuration
}

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
		yaml.Unmarshal(bytes, &mapSlice)
		if err != nil {
			log.Warn("Failed to read configuration file, reading default")
			configuration = DefaultConfiguration()
		}
	}

	return configuration, mapSlice
}
