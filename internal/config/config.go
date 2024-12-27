package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

// we kept the camel case so that it can be exported with the help of this package
//run go get -u github.com/ilyakaznacheev/cleanenv  , we can do annotations with the help of this package
//run go get github.com/go-playground/validator/v10, req validation

type HTTPServer struct {
	Addr string `yaml:"address" env-required:"true"`
}
type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string `yaml:"storage-path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

// config is parsed and stored in the struct here
func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()
		configPath = *flags //dereferencing the flag

		if configPath == "" {
			log.Fatal("Config path is not set ")
		}
	}

	//stat returns the file information
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config path does not exist: %s", configPath)
	}

	var cfg Config

	//first paraemter is configpath and second parameter is the structure in which it will be read
	err := cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		//printing string message of error
		log.Fatalf("Config path does not exist: %s", err.Error())
	}

	//return address
	return &cfg
}
