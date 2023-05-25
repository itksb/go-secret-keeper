package client

import (
	"encoding/json"
	"flag"
	"os"
)

// Config - server config
type Config struct {
	GRPCAddr  string `json:"grpc_address"` // gRPC server address
	Config    string `json:"-"`            // config file path
	Debug     bool   `json:"-"`            // is debug mode
	CryptoKey string `json:"crypto_key"`   // crypto key
}

// NewConfig - create new config
func NewConfig() *Config {
	return &Config{
		GRPCAddr:  "",
		Config:    "",
		Debug:     false,
		CryptoKey: "",
	}
}

// UseFlags - use flags
func (cfg *Config) UseFlags() {
	grpcAddr := flag.String("g", cfg.GRPCAddr, "GRPC_ADDRESS")
	configFile := flag.String("c", cfg.Config, "CONFIG")
	cryptoKey := flag.String("k", cfg.CryptoKey, "CRYPTO_KEY")

	flag.Parse()

	if *grpcAddr != "" {
		cfg.GRPCAddr = *grpcAddr
	}
	if *configFile != "" {
		cfg.Config = *configFile
	}
	if *cryptoKey != "" {
		cfg.CryptoKey = *cryptoKey
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
		mergeConfigs(
			cfg,
			&config,
		)
	}

	return cfg, nil
}

// mergeConfigs merges configs into one
// first config values have priority
func mergeConfigs(result, cfg2 *Config) *Config {
	if result.GRPCAddr == "" {
		result.GRPCAddr = cfg2.GRPCAddr
	}
	if result.CryptoKey == "" {
		result.CryptoKey = cfg2.CryptoKey
	}
	if !result.Debug {
		result.Debug = cfg2.Debug
	}
	return result
}
