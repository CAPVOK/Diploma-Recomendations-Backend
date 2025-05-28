package config

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	Environment   string        `validate:"required,oneof=development staging production"`
	InternalToken string        `mapstructure:"internal_token"`
	Server        ServerConfig  `validate:"required"`
	GRPC          GRPCConfig    `validate:"required"`
	DB            DBConfig      `mapstructure:"postgres"`
	Auth          AuthConfig    `mapstructure:"auth"`
	Logging       LoggingConfig `mapstructure:"logging"`
	KafkaProducer KafkaProducer `mapstructure:"kafka_producer"`
}

type GRPCConfig struct {
	DetectorAddr string `mapstructure:"detector_url"`
}

type ServerConfig struct {
	Host            string        `validate:"required"`
	Port            int           `validate:"required,min=1,max=65535"`
	ReadTimeout     time.Duration `validate:"required" mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `validate:"required" mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `validate:"required" mapstructure:"shutdown_timeout"`
	CORSAllowed     []string
}

type DBConfig struct {
	Host            string        `validate:"required" mapstructure:"host"`
	Port            int           `validate:"required,min=1,max=65535" mapstructure:"port"`
	User            string        `validate:"required" mapstructure:"user"`
	Password        string        `validate:"required" mapstructure:"password"`
	DBName          string        `validate:"required" mapstructure:"dbname"`
	SSLMode         string        `validate:"oneof=disable allow prefer require verify-full" mapstructure:"sslmode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type AuthConfig struct {
	JWTSecret          string        `validate:"required" mapstructure:"jwt_secret"`
	AccessTokenExpire  time.Duration `validate:"required" mapstructure:"access_token_expire"`
	RefreshTokenExpire time.Duration `validate:"required" mapstructure:"refresh_token_expire"`
	PasswordCost       int           `validate:"min=4,max=14" mapstructure:"password_cost"`
}

type LoggingConfig struct {
	Level          string `validate:"oneof=debug info warn error"`
	JSONFormat     bool
	LogFilePath    string
	RotationPolicy struct {
		MaxSize    int `mapstructure:"max_size"`
		MaxBackups int `mapstructure:"max_backups"`
		MaxAge     int `mapstructure:"max_age"`
	}
}

type KafkaProducer struct {
	Broker string `mapstructure:"broker"`
}

func MustLoad() *Config {
	var cfgPath string

	flag.StringVar(&cfgPath, "config", "config/dev.yaml", "path to config file")
	flag.Parse()

	if cfgPath == "" {
		if envPath := os.Getenv("CONFIG_PATH"); envPath != "" {
			cfgPath = envPath
		} else {
			log.Fatal("config path must be set via --config flag or CONFIG_PATH environment variable")
		}
	}

	v := viper.New()
	v.SetConfigFile(cfgPath)
	v.SetConfigType("yaml")

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %s", err.Error())
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("failed to unmarshal config: %s", err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		log.Fatalf("failed to validate config: %s", err.Error())
	}

	return &cfg
}

func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", 15*time.Second)
	v.SetDefault("server.write_timeout", 15*time.Second)
	v.SetDefault("server.shutdown_timeout", 5*time.Second)
	v.SetDefault("server.cors_allowed", []string{"*"})

	// Database defaults
	v.SetDefault("postgres.host", "localhost")
	v.SetDefault("postgres.port", 5432)
	v.SetDefault("postgres.sslmode", "disable")
	v.SetDefault("postgres.max_open_conns", 10)
	v.SetDefault("postgres.max_idle_conns", 5)
	v.SetDefault("postgres.conn_max_lifetime", 30*time.Minute)

	// Auth defaults
	v.SetDefault("auth.access_token_expire", 15*time.Minute)
	v.SetDefault("auth.refresh_token_expire", 24*time.Hour*7) // 1 week
	v.SetDefault("auth.password_cost", 10)

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.json_format", false)
	v.SetDefault("logging.rotation_policy.max_size", 100) // MB
	v.SetDefault("logging.rotation_policy.max_backups", 3)
	v.SetDefault("logging.rotation_policy.max_age", 30) // days
}
