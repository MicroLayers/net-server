package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"

	"net-server/configuration"
	"net-server/listener"
	"net-server/module"

	log "github.com/sirupsen/logrus"
)

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
	configuration configuration.Configuration,
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

	configuration, confMapSlice := configuration.ReadConf(configPath)

	log.WithFields(log.Fields{
		"Listen.Http.Enabled": configuration.Server.Listen.HTTP.Enabled,
		"Listen.Tcp.Enabled":  configuration.Server.Listen.TCP.Enabled,
		"Listen.Unix.Enabled": configuration.Server.Listen.Unix.Enabled,
		"Module":              configuration.Server.Module,
	}).Info("Loaded base configuration data")

	mod, err := module.LoadModule(configuration.Server.Module)

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	mod.Init(confMapSlice)

	startServers(ctx, configuration, mod)
}
