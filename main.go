package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"plugin"
	"sync"

	"net-server/listener"
	"net-server/module"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type config struct {
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

func loadDefaultConfigValues(configuration *config) {
	configuration.Server.Listen.Http.Enabled = true
	configuration.Server.Listen.Http.Port = 8080
	configuration.Server.Listen.Http.Protocol = "json"

	configuration.Server.Listen.Tcp.Enabled = true
	configuration.Server.Listen.Tcp.Port = 15252
	configuration.Server.Listen.Tcp.Protocol = "proto"

	configuration.Server.Listen.Unix.Enabled = true
	configuration.Server.Listen.Unix.Protocol = "proto"
	configuration.Server.Listen.Unix.Socket = "/var/run/net-server.socket"
}

func readConf(path string) (configuration config) {
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
		loadDefaultConfigValues(&configuration)
	} else {
		log.WithField("Configuration file", path).Info("Reading configuration file")
		err = yaml.Unmarshal(bytes, &configuration)

		if err != nil {
			log.Warn("Failed to read configuration file, reading default")
			loadDefaultConfigValues(&configuration)
		}
	}

	return configuration
}

func sigintTrap(ctx context.Context) context.Context {
	appContext, cancel := context.WithCancel(ctx)

	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt)

	go func() {
		<-interruptChannel
		log.Info("Received SIGINT signal, stopping...")
		cancel()
	}()

	return appContext
}

func startServers(
	ctx context.Context,
	configuration config,
	configPath string,
	mod module.NetServerModule,
) {
	serversWaitGroup := sync.WaitGroup{}
	if configuration.Server.Listen.Unix.Enabled {
		log.WithFields(log.Fields{
			"server":   "unix",
			"protocol": configuration.Server.Listen.Unix.Protocol,
			"socket":   configuration.Server.Listen.Unix.Socket,
		}).Info("Starting listening service")
		serversWaitGroup.Add(1)
		go func() {
			err := listener.ListenUnix(
				ctx,
				mod,
				configuration.Server.Listen.Unix.Socket,
				configuration.Server.Listen.Unix.Protocol,
				configPath,
			)
			if err != nil {
				log.WithFields(log.Fields{
					"server": "unix",
					"error":  err,
				}).Error("Unexpected error")
			}
			serversWaitGroup.Done()
		}()
	}
	serversWaitGroup.Wait()
}

func main() {
	ctx := sigintTrap(context.Background())
	configPtr := flag.String("config", "", "Configuration file")
	flag.Parse()

	configPath := *configPtr

	configuration := readConf(configPath)

	log.WithFields(log.Fields{
		"Listen.Http.Enabled": configuration.Server.Listen.Http.Enabled,
		"Listen.Tcp.Enabled":  configuration.Server.Listen.Tcp.Enabled,
		"Listen.Unix.Enabled": configuration.Server.Listen.Unix.Enabled,
		"Module":              configuration.Server.Module,
	}).Info("Loaded base configuration data")

	plug, err := plugin.Open(configuration.Server.Module)

	if err != nil {
		log.WithError(err).Error("Impossible to open the plugin file")
		os.Exit(1)
	}

	symPlugin, err := plug.Lookup("NetServerModule")

	if err != nil {
		log.WithError(err).Error("Impossible to find the symbol NetServerModule")
		os.Exit(2)
	}

	mod, ok := symPlugin.(module.NetServerModule)

	if !ok {
		log.Error(fmt.Sprintf(
			"Module %s does not implement the required interface",
			configuration.Server.Module,
		))
		os.Exit(3)
	}

	if ok {
		startServers(ctx, configuration, configPath, mod)
	}
}
