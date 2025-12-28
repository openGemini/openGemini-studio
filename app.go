package main

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	bolt "go.etcd.io/bbolt"
)

// App struct
type App struct {
	ctx      context.Context
	db       *bolt.DB
	connects sync.Map
	logger   *Logger
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
}

func (app *App) shutdown(ctx context.Context) {
	if err := app.db.Close(); err != nil {
		app.logger.Error("close database failed", "reason", err)
	}

	app.logger.Close()
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
