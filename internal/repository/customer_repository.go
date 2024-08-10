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

package repository

import (
	"context"
	"goboil/internal/model"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerRepository interface {
	// Create(ctx context.Context, cust *model.CustomerRequest) error
	GetAll(ctx context.Context) ([]*model.CustomerResponse, error)
	// GetByID(ctx context.Context, ID string) (*model.CustomerResponse, error)
	// Update(ctx context.Context, cust *model.CustomerRequest) error
	// Delete(ctx context.Context, ID string) error
}

type psqlCustomerRepo struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewPsqlCustomerRepo(db *pgxpool.Pool, logger *slog.Logger) *psqlCustomerRepo {
	return &psqlCustomerRepo{
		db:     db,
		logger: logger,
	}
}

func (p *psqlCustomerRepo) GetAll(ctx context.Context) ([]*model.CustomerResponse, error) {
	query := `
        SELECT 
            id, 
            CONCAT(first_name, ' ', last_name) AS full_name, 
            email, 
            phone, 
            address, 
            created_at AT TIME ZONE 'UTC', 
            updated_at AT TIME ZONE 'UTC'
        FROM customers
        ORDER BY created_at DESC
    `
	rows, err := p.db.Query(ctx, query)
	if err != nil {
		p.logger.Error("Failed to execute query", slog.String("query", query), slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var customers []*model.CustomerResponse
	for rows.Next() {
		cust := &model.CustomerResponse{}
		if err := rows.Scan(
			&cust.ID,
			&cust.FullName,
			&cust.Email,
			&cust.Phone,
			&cust.Address,
			&cust.CreatedAt,
			&cust.UpdatedAt,
		); err != nil {
			p.logger.Error("Failed to scan row", slog.String("error", err.Error()))
			return nil, err
		}
		customers = append(customers, cust)
	}

	if err := rows.Err(); err != nil {
		p.logger.Error("Row iteration error", slog.String("error", err.Error()))
		return nil, err
	}

	p.logger.Info("Successfully fetched all customers", slog.Int("count", len(customers)))
	return customers, nil
}
