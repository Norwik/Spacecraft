// Copyright 2021 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/clevenio/spacecraft/orders/core/service"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

var token *service.Token

func init() {
	token = &service.Token{}
}

// Order controller
func Order(c echo.Context) error {
	var err error
	var t *service.Token

	id, _ := strconv.Atoi(c.Param("id"))

	helmet := service.NewHelmetClient(viper.GetString("app.apigw.url"), 20)

	if token.IsExpired() {
		t, err = helmet.FetchAccessToken(context.Background())
		token.Update(t.AccessToken, t.TokenType, t.ExpireAt)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Fetch customer data
	customer, err := helmet.GetCustomerById(id, token.AccessToken)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Return order data
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":       1,
		"customer": customer,
	})
}
