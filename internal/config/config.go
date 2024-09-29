package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	//here slash quotes are for encoding field  names in another format
	//for example: yaml, json, xml...
	//it is called struct tags

	//struct tags are like annotations, which can say to libs
	//how to work with these fields
	//yaml determines field name in yaml
	//env is field name if we will read it from environment vars of os
	//env-default - default value if there is no value in config.yaml
	//env-required - this field must have any value. Otherwise error
	Env         string `yaml:"env" env:"ENV" env-default:"local" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	//must prefix is used in case of panic inside function without returning error
	//it is ok not to return error in config-loading state, where the app has not loaded yet

	//get value from environment var of os
	configPath := os.Getenv("CONFIG_PRETTYLINK_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PRETTYLINK_PATH has not set")
	}

	//get file stat and if it is not exist, log fatal
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var cfg Config
	//read config file using cleanenv and handle possible error
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}

	return &cfg
}
