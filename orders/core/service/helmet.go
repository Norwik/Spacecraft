// Copyright 2021 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

var (
	OauthURI                 = "apigw/token"
	CustomersServiceEndpoint = "customers/v1/customer"
)

type Customer struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Token struct {
	sync.RWMutex

	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpireAt    time.Time `json:"expire_at"`
}

// LoadFromJSON update object from json
func (c *Customer) LoadFromJSON(data []byte) error {
	err := json.Unmarshal(data, &c)
	if err != nil {
		return err
	}
	return nil
}

// ConvertToJSON convert object to json
func (c *Customer) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&c)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// LoadFromJSON update object from json
func (t *Token) LoadFromJSON(data []byte) error {
	err := json.Unmarshal(data, &t)
	if err != nil {
		return err
	}
	return nil
}

// ConvertToJSON convert object to json
func (t *Token) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&t)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// IsExpired ..
func (t *Token) IsExpired() bool {
	if t.AccessToken == "" || t.ExpireAt.Before(time.Now()) {
		return true
	}

	return false
}

// Update ..
func (t *Token) Update(accessToken, tokenType string, expireAt time.Time) {
	t.Lock()
	defer t.Unlock()

	t.TokenType = tokenType
	t.AccessToken = accessToken
	t.ExpireAt = expireAt
}

// Helmet struct
type Helmet struct {
	Timeout time.Duration
	APIGW   string
}

// NewHelmetClient creates an instance of http client
func NewHelmetClient(apigwURL string, timeout int) *Helmet {
	return &Helmet{
		Timeout: time.Duration(timeout),
		APIGW:   apigwURL,
	}
}

// FetchAccessToken ..
func (h *Helmet) FetchAccessToken(ctx context.Context) (*Token, error) {
	var err error

	token := &Token{}
	params := url.Values{}
	params.Add("grant_type", `client_credentials`)
	reqBody := strings.NewReader(params.Encode())

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/%s", h.APIGW, OauthURI),
		reqBody,
	)

	if err != nil {
		return token, err
	}

	req = req.WithContext(ctx)

	req.SetBasicAuth(
		viper.GetString("app.apigw.client_id"),
		viper.GetString("app.apigw.client_secret"),
	)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{
		Timeout: time.Second * h.Timeout,
	}

	resp, err := client.Do(req)

	if err != nil {
		return token, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return token, err
	}

	err = token.LoadFromJSON(body)

	token.ExpireAt = time.Now().Add(time.Second * 3600)

	if err != nil {
		return token, err
	}

	return token, nil
}

// GetCustomerById ..
func (h *Helmet) GetCustomerById(id int, accessToken string) (Customer, error) {
	var err error
	var customer Customer

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/%s/%d", h.APIGW, CustomersServiceEndpoint, id),
		nil,
	)

	if err != nil {
		return customer, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := http.Client{
		Timeout: time.Second * h.Timeout,
	}

	resp, err := client.Do(req)

	if err != nil {
		return customer, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return customer, err
	}

	err = customer.LoadFromJSON(body)

	if err != nil {
		return customer, err
	}

	return customer, nil
}
