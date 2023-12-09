package configs

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/pelletier/go-toml"
)

var (
	ErrEnvVarNotSet = errors.New("environment variable is not set")
)

const (
	envConfig  = "deploy/docker-compose/.env"
	tomlConfig = "configs/default.toml"

	serverPortKey = "SERVICE_PORT"
	dbUserKey     = "POSTGRES_USER"
	dbPassKey     = "POSTGRES_PASSWORD"
	dbDatabaseKey = "POSTGRES_DB"
)

type AppConfig struct {
	Server struct {
		ListenAddr          string `toml:"listen_addr"`
		Port                int
		ReadTimeoutSeconds  int `toml:"read_timeout_seconds"`
		WriteTimeoutSeconds int `toml:"write_timeout_seconds"`
		IdleTimeoutSeconds  int `toml:"idle_timeout_seconds"`
	} `toml:"server"`

	Database struct {
		Host     string `toml:"host"`
		Port     int    `toml:"port"`
		User     string
		Password string
		DB       string
	} `toml:"database"`

	Logger struct {
		Level             string `toml:"level"`
		StackTraceEnabled bool   `toml:"stack_trace_enabled"`
	} `toml:"logger"`
}

func LoadDefault() (AppConfig, error) {
	configFile, err := os.OpenFile(tomlConfig, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return AppConfig{}, fmt.Errorf("opening toml config file: %w", err)
	}

	defer func() {
		err = errors.Join(err, configFile.Close())
		if err != nil {
			err = fmt.Errorf("deferred closing of toml config file: %w", err)
			log.Println(err)
		}
	}()

	cfg := AppConfig{}
	err = toml.NewDecoder(configFile).Decode(&cfg)
	if err != nil {
		return AppConfig{}, fmt.Errorf("decoding toml config: %w", err)
	}

	err = godotenv.Load(envConfig)
	if err != nil {
		return AppConfig{}, fmt.Errorf("loading env config: %w", err)
	}

	portStr, serverPortFound := os.LookupEnv(serverPortKey)
	if !serverPortFound {
		return AppConfig{},
			fmt.Errorf("looking for server port value in .env file: %w", ErrEnvVarNotSet)
	}

	cfg.Server.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return AppConfig{}, fmt.Errorf("parsing server port value: %w", err)
	}

	var foundUser, foundPass, foundDb bool
	cfg.Database.User, foundUser = os.LookupEnv(dbUserKey)
	cfg.Database.Password, foundPass = os.LookupEnv(dbPassKey)
	cfg.Database.DB, foundDb = os.LookupEnv(dbDatabaseKey)

	if !foundUser || !foundPass || !foundDb {
		return AppConfig{},
			fmt.Errorf("looking for database client configuration in .env file: %w", ErrEnvVarNotSet)
	}

	return cfg, nil
}
