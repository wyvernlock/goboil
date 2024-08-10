/*
Goboil
Boilerplate for building Go REST API with net/http

Copyright (C) 2024 Rian

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	MaxConns int32
	MinConns int32
}

func NewPool(logger *slog.Logger, cfg Config) (*pgxpool.Pool, error) {
	logger.Info("Initializing database connection")

	dbURL, err := getDBURL()
	if err != nil {
		logger.Error("Failed to get DATABASE_URL", "error", err)
		return nil, err
	}

	poolConfig, err := createPoolConfig(dbURL, cfg)
	if err != nil {
		logger.Error("Failed to create pool config", "error", err)
		return nil, err
	}

	db, err := createPool(poolConfig)
	if err != nil {
		logger.Error("Failed to create database connection pool", "error", err)
		return nil, err
	}

	if err := pingDatabase(db); err != nil {
		logger.Error("Failed to ping database", "error", err)
		db.Close()
		return nil, err
	}

	logger.Info("Successfully connected to database",
		"max_conns", poolConfig.MaxConns,
		"min_conns", poolConfig.MinConns)
	return db, nil
}

func getDBURL() (string, error) {
	dbURL, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return "", errors.New("DATABASE_URL environment variable is not set")
	}
	return dbURL, nil
}

func createPoolConfig(dbURL string, cfg Config) (*pgxpool.Config, error) {
	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}
	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	return poolConfig, nil
}

func createPool(config *pgxpool.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return pgxpool.NewWithConfig(ctx, config)
}

func pingDatabase(db *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.Ping(ctx)
}
