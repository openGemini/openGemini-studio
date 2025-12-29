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
	"path/filepath"
	"time"

	bolt "go.etcd.io/bbolt"
)

const (
	BucketConnections = "connections"
	BucketSettings    = "settings"
	BucketHistories   = "histories"
)

func ConnectDatabase() (*bolt.DB, error) {
	db, err := bolt.Open(filepath.Join(workDirectory, "config.db"), 0600, &bolt.Options{
		Timeout: time.Second * 3,
		Logger:  nil,
	})
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(BucketConnections)); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(BucketSettings)); err != nil {
			return err
		}
		settingsBucket := tx.Bucket([]byte(BucketSettings))
		if settingsBucket.Get([]byte("system")) == nil {
			err := settingsBucket.Put([]byte("system"), defaultAppSetting.Marshal())
			if err != nil {
				return err
			}
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(BucketHistories)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
