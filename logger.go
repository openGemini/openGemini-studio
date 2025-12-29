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
	"log/slog"
	"os"
	"path/filepath"

	slogmulti "github.com/samber/slog-multi"
)

type Logger struct {
	logger        *slog.Logger
	defaultLogger *slog.Logger
	fileHandler   *os.File
}

func NewLogger() *Logger {
	var log = new(Logger)
	log.defaultLogger = slog.Default()

	logName := filepath.Join(workDirectory, "app.log")
	handle, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.defaultLogger.Error("open log file failed, use default logger", "reason", err)
		return log
	}
	log.fileHandler = handle
	defaultOpts := &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}
	log.logger = slog.New(slogmulti.Fanout(
		slog.NewTextHandler(os.Stdout, defaultOpts),
		slog.NewTextHandler(handle, defaultOpts),
	))

	slog.SetDefault(log.logger)

	return log
}

func (log *Logger) Close() {
	if log.fileHandler == nil {
		return
	}
	_ = log.fileHandler.Close()
}

func (log *Logger) Info(msg string, args ...interface{}) {
	slog.Info(msg, args...)
}

func (log *Logger) Warn(msg string, args ...interface{}) {
	slog.Warn(msg, args...)
}

func (log *Logger) Error(msg string, args ...interface{}) {
	slog.Error(msg, args...)
}
