// Copyright 2021 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Customer controller
func Customer(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":      id,
		"name":    "Joe Doe",
		"address": "123BT",
	})
}
