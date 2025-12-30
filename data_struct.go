// Copyright 2025 openGemini Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

type RetentionPolicy struct {
	Name     string `json:"name"`
	Duration string `json:"duration"`
}

type DatabaseMetadata struct {
	RetentionPolicy []*RetentionPolicy `json:"retention_policies"`
	Measurements    []string           `yaml:"measurements"`
}
