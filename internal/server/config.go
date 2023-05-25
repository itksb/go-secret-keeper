package server

import (
	"encoding/json"
	"flag"
	"os"
	"strings"
)

// Config - server config
type Config struct {
	GRPCAddr string `json:"grpc_address"` // gRPC server address
	Dsn      string `json:"database_dsn"` // data source name
	Config   string `json:"-"`            // config file path
	Debug    bool   `json:"-"`            // is debug mode
	TokenKey string `json:"token_key"`    // token key
}

// NewConfig - create new config
func NewConfig() *Config {
	return &Config{
		GRPCAddr: "",
		Dsn:      "",
		Config:   "",
		Debug:    false,
		TokenKey: "",
	}
}

// UseFlags - use flags
func (cfg *Config) UseFlags() {
	grpcAddr := flag.String("g", cfg.GRPCAddr, "GRPC_ADDRESS")
	dsn := flag.String("d", cfg.Dsn, "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable")
	configFile := flag.String("c", cfg.Config, "CONFIG")
	debugStr := flag.String("e", func() string {
		if cfg.Debug {
			return "debug"
		}
		return "prod"
	}(), "prod|debug")
	token := flag.String("t", cfg.TokenKey, "Token secret key")

	flag.Parse()

	cfg.Dsn = *dsn
	cfg.Debug = strings.ToLower(*debugStr) == "debug"

	if *grpcAddr != "" {
		cfg.GRPCAddr = *grpcAddr
	}
	if *configFile != "" {
		cfg.Config = *configFile
	}
	if *token != "" {
		cfg.TokenKey = *token
	}
}

// UseJsonConfigFile - use config file
func (cfg *Config) UseJsonConfigFile() (*Config, error) {
	if cfg.Config != "" {
		configFile, err := os.Open(cfg.Config)
		if err != nil {
			return cfg, err
		}
		defer configFile.Close()

		config := Config{}

		jsonParser := json.NewDecoder(configFile)
		jsonParser.DisallowUnknownFields() // отклонение неизвестных полей
		err = jsonParser.Decode(&config)
		if err != nil {
			return cfg, err
		}
		cfg = mergeConfigs(cfg, &config)
	}

	return cfg, nil
}

// mergeConfigs merges configs into one
// first config values have priority
func mergeConfigs(result, cfg2 *Config) *Config {
	if result.Dsn == "" {
		result.Dsn = cfg2.Dsn
	}
	if result.GRPCAddr == "" {
		result.GRPCAddr = cfg2.GRPCAddr
	}
	if !result.Debug {
		result.Debug = cfg2.Debug
	}
	if result.TokenKey == "" {
		result.TokenKey = cfg2.TokenKey
	}

	return result
}
