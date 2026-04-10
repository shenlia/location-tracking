package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Tracking TrackingConfig `yaml:"tracking"`
	Geo      GeoConfig      `yaml:"geo"`
	Admin    AdminConfig    `yaml:"admin"`
}

type ServerConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	StaticDir string `yaml:"static_dir"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

type TrackingConfig struct {
	DefaultDomain     string `yaml:"default_domain"`
	CodeLength        int    `yaml:"code_length"`
	HeartbeatInterval int    `yaml:"heartbeat_interval"`
}

type GeoConfig struct {
	IPAPIURL   string `yaml:"ip_api_url"`
	IPFallback bool   `yaml:"ip_fallback"`
}

type AdminConfig struct {
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	SessionSecret string `yaml:"session_secret"`
}

var AppConfig *Config

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	AppConfig = &cfg
	return &cfg, nil
}

func Get() *Config {
	return AppConfig
}
