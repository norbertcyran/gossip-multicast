package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/norbertcyran/gossip-multicast/gossip"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	cfgPath := flag.String("c", "", "TOML configuration file")
	flag.Parse()

	if *cfgPath == "" {
		panic("Config file is required!")
	}

	cfg, err := readConfig(*cfgPath)
	if err != nil {
		panic(err)
	}

	udpCfg := &gossip.UDPTransportConfig{
		BindAddr: cfg.ListenAddr, BindPort: cfg.Port,
	}
	udp, err := gossip.NewUDPTransport(udpCfg)
	if err != nil {
		panic(err)
	}
	gossipCfg := &gossip.Config{
		Transport:            udp,
		Fanout:               cfg.Fanout,
		GossipInterval:       cfg.GossipInterval,
		RetransmitMultiplier: cfg.RetransmitMultiplier,
		Neighbours:           cfg.Neighbours,
	}
	_, err = gossip.StartService(gossipCfg)
	if err != nil {
		panic(err)
	}

	<-ctx.Done()
	fmt.Println("Shutting down...")
}
