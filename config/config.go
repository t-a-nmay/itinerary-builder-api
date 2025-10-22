package config

import "os"

//holds the config values for the application
type Config struct {
	ServerAddress string
	OutputDir     string
}

//func NewConfig() initializes a new Config instance 
func NewConfig() *Config {
	//gets port from .env file if not found default port is 8080
	serverAddr := os.Getenv("SERVER_ADDRESS")
	if serverAddr == "" {
		serverAddr = ":8080"
	}

	//gets output directory from .env default is ./output
	outputDir := os.Getenv("OUTPUT_DIR")
	if outputDir == "" {
		outputDir = "./output"
	}
	
	//returns pointer to new Config instance
	return &Config{
		ServerAddress: serverAddr,
		OutputDir:     outputDir,
	}
}