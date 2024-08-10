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

package model

import "time"

type (
	CustomerRequest struct {
		ID             string `json:"id" validate:"omitempty,numeric"`
		FirstName      string `json:"first_name" validate:"required,min=2,max=100"`
		LastName       string `json:"last_name" validate:"required,min=2,max=100"`
		Email          string `json:"email" validate:"required,email"`
		Phone          string `json:"phone" validate:"required,e164"`
		Address        string `json:"address" validate:"required"`
		HashedPassword string `json:"hashed_password" validate:"required"`
	}

	CustomerResponse struct {
		ID        string    `json:"id"`
		FullName  string    `json:"full_name"`
		Email     string    `json:"email"`
		Phone     string    `json:"phone"`
		Address   string    `json:"address"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)
