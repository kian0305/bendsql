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
	"net/http"
	"strings"
	"time"

	"github.com/avast/retry-go"
	dc "github.com/databendcloud/databend-go"
	"github.com/pkg/errors"
)

func (c *Client) Query(warehouseName, query string) (*dc.QueryResponse, error) {
	headers := make(http.Header)
	headers.Set("X-DATABENDCLOUD-WAREHOUSE", warehouseName)
	headers.Set("X-DATABENDCLOUD-ORG", string(c.cfg.Org))
	request := dc.QueryRequest{
		SQL: query,
	}
	path := "/v1/query"
	var result dc.QueryResponse
	err := c.DoRequest("POST", path, headers, request, &result)
	if err != nil {
		return nil, err
	}
	if result.Error != nil {
		return &result, errors.Wrapf(result.Error, "query %s in org %s: %s", warehouseName, c.cfg.Org)
	}
	return &result, nil
}

func (c *Client) QuerySync(warehouseName string, sql string, respCh chan dc.QueryResponse) error {
	var r0 *dc.QueryResponse
	err := retry.Do(
		func() error {
			r, err := c.Query(warehouseName, sql)
			if err != nil {
				return errors.Wrap(err, "query failed")
			}
			r0 = r
			return nil
		},
		// other err no need to retry
		retry.RetryIf(func(err error) bool {
			if err != nil && !(dc.IsProxyErr(err) || strings.Contains(err.Error(), dc.ProvisionWarehouseTimeout)) {
				return false
			}
			return true
		}),
		retry.Delay(2*time.Second),
		retry.Attempts(10),
	)
	if err != nil {
		return errors.Wrap(err, "query failed")
	}
	if r0.Error != nil {
		return errors.Wrap(r0.Error, "query has error")
	}
	if err != nil {
		return err
	}
	respCh <- *r0
	nextUri := r0.NextURI
	for len(nextUri) != 0 {
		p, err := c.QueryPage(warehouseName, r0.Id, nextUri)
		if err != nil {
			return errors.Wrap(err, "query page failed")
		}
		if p.Error != nil {
			return errors.Wrap(p.Error, "query has error")
		}
		nextUri = p.NextURI
		respCh <- *p
	}
	return nil
}

func (c *Client) QueryPage(warehouseName, queryId, path string) (*dc.QueryResponse, error) {
	headers := make(http.Header)
	headers.Set("queryID", queryId)
	headers.Set("X-DATABENDCLOUD-WAREHOUSE", warehouseName)
	headers.Set("X-DATABENDCLOUD-ORG", string(c.cfg.Org))
	var result dc.QueryResponse
	err := retry.Do(
		func() error {
			err := c.DoRequest("GET", path, headers, nil, &result)
			if err != nil {
				return errors.Wrap(err, "query page failed")
			}
			return nil
		},
		retry.Delay(2*time.Second),
		retry.Attempts(10),
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
