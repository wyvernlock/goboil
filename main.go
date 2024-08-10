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

package main

import (
	"goboil/internal/database"
	"goboil/internal/handler"
	"goboil/internal/repository"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	defaultAddr = ":8080"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	dbConfig := database.Config{
		MaxConns: 10,
		MinConns: 2,
	}

	db, err := database.NewPool(logger, dbConfig)
	if err != nil {
		logger.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	repo := repository.NewPsqlCustomerRepo(db, logger)
	customerHandler := handler.NewCustomerHandler(repo, logger)

	mux := http.NewServeMux()
	srv := http.Server{
		Addr:         defaultAddr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	mux.HandleFunc("GET /customers", customerHandler.GetAll)
	logger.Info("Route registered", "method", "GET", "path", "/customers")

	logger.Info("Server starting", "address", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Failed to start HTTP server", "error", err)
		os.Exit(1)
	}

	startTime := time.Now()
	logger.Info("HTTP server stopped", "uptime", time.Since(startTime).String())
}
