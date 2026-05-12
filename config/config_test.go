package config

import (
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}
	if cfg.Database.Driver != "sqlite3" {
		t.Errorf("DefaultConfig().Database.Driver = %v, want sqlite3", cfg.Database.Driver)
	}
	if cfg.Database.Name != "fate" {
		t.Errorf("DefaultConfig().Database.Name = %v, want fate", cfg.Database.Name)
	}
}

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if cfg == nil {
		t.Fatal("LoadConfig() returned nil")
	}

	def := DefaultConfig()
	if cfg.Database.Driver != def.Database.Driver {
		t.Errorf("LoadConfig().Database.Driver = %v, want %v", cfg.Database.Driver, def.Database.Driver)
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name:    "valid default",
			cfg:     DefaultConfig(),
			wantErr: false,
		},
		{
			name:    "nil config",
			cfg:     nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
