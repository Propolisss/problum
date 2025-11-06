package config

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	// server
	defaultServerHost            = "0.0.0.0"
	defaultServerPort            = 8080
	defaultServerReadTimeout     = time.Duration(30) * time.Second
	defaultServerWriteTimeout    = time.Duration(30) * time.Second
	defaultServerShutdownTimeout = time.Duration(10) * time.Second
	defaultServerEnvironment     = "dev"

	// db
	defaultDBHost            = "postgres"
	defaultDBPort            = 5432
	defaultDBUser            = "problum"
	defaultDBName            = "problum"
	defaultDBSSLMode         = "disable"
	defaultDBMaxOpenConns    = 200
	defaultDBMaxIdleConns    = 20
	defaultDBConnMaxLifetime = time.Duration(300) * time.Second

	// redis
	defaultRedisHost     = "0.0.0.0"
	defaultRedisPort     = 6379
	defaultRedisPassword = ""
	defaultRedisDB       = 0

	// nats
	defaultNatsHost = "0.0.0.0"
	defaultNatsPort = 4222
)

type Config struct {
	Server *Server
	DB     *DB
	Redis  *Redis
	Nats   *Nats
}

type Server struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
	Environment     string        `mapstructure:"environment"`
}

type DB struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"db_name"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type Nats struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func readServerConfig() *Server {
	return &Server{
		Host:            viper.GetString("server.host"),
		Port:            viper.GetInt("server.port"),
		ReadTimeout:     viper.GetDuration("server.read_timeout"),
		WriteTimeout:    viper.GetDuration("server.write_timeout"),
		ShutdownTimeout: viper.GetDuration("server.shutdown_timeout"),
		Environment:     viper.GetString("server.environment"),
	}
}

func readDBConfig() *DB {
	return &DB{
		Host:            viper.GetString("db.host"),
		Port:            viper.GetInt("db.port"),
		User:            viper.GetString("db.user"),
		Password:        viper.GetString("db.password"),
		DBName:          viper.GetString("db.db_name"),
		SSLMode:         viper.GetString("db.ssl_mode"),
		MaxOpenConns:    viper.GetInt("db.max_open_conns"),
		MaxIdleConns:    viper.GetInt("db.max_idle_conns"),
		ConnMaxLifetime: viper.GetDuration("db.conn_max_lifetime"),
	}
}

func readRedisConfig() *Redis {
	return &Redis{
		Host:     viper.GetString("redis.host"),
		Port:     viper.GetInt("redis.port"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	}
}

func readNatsConfig() *Nats {
	return &Nats{
		Host: viper.GetString("nats.host"),
		Port: viper.GetInt("nats.port"),
	}
}

func setDefault() {
	// server
	viper.SetDefault("server.host", defaultServerHost)
	viper.SetDefault("server.port", defaultServerPort)
	viper.SetDefault("server.read_timeout", defaultServerReadTimeout)
	viper.SetDefault("server.write_timeout", defaultServerWriteTimeout)
	viper.SetDefault("server.shutdown_timeout", defaultServerShutdownTimeout)
	viper.SetDefault("server.environment", defaultServerEnvironment)

	// db
	viper.SetDefault("db.host", defaultDBHost)
	viper.SetDefault("db.port", defaultDBPort)
	viper.SetDefault("db.user", defaultDBUser)
	viper.SetDefault("db.db_name", defaultDBName)
	viper.SetDefault("db.ssl_mode", defaultDBSSLMode)
	viper.SetDefault("db.max_open_conns", defaultDBMaxOpenConns)
	viper.SetDefault("db.max_idle_conns", defaultDBMaxIdleConns)
	viper.SetDefault("db.max_lifetime", defaultDBConnMaxLifetime)

	// redis
	viper.SetDefault("redis.host", defaultRedisHost)
	viper.SetDefault("redis.port", defaultRedisPort)
	viper.SetDefault("redis.password", defaultRedisPassword)
	viper.SetDefault("redis.db", defaultRedisDB)

	// nats
	viper.SetDefault("nats.host", defaultNatsHost)
	viper.SetDefault("nats.port", defaultNatsPort)
}

func (c *DB) GetDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DBName,
		c.SSLMode,
	)
}

func (s *Server) IsProduction() bool {
	return s.Environment == "prod"
}

func (cfg *Config) IsProduction() bool {
	return cfg.Server.IsProduction()
}

func New() (*Config, error) {
	configPath := os.Getenv("PROBLUM_CONFIG_FILE")
	log.Info().Msgf("config path: %s", configPath)
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Error().Err(err).Msg("failed to read config file")
		return nil, err
	}

	setDefault()
	serverConfig := readServerConfig()
	dbConfig := readDBConfig()
	redisConfig := readRedisConfig()
	natsConifg := readNatsConfig()

	return &Config{
		Server: serverConfig,
		DB:     dbConfig,
		Redis:  redisConfig,
		Nats:   natsConifg,
	}, nil
}
