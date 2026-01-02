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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/openGemini/opengemini-client-go/opengemini"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	bolt "go.etcd.io/bbolt"
)

// App struct
type App struct {
	ctx      context.Context
	db       *bolt.DB
	connects sync.Map
	logger   *Logger
	debug    bool
}

// NewApp creates a new App application struct
func NewApp() *App {
	var app = &App{}

	if err := CreateWorkDirectory(); err != nil {
		panic("create work directory failed: " + err.Error())
	}

	app.logger = NewLogger()

	database, err := ConnectDatabase()
	if err != nil {
		app.logger.Error("open database failed", "reason", err)
		panic("open database failed")
	}
	app.db = database

	return app
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (app *App) startup(ctx context.Context) {
	app.ctx = ctx
	setting, err := app.GetSetting()
	if err != nil {
		app.logger.Error("get setting failed", "reason", err)
	} else {
		app.debug = setting.Debug
	}
}

func (app *App) shutdown(ctx context.Context) {
	// Close all HTTP client connections first
	app.connects.Range(func(key, value interface{}) bool {
		app.connects.Delete(key)
		return true
	})

	// Close database connection
	if app.db != nil {
		if err := app.db.Close(); err != nil {
			app.logger.Error("close database failed", "reason", err)
		}
		app.db = nil
	}

	// Close logger last
	if app.logger != nil {
		app.logger.Close()
	}
}

func (app *App) ListConnects() []*ConnectConfig {
	var connects []*ConnectConfig
	err := app.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketConnections))
		err := bucket.ForEach(func(k, v []byte) error {
			var connect = &ConnectConfig{}
			err := json.Unmarshal(v, connect)
			if err != nil {
				return err
			}
			connects = append(connects, connect)
			return nil
		})
		return err
	})
	if err != nil {
		app.logger.Error("list connects failed", "reason", err)
		return connects
	}
	return connects
}

func (app *App) AddConnect(cc *ConnectConfig) error {
	err := app.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketConnections))
		existConnect := bucket.Get([]byte(cc.Name))
		if existConnect != nil {
			app.logger.Error("connect already exists", "connection", cc.Name)
			return errors.New("connect already exist")
		}
		content, err := json.Marshal(cc)
		if err != nil {
			return err
		}
		err = bucket.Put([]byte(cc.Name), content)
		return err
	})
	if err == nil {
		app.logger.Info("add connect", "name", cc.Name)
		return nil
	}
	return err
}

func (app *App) UpdateConnect(name string, cc *ConnectConfig) error {
	err := app.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketConnections))
		existConnect := bucket.Get([]byte(name))
		if existConnect == nil {
			app.logger.Error("connect does not exist", "name", name)
			return errors.New("connect does not exist")
		}
		content, err := json.Marshal(cc)
		if err != nil {
			return err
		}
		err = bucket.Put([]byte(name), content)
		return err
	})
	if err == nil {
		app.logger.Info("update connect", "name", name)
		return nil
	}
	return err
}

func (app *App) DeleteConnect(name string) error {
	err := app.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketConnections))
		return bucket.Delete([]byte(name))
	})
	return err
}

func (app *App) GetConnect(name string) (*ConnectConfig, error) {
	var connect = &ConnectConfig{}
	err := app.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketConnections))
		content := bucket.Get([]byte(name))
		if content == nil {
			return errors.New("connection not found")
		}
		err := json.Unmarshal(content, connect)
		return err
	})
	if err != nil {
		app.logger.Error("get connect failed", "reason", err)
		return connect, err
	}
	return connect, nil
}

func (app *App) UpdateSetting(settings *AppSetting) error {
	app.debug = settings.Debug
	data, err := json.Marshal(settings)
	if err != nil {
		app.logger.Error("update settings failed", "reason", err)
		return err
	}
	err = app.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketSettings))
		return bucket.Put([]byte("system"), data)
	})
	return err
}

func (app *App) GetSetting() (*AppSetting, error) {
	var setting = &AppSetting{}
	err := app.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketSettings))
		data := bucket.Get([]byte("system"))
		if data == nil {
			app.logger.Error("get settings failed", "reason", "setting not found")
			return errors.New("setting not found")
		}
		return json.Unmarshal(data, setting)
	})
	return setting, err
}

func (app *App) OpenFileDialog() (string, error) {
	filePath, err := runtime.OpenFileDialog(app.ctx, runtime.OpenDialogOptions{
		Title: "Select File",
	})
	return filePath, err
}

func (app *App) DialConnect(name string) ([]string, error) {
	app.connects.Delete(name)
	app.logger.Info("dial connect", "name", name)
	cc, err := app.GetConnect(name)
	if err != nil {
		app.logger.Error("dial connect failed: find connect config failed", "reason", err)
		return nil, err
	}
	if app.debug {
		cc.debug = true
	}
	httpClient, err := NewHttpClient(cc)
	if err != nil {
		app.logger.Error("dial connect failed: create http client failed", "reason", err)
		return nil, err
	}
	err = httpClient.Ping()
	if err != nil {
		app.logger.Error("dial connect failed: ping http client failed", "reason", err)
		return nil, err
	}
	app.connects.Store(name, httpClient)

	databases, err := httpClient.Databases(app.ctx)
	if err != nil {
		app.logger.Error("dial connect failed: databases failed", "reason", err)
		return nil, err
	}
	return databases, nil
}

func (app *App) getDialer(name string) (HttpClient, error) {
	value, ok := app.connects.Load(name)
	if !ok {
		return nil, errors.New("dialer not found")
	}
	return value.(HttpClient), nil
}

func (app *App) GetDatabaseMetadata(connectName, databaseName string) (*DatabaseMetadata, error) {
	httpClient, err := app.getDialer(connectName)
	if err != nil {
		app.logger.Error("get opengemini client failed", "reason", err, "name", connectName)
		return nil, err
	}
	policies, err := httpClient.RetentionPolicies(app.ctx, databaseName)
	if err != nil {
		app.logger.Error("get retention policies failed", "reason", err, "db", databaseName)
		return nil, err
	}
	measurements, err := httpClient.Measurements(app.ctx, databaseName)
	if err != nil {
		app.logger.Error("get measurements failed", "reason", err, "db", databaseName)
		return nil, err
	}

	return &DatabaseMetadata{RetentionPolicy: policies, Measurements: measurements}, nil
}

func (app *App) CloseConnect(connectName string) {
	app.logger.Info("close connect", "name", connectName)
	app.connects.Delete(connectName)
}

func (app *App) ExecuteCommand(data *ExecuteRequest) (*ExecuteResponse, error) {
	if data.ConnectName == "" {
		return nil, errors.New("connect name required")
	}
	if data.Command == "" {
		return nil, errors.New("command required")
	}
	httpClient, err := app.getDialer(data.ConnectName)
	if err != nil {
		app.logger.Error("get opengemini client failed", "reason", err, "command", data.Command)
		return nil, err
	}

	app.logger.Info("execute command data", "data", data)

	var startTime = time.Now()

	if strings.HasPrefix(strings.ToLower(data.Command), "insert") {
		if data.Database == "" {
			return nil, errors.New("database required")
		}
		if data.RetentionPolicy == "" {
			data.RetentionPolicy = "autogen"
		}
		if data.Precision == "" {
			data.Precision = "ns"
		}
		lineProtocol := strings.TrimSpace(data.Command[6:])
		err := httpClient.Write(app.ctx, data.Database, data.RetentionPolicy, lineProtocol, data.Precision)
		executionTime := time.Since(startTime).Milliseconds()

		// Save history record
		history := &History{
			ID:              strconv.FormatInt(time.Now().UnixMilli(), 10),
			Query:           data.Command,
			Timestamp:       time.Now().UnixMilli(),
			ExecutionTime:   float64(executionTime),
			Database:        data.Database,
			RetentionPolicy: data.RetentionPolicy,
			Success:         err == nil,
		}
		if err != nil {
			history.Error = err.Error()
			app.logger.Error("execute command failed", "reason", err, "command", data.Command)
			_ = app.AddHistory(history)
			return nil, err
		}
		_ = app.AddHistory(history)
		return &ExecuteResponse{NoContent: true, Message: "write success", ExecutionTime: float64(executionTime)}, nil
	}
	response, err := httpClient.Query(app.ctx, &opengemini.Query{
		Database:        data.Database,
		Command:         data.Command,
		RetentionPolicy: data.RetentionPolicy,
		Precision:       opengemini.ToPrecision(data.Precision),
	})
	executionTime := time.Since(startTime).Milliseconds()

	// Save history record
	history := &History{
		ID:              strconv.FormatInt(time.Now().UnixMilli(), 10),
		Query:           data.Command,
		Timestamp:       time.Now().UnixMilli(),
		ExecutionTime:   float64(executionTime),
		Database:        data.Database,
		RetentionPolicy: data.RetentionPolicy,
		Success:         false,
	}

	if err != nil {
		history.Error = err.Error()
		app.logger.Error("execute command failed", "reason", err, "name", data.ConnectName)
		_ = app.AddHistory(history)
		return nil, err
	}
	if response.Error != "" {
		history.Error = response.Error
		app.logger.Error("execute command failed", "reason", response.Error)
		_ = app.AddHistory(history)
		return nil, fmt.Errorf("execute command failed: %s", response.Error)
	}
	if len(response.Results) == 0 {
		history.Success = true
		_ = app.AddHistory(history)
		return &ExecuteResponse{NoContent: true, ExecutionTime: float64(executionTime)}, nil
	}

	if response.Results[0].Error != "" {
		history.Error = response.Results[0].Error
		app.logger.Error("execute command failed", "reason", response.Results[0].Error)
		_ = app.AddHistory(history)
		return nil, fmt.Errorf("execute command failed: %s", response.Results[0].Error)
	}

	if len(response.Results[0].Series) == 0 {
		history.Success = true
		_ = app.AddHistory(history)
		return &ExecuteResponse{NoContent: true, ExecutionTime: float64(executionTime)}, nil
	}

	var (
		seriesValues = response.Results[0].Series[0].Values
		columns      = response.Results[0].Series[0].Columns
		values       = make([][]any, 0, len(seriesValues))
	)
	for _, v := range seriesValues {
		// Directly append the row data from seriesValues
		values = append(values, v)
	}

	history.Success = true
	_ = app.AddHistory(history)
	return &ExecuteResponse{NoContent: false, Columns: columns, Values: values, ExecutionTime: float64(executionTime)}, nil
}

func (app *App) GetHistories() ([]*History, error) {
	var histories []*History
	err := app.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketHistories))
		err := bucket.ForEach(func(k, v []byte) error {
			var history = &History{}
			err := json.Unmarshal(v, history)
			if err != nil {
				return err
			}
			histories = append(histories, history)
			return nil
		})
		return err
	})
	if err != nil {
		app.logger.Error("list histories failed", "reason", err)
		return nil, err
	}

	// Sort by timestamp descending (newest first)
	for i := 0; i < len(histories)-1; i++ {
		for j := i + 1; j < len(histories); j++ {
			if histories[i].Timestamp < histories[j].Timestamp {
				histories[i], histories[j] = histories[j], histories[i]
			}
		}
	}

	// Limit to maxHistoryCount from settings
	settings, err := app.GetSetting()
	if err == nil && settings.MaxHistoryCount > 0 && len(histories) > settings.MaxHistoryCount {
		histories = histories[:settings.MaxHistoryCount]
	}

	return histories, nil
}

func (app *App) AddHistory(history *History) error {
	data, err := json.Marshal(history)
	if err != nil {
		app.logger.Error("marshal history failed", "reason", err)
		return err
	}

	err = app.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketHistories))
		return bucket.Put([]byte(history.ID), data)
	})

	if err != nil {
		app.logger.Error("save history failed", "reason", err)
		return err
	}

	// Clean up old histories if exceeds max count
	settings, err := app.GetSetting()
	if err == nil && settings.MaxHistoryCount > 0 {
		histories, err := app.GetHistories()
		if err == nil && len(histories) > settings.MaxHistoryCount {
			// Delete oldest histories
			toDelete := histories[settings.MaxHistoryCount:]
			err = app.db.Update(func(tx *bolt.Tx) error {
				bucket := tx.Bucket([]byte(BucketHistories))
				for _, h := range toDelete {
					if err := bucket.Delete([]byte(h.ID)); err != nil {
						return err
					}
				}
				return nil
			})
			if err != nil {
				app.logger.Error("delete old histories failed", "reason", err)
			}
		}
	}

	return nil
}
