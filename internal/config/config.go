package config

import (
	"log"
	"os"

	yaml "github.com/goccy/go-yaml"
)

const (
	defaultPort       = 8080
	defaultLogDir     = "logs"
	defaultDBHost     = "localhost"
	defaultDBPort     = 3306
	defaultDBUsername = "root"
	defaultDBName     = "sealos-complik-admin"
	defaultDBPassword = "123456"
)

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Config struct {
	Port     int            `yaml:"port"`
	LogDir   string         `yaml:"log_dir"`
	Database DatabaseConfig `yaml:"database"`
}

// LoadConfig loads the configuration from the specified YAML file and environment variables.
func LoadConfig(configFile string) *Config {
	// Set default values
	cfg := &Config{
		Port:   defaultPort,
		LogDir: defaultLogDir,
		Database: DatabaseConfig{
			Host:     defaultDBHost,
			Port:     defaultDBPort,
			Username: defaultDBUsername,
			Password: defaultDBPassword, // Get DB password from environment variable
			Name:     defaultDBName,
		},
	}
	// Load base config from file
	if err := loadConfigInto(configFile, cfg, false); err != nil {
		log.Printf("read config file %q failed: %v, using default config", configFile, err)
	}

	return cfg
}

// loadConfigInto loads the YAML configuration from the specified file into the provided Config struct.
func loadConfigInto(configFile string, cfg *Config, optional bool) error {
	content, err := os.ReadFile(configFile)
	if err != nil {
		if optional && os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if err := yaml.Unmarshal(content, cfg); err != nil {
		return err
	}

	return nil
}
