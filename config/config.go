package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server       ServerConfig       `mapstructure:"server"`
	Database     DatabaseConfig     `mapstructure:"database"`
	ERPDatabase  ERPDatabaseConfig  `mapstructure:"erp_database"`
	ERPDBMapping map[string]string  `mapstructure:"erp_db_mapping"`
	SignatureKey SignatureKeyConfig `mapstructure:"signature"`
	JWT          JWTConfig          `mapstructure:"jwt"`
	Logger       LoggerConfig       `mapstructure:"logger"`
}

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
	ENV  string `mapstructure:"env"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"name"`
	Timeout  int    `mapstructure:"timeout"`
}

type ERPDatabaseConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DefaultDB string `mapstructure:"default_db"`
	Timeout   int    `mapstructure:"timeout"`
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	ExpiryHour int    `mapstructure:"expiry_hour"`
}

type LoggerConfig struct {
	Level string `mapstructure:"level"`
	Path  string `mapstructure:"path"`
}

type SignatureKeyConfig struct {
	Secret string `mapstructure:"signature_key"`
}

func LoadConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "."
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshal config: %w", err)
	}
	return config, nil
}

func MustConfig() *Config {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Fatal error loading config: %s", err)
	}
	return cfg
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
	)
}

func (c *Config) GetERPDatabaseDSN() string {
	return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&encrypt=disable&trustServerCertificate=true&connection+timeout=%d",
		c.ERPDatabase.User,
		c.ERPDatabase.Password,
		c.ERPDatabase.Host,
		c.ERPDatabase.Port,
		c.ERPDatabase.DefaultDB,
		c.ERPDatabase.Timeout,
	)
}

func (c *Config) GetJWTExpiry() time.Duration {
	return time.Duration(c.JWT.ExpiryHour) * time.Hour
}
