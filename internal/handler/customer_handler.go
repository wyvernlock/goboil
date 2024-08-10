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

package handler

import (
	"encoding/json"
	"goboil/internal/repository"
	"log/slog"
	"net/http"
)

type CustomerHandler struct {
	repo   repository.CustomerRepository
	logger *slog.Logger
}

func NewCustomerHandler(repo repository.CustomerRepository, logger *slog.Logger) *CustomerHandler {
	return &CustomerHandler{
		repo:   repo,
		logger: logger,
	}
}

func (h *CustomerHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	customers, err := h.repo.GetAll(r.Context())
	if err != nil {
		h.logger.Error("Failed to fetch customers", slog.String("error", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("Fetched all customers", slog.Int("count", len(customers)))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}
