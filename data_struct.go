package main

import (
	"encoding/json"
)

type ConnectConfig struct {
	Name              string `json:"name"`
	Address           string `json:"address"`
	HTTPSchema        string `json:"http_schema"`
	CACertificate     string `json:"ca_certificate"`
	ClientCertificate string `json:"client_certificate"`
	ClientKey         string `json:"client_key"`
	InsecureTls       bool   `json:"insecure_tls"`
	InsecureHostname  bool   `json:"insecure_hostname"`
	EnableAuth        bool   `json:"enable_auth"`
	Username          string `json:"username"`
	Password          string `json:"password"`
}

var defaultAppSetting = &AppSetting{
	Language:        "en",
	ThemeMode:       "light",
	MaxHistoryCount: 100,
	DataDirectory:   "./data",
}

type AppSetting struct {
	Language        string `json:"language"`
	ThemeMode       string `json:"theme_mode"`
	CustomFont      string `json:"custom_font"`
	MaxHistoryCount int    `json:"max_history_count"`
	DataDirectory   string `json:"data_dir"`
}

func (as *AppSetting) Marshal() []byte {
	data, err := json.Marshal(as)
	if err != nil {
		return []byte("{}")
	}
	return data
}
