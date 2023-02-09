// Copyright 2022 Datafuse Labs.
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

package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/databendcloud/bendsql/internal/config"
	dc "github.com/databendcloud/databend-go"
	"github.com/pkg/errors"
)

type Client struct {
	cfg *config.CloudConfig
}

const (
	accept          = "Accept"
	authorization   = "Authorization"
	contentType     = "Content-Type"
	jsonContentType = "application/json; charset=utf-8"
	timeZone        = "Time-Zone"
	userAgent       = "User-Agent"

	EndpointGlobal = "https://app.databend.com"
	EndpointCN     = "https://app.databend.cn"
)

func NewClient() (*Client, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config")
	}

	if cfg.Cloud == nil {
		cfg.Cloud = &config.CloudConfig{
			Endpoint: EndpointGlobal,
		}
	}

	client := &Client{
		cfg: cfg.Cloud,
	}
	return client, nil
}

func (c *Client) WriteConfig() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return errors.Wrap(err, "failed to get config")
	}
	cfg.Target = config.TARGET_CLOUD
	cfg.Cloud = c.cfg
	return config.FlushConfig(cfg)
}

func (c *Client) CurrentWarehouse() string {
	return c.cfg.Warehouse
}

func (c *Client) CurrentOrganization() string {
	return c.cfg.Org
}

func (c *Client) CurrentEndpoint() string {
	return c.cfg.Endpoint
}

func (c *Client) SetCurrentWarehouse(warehouse string) error {
	warehouseList, err := c.ListWarehouses()
	if err != nil {
		return errors.Wrap(err, "failed to list warehouses")
	}
	if len(warehouseList) == 0 {
		return errors.New("no warehouse found")
	}
	if warehouse == "" {
		c.cfg.Warehouse = warehouseList[0].Name
		return nil
	}
	for i := range warehouseList {
		if warehouse == warehouseList[i].Name {
			c.cfg.Warehouse = warehouse
			return nil
		}
	}
	return errors.Errorf("warehouse %s not found", warehouse)
}

func (c *Client) SetEndpoint(endpoint string) {
	c.cfg.Endpoint = endpoint
}

func (c *Client) SetCurrentOrg(org, tenant, gateway string) {
	c.cfg.Org = org
	c.cfg.Tenant = tenant
	c.cfg.Gateway = gateway
}

func (c *Client) DoRequest(method, path string, headers http.Header, req interface{}, resp interface{}) error {
	if c.cfg.Token == nil {
		return errors.New("please use `bendsql cloud login` to login your account first")
	}
	if c.cfg.Token.ExpiresAt.Before(time.Now()) {
		err := c.RefreshToken()
		if err != nil {
			return errors.Wrap(err, "failed to refresh token")
		}
	}
	if headers != nil {
		headers = headers.Clone()
	} else {
		headers = http.Header{}
	}
	headers.Set(authorization, "Bearer "+c.cfg.Token.AccessToken)
	return c.request(method, path, headers, req, resp)
}

func (c *Client) request(method, path string, headers http.Header, req interface{}, resp interface{}) error {
	var err error

	reqBody := []byte{}
	if req != nil {
		reqBody, err = json.Marshal(req)
		if err != nil {
			return errors.Wrap(err, "failed to marshal request body")
		}
	}

	url, err := c.makeURL(path)
	if err != nil {
		return errors.Wrap(err, "failed to make url")
	}
	httpReq, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return errors.Wrap(err, "failed to create http request")
	}

	httpReq.Header = headers
	httpReq.Header.Set(contentType, jsonContentType)
	httpReq.Header.Set(accept, jsonContentType)

	httpClient := &http.Client{}
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return errors.Wrap(err, "http request error")
	}
	defer httpResp.Body.Close()

	httpRespBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read http response body")
	}

	if httpResp.StatusCode == http.StatusUnauthorized {
		return dc.NewAPIError("please use `bendsql cloud login` to login your account.", httpResp.StatusCode, httpRespBody)
	} else if httpResp.StatusCode >= 500 {
		return dc.NewAPIError("please retry again later.", httpResp.StatusCode, httpRespBody)
	} else if httpResp.StatusCode >= 400 {
		return dc.NewAPIError("please check your arguments.", httpResp.StatusCode, httpRespBody)
	}

	if resp != nil {
		if err := json.Unmarshal(httpRespBody, &resp); err != nil {
			return errors.Wrap(err, "failed to unmarshal http response body")
		}
	}

	return nil
}

func (c *Client) makeURL(path string) (string, error) {
	apiEndpoint := os.Getenv("BENDSQL_API_ENDPOINT")
	if apiEndpoint == "" {
		apiEndpoint = c.cfg.Endpoint
	}
	if apiEndpoint == "" {
		apiEndpoint = EndpointGlobal
	}
	u, err := url.Parse(apiEndpoint)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse api endpoint")
	}
	u.Path = path
	return u.String(), nil
}

func (c *Client) GetCloudDSN() (dsn string, err error) {
	if c.cfg.Token == nil {
		return "", errors.New("please use `bendsql cloud login` to login your account first")
	}

	cfg := dc.NewConfig()
	if strings.HasPrefix(c.cfg.Endpoint, "http://") {
		cfg.SSLMode = dc.SSL_MODE_DISABLE
	}
	cfg.Host = c.cfg.Gateway
	cfg.Tenant = c.cfg.Tenant
	cfg.Warehouse = c.cfg.Warehouse
	cfg.AccessToken = c.cfg.Token.AccessToken

	dsn = cfg.FormatDSN()
	return
}
