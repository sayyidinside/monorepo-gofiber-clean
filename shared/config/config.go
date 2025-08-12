package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	// General Config
	AppName   string `mapstructure:"APP_NAME"`
	Env       string `mapstructure:"ENV"`
	Port      string `mapstructure:"PORT"`
	AdminPass string `mapstructure:"ADMIN_PASS"`
	Debug     bool   `mapstructure:"DEBUG"`

	// JWT
	JwtRefreshPrivateSecret string `mapstructure:"JWT_REFRESH_PRIVATE_SECRET"`
	JwtRefreshPublicSecret  string `mapstructure:"JWT_REFRESH_PUBLIC_SECRET"`
	JwtAccessPrivateSecret  string `mapstructure:"JWT_ACCESS_PRIVATE_SECRET"`
	JwtAccessPublicSecret   string `mapstructure:"JWT_ACCESS_PUBLIC_SECRET"`
	JwtAccessTime           int    `mapstructure:"JWT_ACCESS_TIME"`
	JwtRefreshTime          int    `mapstructure:"JWT_REFRESH_TIME"`

	// Cors
	CorsMaxAge       int    `mapstructure:"CORS_MAX_AGE"`
	CorsAllowOrigins string `mapstructure:"CORS_ALLOW_ORIGINS"`
	CorsAllowMethods string `mapstructure:"CORS_ALLOW_METHODS"`

	// Rate Limit
	RateLimitMax int `mapstructure:"RATE_LIMIT_MAX"`
	RateLimitExp int `mapstructure:"RATE_LIMIT_EXPIRATION"`

	// Caching
	CacheExp int `mapstructure:"CACHE_EXPIRATION"`

	// Whitelist IP
	AllowedIPs string `mapstructure:"ALLOWED_IP"`

	// Email
	SmtpEmail    string `mapstructure:"SMTP_EMAIL"`
	SmtpPassword string `mapstructure:"SMTP_PASSWORD"`
	SmtpHost     string `mapstructure:"SMTP_HOST"`
	SmtPort      string `mapstructure:"SMTP_PORT"`

	// Redis
	RedisAddress  string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`

	// Database
	ProdDbUsername  string `mapstructure:"PROD_DB_USERNAME"`
	ProdDbPassword  string `mapstructure:"PROD_DB_PASSWORD"`
	ProdDbName      string `mapstructure:"PROD_DB_NAME"`
	ProdDbHost      string `mapstructure:"PROD_DB_HOST"`
	ProdDbPort      string `mapstructure:"PROD_DB_PORT"`
	DevDbUsername   string `mapstructure:"DEV_DB_USERNAME"`
	DevDbPassword   string `mapstructure:"DEV_DB_PASSWORD"`
	DevDbName       string `mapstructure:"DEV_DB_NAME"`
	DevDbHost       string `mapstructure:"DEV_DB_HOST"`
	DevDbPort       string `mapstructure:"DEV_DB_PORT"`
	LocalDbUsername string `mapstructure:"LOCAL_DB_USERNAME"`
	LocalDbPassword string `mapstructure:"LOCAL_DB_PASSWORD"`
	LocalDbName     string `mapstructure:"LOCAL_DB_NAME"`
	LocalDbHost     string `mapstructure:"LOCAL_DB_HOST"`
	LocalDbPort     string `mapstructure:"LOCAL_DB_PORT"`
}

var AppConfig *Config

func LoadConfig() error {
	viper.SetConfigFile("./.env")

	// Enable Viper to read environment variables
	viper.AutomaticEnv()

	// Set default for env
	viper.SetDefault("PORT", "4000")
	viper.SetDefault("JWT_ACCESS_TIME", 30)
	viper.SetDefault("JWT_REFRESH_TIME", 168)

	// Try to read the configuration file (optional)
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
	}

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		return err
	}

	// Watch for changes in the config file and reload AppConfig when changes occur
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// Reload the config file upon changes
		if err := viper.Unmarshal(AppConfig); err != nil {
			log.Printf("Error unmarshaling updated config: %s", err)
		}
	})

	return nil
}
