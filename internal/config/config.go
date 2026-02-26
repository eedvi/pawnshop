package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Storage  StorageConfig
	Logging  LoggingConfig
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

type LoggingConfig struct {
	Level              string        // debug, info, warn, error
	Format             string        // json, console
	SlowQueryThreshold time.Duration // threshold for slow query logging
	LogAllQueries      bool          // log all database queries (debug mode)
}

func Load() (*Config, error) {
	// Load .env file if exists (development)
	// In production, use real environment variables
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	if err := godotenv.Load(envFile); err != nil {
		// .env is optional, only warn if explicitly set
		if envFile != ".env" {
			return nil, fmt.Errorf("error loading env file %s: %w", envFile, err)
		}
		// .env not found is OK - will use config.yaml and env vars
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/pawnshop")

	// Environment variables override config file
	viper.AutomaticEnv()
	viper.SetEnvPrefix("PAWN")

	// Map environment variables to config keys
	// This allows both PAWN_DATABASE_HOST and DB_HOST to work
	bindEnvVariables()

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

	// Logging
	config.Logging = LoggingConfig{
		Level:              viper.GetString("logging.level"),
		Format:             viper.GetString("logging.format"),
		SlowQueryThreshold: viper.GetDuration("logging.slow_query_threshold"),
		LogAllQueries:      viper.GetBool("logging.log_all_queries"),
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

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "console")
	viper.SetDefault("logging.slow_query_threshold", "1s")
	viper.SetDefault("logging.log_all_queries", false)
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

// bindEnvVariables maps common environment variable names to viper keys
// This allows using standard names like DB_HOST instead of PAWN_DATABASE_HOST
func bindEnvVariables() {
	// Database
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.dbname", "DB_NAME")
	viper.BindEnv("database.sslmode", "DB_SSL_MODE")

	// Redis
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("redis.db", "REDIS_DB")

	// JWT
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("jwt.issuer", "JWT_ISSUER")

	// App
	viper.BindEnv("app.environment", "APP_ENV")
	viper.BindEnv("app.debug", "APP_DEBUG")

	// Storage
	viper.BindEnv("storage.endpoint", "S3_ENDPOINT")
	viper.BindEnv("storage.access_key", "S3_ACCESS_KEY")
	viper.BindEnv("storage.secret_key", "S3_SECRET_KEY")
	viper.BindEnv("storage.bucket", "S3_BUCKET")
	viper.BindEnv("storage.region", "S3_REGION")
}
