package main

import "github.com/BurntSushi/toml"

type Config struct {
	ListenAddr string
	Port       int
	Neighbours []string
	Fanout     int
}

func readConfig(cfgPath string) (*Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(cfgPath, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
