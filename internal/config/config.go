package config

import (
	"fmt"
	"os"
	"path"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config - structure for storing application configuration
type Config struct {
	App           App
	Logger        Logger
	HTTPServer    HTTPServer
	Postgres      Postgres
	Authorization Authorization
	Hasher        Hasher
	Documents     Documents
}

// Init - a function for initializing the application configuration
func Init(configPath string, envFile string) (*Config, error) {
	config := Config{}

	// Define a list of configuration files
	configFiles := []struct {
		fileName string
		key      string
		rawVal   any
	}{
		{fileName: "app.yaml", key: "app", rawVal: &config.App},
		{fileName: "logger.yaml", key: "logger", rawVal: &config.Logger},
		{fileName: "http_server.yaml", key: "http_server", rawVal: &config.HTTPServer},
		{fileName: "postgres.yaml", key: "postgres", rawVal: &config.Postgres},
		{fileName: "auth.yaml", key: "auth", rawVal: &config.Authorization},
		{fileName: "documents.yaml", key: "documents", rawVal: &config.Documents},
	}

	// Reading configuration from YAML files
	for _, configFile := range configFiles {
		if err := readConfigFile(configPath, configFile.fileName, configFile.key, configFile.rawVal); err != nil {
			return nil, err
		}
	}

	// Reading configuration from environment variables
	if err := configFromEnv(&config, envFile); err != nil {
		return nil, err
	}

	return &config, nil
}

// readConfigFile - function for reading configuration from YAML file
func readConfigFile(baseDir string, fileName string, key string, rawVal any) (err error) {
	// Set the path to the configuration file
	viper.SetConfigFile(path.Join(baseDir, fileName))

	// Read the configuration file
	if err = viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read in config file: file_name=%v: %v", fileName, err)
	}

	// Parse the configuration from the file
	if err := viper.UnmarshalKey(key, rawVal); err != nil {
		return fmt.Errorf("failed to parse config file: file_name=%v: key=%v: %v", fileName, key, err)
	}

	return nil
}

// configFromEnv - function for reading configuration from environment variables
func configFromEnv(cfg *Config, envFile string) error {
	// Read environment variables from a file
	if err := godotenv.Load(envFile); err != nil {
		return fmt.Errorf("failed to read env file: %s: %s", envFile, err)
	}

	// Fill the configuration structure with values ​​from environment variables
	cfg.Postgres.Username = os.Getenv("POSTGRES_USER")
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Postgres.DBName = os.Getenv("POSTGRES_DB")

	cfg.Authorization.AdminToken = os.Getenv("AUTH_ADMIN_TOKEN")
	cfg.Authorization.JWT.SigningKey = os.Getenv("AUTH_SIGNING_KEY")

	cfg.Hasher.Salt = os.Getenv("HASHER_SALT")

	return nil
}
