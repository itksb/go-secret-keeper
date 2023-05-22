package server

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
		"grpc_address": ":5000",
		"database_dsn": "host=localhost port=5432 user=postgres password=pass dbname=mydb sslmode=disable"
	}`
	err = os.WriteFile(tempFile.Name(), []byte(jsonContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	cfg := &Config{
		GRPCAddr: "",
		Dsn:      "",
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

	expectedDsn := "host=localhost port=5432 user=postgres password=pass dbname=mydb sslmode=disable"
	if cfg.Dsn != expectedDsn {
		t.Errorf("Expected Dsn to be %s, but got %s", expectedDsn, cfg.Dsn)
	}

	// Test error handling when config file doesn't exist
	cfg = &Config{
		GRPCAddr: ":3200",
		Dsn:      "",
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
		Dsn:      "",
		Config:   "",
	}
	cfg2 := &Config{
		GRPCAddr: ":5000",
		Dsn:      "host=localhost port=5432 user=postgres password=pass dbname=mydb sslmode=disable",
		Config:   "",
	}

	// Test merging when result has empty values
	merged := mergeConfigs(result, cfg2)

	// Check if values are correctly merged
	expectedGRPCAddr := ":5000"
	if merged.GRPCAddr != expectedGRPCAddr {
		t.Errorf("Expected GRPCAddr to be %s, but got %s", expectedGRPCAddr, merged.GRPCAddr)
	}

	expectedDsn := "host=localhost port=5432 user=postgres password=pass dbname=mydb sslmode=disable"
	if merged.Dsn != expectedDsn {
		t.Errorf("Expected Dsn to be %s, but got %s", expectedDsn, merged.Dsn)
	}

	// Test merging when result has non-empty values
	result.GRPCAddr = ":3200"
	result.Dsn = "old-dsn"
	merged = mergeConfigs(result, cfg2)

	// Check if values are correctly merged
	if merged.GRPCAddr != ":3200" {
		t.Errorf("Expected GRPCAddr to be :3200, but got %s", merged.GRPCAddr)
	}

	if merged.Dsn != "old-dsn" {
		t.Errorf("Expected Dsn to be old-dsn, but got %s", merged.Dsn)
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

	expectedDsn := ""
	if config.Dsn != expectedDsn {
		t.Errorf("Expected Dsn to be %s, but got %s", expectedDsn, config.Dsn)
	}

	expectedConfig := ""
	if config.Config != expectedConfig {
		t.Errorf("Expected Config to be %s, but got %s", expectedConfig, config.Config)
	}
}
