package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	pb "github.com/averak/gamebox/protobuf/config"
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	config *Config
	once   sync.Once
)

const (
	envKey = "GAMEBOX_CONFIG_FILEPATH"
)

type Config = pb.Config

func Get() *Config {
	once.Do(func() {
		c, err := load()
		if err != nil {
			log.Panicf("Failed to load config: %v", err)
		}
		config = c
	})
	return config
}

func load() (*Config, error) {
	filepath, ok := os.LookupEnv(envKey)
	if !ok {
		return nil, fmt.Errorf("ENV[%s] is not found", envKey)
	}

	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var conf Config
	err = protojson.Unmarshal(bytes, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
