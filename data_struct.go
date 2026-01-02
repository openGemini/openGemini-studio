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
	debug             bool   `json:"-"`
}

var defaultAppSetting = &AppSetting{
	Language:        "en",
	ThemeMode:       "light",
	MaxHistoryCount: 100,
	DataDirectory:   "./data",
	Debug:           false,
}

type AppSetting struct {
	Language        string `json:"language"`
	ThemeMode       string `json:"theme_mode"`
	CustomFont      string `json:"custom_font"`
	MaxHistoryCount int    `json:"max_history_count"`
	DataDirectory   string `json:"data_dir"`
	Debug           bool   `json:"debug"`
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
	Measurements    []string           `json:"measurements"`
}

type ExecuteRequest struct {
	ConnectName     string `json:"connect_name"`
	Database        string `json:"database"`
	RetentionPolicy string `json:"retention_policy"`
	Measurement     string `json:"measurement"`
	Precision       string `json:"precision"`
	Command         string `json:"command"`
}

type ExecuteResponse struct {
	NoContent     bool     `json:"no_content"`
	Message       string   `json:"message"`
	ExecutionTime float64  `json:"execution_time"` // Execution time in milliseconds
	Columns       []string `json:"columns"`
	Values        [][]any  `json:"values"`
}

type History struct {
	ID              string  `json:"id"`
	Query           string  `json:"query"`
	Timestamp       int64   `json:"timestamp"`      // Unix timestamp in milliseconds
	ExecutionTime   float64 `json:"execution_time"` // Execution time in milliseconds
	Database        string  `json:"database"`
	RetentionPolicy string  `json:"retention_policy"`
	Success         bool    `json:"success"`
	Error           string  `json:"error"`
}

func (h *History) Marshal() []byte {
	data, err := json.Marshal(h)
	if err != nil {
		return []byte("{}")
	}
	return data
}
