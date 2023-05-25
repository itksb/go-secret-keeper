package client

import (
	"errors"
	"os"
	"testing"
)

// TestConfig_UseEnvVars tests the UseEnvVars method
func TestConfig_UseJsonConfigFile(t *testing.T) {
	// Create a temporary config file
	tempFile, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// Write the JSON content to the temporary file
	jsonContent := `{
		"grpc_address": ":5000"
	}`
	err = os.WriteFile(tempFile.Name(), []byte(jsonContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	cfg := &Config{
		GRPCAddr: "",
		Config:   tempFile.Name(),
	}

	// Use JSON config file
	_, err = cfg.UseJsonConfigFile()
	if err != nil {
		t.Errorf("Error while using JSON config file: %v", err)
	}

	// Check if values are correctly updated
	expectedGRPCAddr := ":5000"
	if cfg.GRPCAddr != expectedGRPCAddr {
		t.Errorf("Expected GRPCAddr to be %s, but got %s", expectedGRPCAddr, cfg.GRPCAddr)
	}

	// Test error handling when config file doesn't exist
	cfg = &Config{
		GRPCAddr: ":3200",
		Config:   "nonexistent.json",
	}

	_, err = cfg.UseJsonConfigFile()
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected error os.ErrNotExist, but got %v", err)
	}
}

// TestMergeConfigs tests the mergeConfigs function
func TestMergeConfigs(t *testing.T) {
	result := &Config{
		GRPCAddr: "",
		Config:   "",
	}
	cfg2 := &Config{
		GRPCAddr: ":5000",
		Config:   "",
	}

	// Test merging when result has empty values
	merged := mergeConfigs(result, cfg2)

	// Check if values are correctly merged
	expectedGRPCAddr := ":5000"
	if merged.GRPCAddr != expectedGRPCAddr {
		t.Errorf("Expected GRPCAddr to be %s, but got %s", expectedGRPCAddr, merged.GRPCAddr)
	}

	// Test merging when result has non-empty values
	result.GRPCAddr = ":3200"
	merged = mergeConfigs(result, cfg2)

	// Check if values are correctly merged
	if merged.GRPCAddr != ":3200" {
		t.Errorf("Expected GRPCAddr to be :3200, but got %s", merged.GRPCAddr)
	}

}

// TestNewConfig tests the NewConfig function
func TestNewConfig(t *testing.T) {
	config := NewConfig()

	// Check initial values
	expectedGRPCAddr := ""
	if config.GRPCAddr != expectedGRPCAddr {
		t.Errorf("Expected GRPCAddr to be %s, but got %s", expectedGRPCAddr, config.GRPCAddr)
	}

	expectedConfig := ""
	if config.Config != expectedConfig {
		t.Errorf("Expected Config to be %s, but got %s", expectedConfig, config.Config)
	}
}
