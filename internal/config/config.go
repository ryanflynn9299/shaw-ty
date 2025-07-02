package config

import (
	"github.com/spf13/viper"
	"log"
)

type AppConfig struct {
	MachineID  int    `mapstructure:"machine_id"`
	DBType     string `mapstructure:"db_type"`
	DBHost     string `mapstructure:"db_host"`
	DBPort     int    `mapstructure:"db_port"`
	DBUser     string `mapstructure:"db_user"`
	DBPassword string `mapstructure:"db_password"`
	DBName     string `mapstructure:"db_name"`
	DBPath     string `mapstructure:"db_path"` // enable SQLite support
	ServerHost string `mapstructure:"server_host"`
	ServerPort int    `mapstructure:"server_port"`
	IsDevMode  bool   `mapstructure:"is_dev_mode"`
}

// LoadConfig loads the configuration using viper
func LoadConfig() (*AppConfig, error) {
	viper.SetConfigName("config")          // Name of the config file (without extension)
	viper.SetConfigType("yaml")            // File format
	viper.AddConfigPath("internal/config") // Look for config in the current directory
	viper.AutomaticEnv()                   // Override with environment variables

	// Default values (if nothing is provided)
	viper.SetDefault("machine_id", 0)
	viper.SetDefault("db_type", "sqlite")
	viper.SetDefault("db_host", "localhost")
	viper.SetDefault("db_path", "localhost")
	viper.SetDefault("db_port", 3306)
	viper.SetDefault("db_user", "root")
	viper.SetDefault("db_password", "root")
	viper.SetDefault("db_name", "shortLink")
	viper.SetDefault("server_port", 8080)
	viper.SetDefault("server_host", "localhost")
	viper.SetDefault("is_dev_mode", true)
	viper.SetDefault("jwt_expiry_duration", 4) // hours

	// Attempt to read from the configuration file
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file not found, using environment variables and defaults: %v", err)
	}

	var config AppConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
