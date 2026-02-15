package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Storage  StorageConfig
}

type AppConfig struct {
	Name        string
	Environment string
	Debug       bool
	Version     string
}

type ServerConfig struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type JWTConfig struct {
	Secret           string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
	Issuer           string
}

type StorageConfig struct {
	Type      string // local, s3, minio
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/pawnshop")

	// Environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("PAWN")

	// Defaults
	setDefaults()

	// Read config file if exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config

	// App
	config.App = AppConfig{
		Name:        viper.GetString("app.name"),
		Environment: viper.GetString("app.environment"),
		Debug:       viper.GetBool("app.debug"),
		Version:     viper.GetString("app.version"),
	}

	// Server
	config.Server = ServerConfig{
		Host:         viper.GetString("server.host"),
		Port:         viper.GetInt("server.port"),
		ReadTimeout:  viper.GetDuration("server.read_timeout"),
		WriteTimeout: viper.GetDuration("server.write_timeout"),
		IdleTimeout:  viper.GetDuration("server.idle_timeout"),
	}

	// Database
	config.Database = DatabaseConfig{
		Host:            viper.GetString("database.host"),
		Port:            viper.GetInt("database.port"),
		User:            viper.GetString("database.user"),
		Password:        viper.GetString("database.password"),
		DBName:          viper.GetString("database.dbname"),
		SSLMode:         viper.GetString("database.sslmode"),
		MaxOpenConns:    viper.GetInt("database.max_open_conns"),
		MaxIdleConns:    viper.GetInt("database.max_idle_conns"),
		ConnMaxLifetime: viper.GetDuration("database.conn_max_lifetime"),
	}

	// Redis
	config.Redis = RedisConfig{
		Host:     viper.GetString("redis.host"),
		Port:     viper.GetInt("redis.port"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	}

	// JWT
	config.JWT = JWTConfig{
		Secret:          viper.GetString("jwt.secret"),
		AccessTokenTTL:  viper.GetDuration("jwt.access_token_ttl"),
		RefreshTokenTTL: viper.GetDuration("jwt.refresh_token_ttl"),
		Issuer:          viper.GetString("jwt.issuer"),
	}

	// Storage
	config.Storage = StorageConfig{
		Type:      viper.GetString("storage.type"),
		Endpoint:  viper.GetString("storage.endpoint"),
		AccessKey: viper.GetString("storage.access_key"),
		SecretKey: viper.GetString("storage.secret_key"),
		Bucket:    viper.GetString("storage.bucket"),
		Region:    viper.GetString("storage.region"),
	}

	return &config, nil
}

func setDefaults() {
	// App defaults
	viper.SetDefault("app.name", "Pawnshop")
	viper.SetDefault("app.environment", "development")
	viper.SetDefault("app.debug", true)
	viper.SetDefault("app.version", "1.0.0")

	// Server defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "15s")
	viper.SetDefault("server.write_timeout", "15s")
	viper.SetDefault("server.idle_timeout", "60s")

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.dbname", "pawnshop")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", "5m")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// JWT defaults
	viper.SetDefault("jwt.secret", "your-super-secret-key-change-in-production")
	viper.SetDefault("jwt.access_token_ttl", "15m")
	viper.SetDefault("jwt.refresh_token_ttl", "168h") // 7 days
	viper.SetDefault("jwt.issuer", "pawnshop")

	// Storage defaults
	viper.SetDefault("storage.type", "local")
	viper.SetDefault("storage.bucket", "pawnshop")
}

// DSN returns the PostgreSQL connection string
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// URL returns the PostgreSQL connection URL
func (c *DatabaseConfig) URL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode,
	)
}

// RedisAddr returns the Redis address
func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
