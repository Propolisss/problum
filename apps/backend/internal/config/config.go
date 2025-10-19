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

	// db
	defaultDBHost          = "postgres"
	defaultDBPort          = 5432
	defaultUser            = "problum"
	defaultDBName          = "problum"
	defaultSSLMode         = "disable"
	defaultMaxOpenConns    = 200
	defaultMaxIdleConns    = 20
	defaultConnMaxLifetime = time.Duration(300) * time.Second
)

type Config struct {
	Server *Server
	DB     *DB
}

type Server struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
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

func readServerConfig() *Server {
	return &Server{
		Host:            viper.GetString("server.host"),
		Port:            viper.GetInt("server.port"),
		ReadTimeout:     viper.GetDuration("server.read_timeout"),
		WriteTimeout:    viper.GetDuration("server.write_timeout"),
		ShutdownTimeout: viper.GetDuration("server.shutdown_timeout"),
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

func setDefault() {
	viper.SetDefault("server.host", defaultServerHost)
	viper.SetDefault("server.port", defaultServerPort)
	viper.SetDefault("server.read_timeout", defaultServerReadTimeout)
	viper.SetDefault("server.write_timeout", defaultServerWriteTimeout)
	viper.SetDefault("server.shutdown_timeout", defaultServerShutdownTimeout)
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

	return &Config{
		Server: serverConfig,
		DB:     dbConfig,
	}, nil
}
