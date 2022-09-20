package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/datafuselabs/bendcloud-cli/api/apierrors"
	"github.com/datafuselabs/bendcloud-cli/internal/config"
)

type APIClient struct {
	UserEmail        string
	Password         string
	AccessToken      string
	RefreshToken     string
	CurrentOrgSlug   string
	CurrentWarehouse string
}

const (
	accept             = "Accept"
	authorization      = "Authorization"
	contentType        = "Content-Type"
	jsonContentType    = "application/json; charset=utf-8"
	timeZone           = "Time-Zone"
	userAgent          = "User-Agent"
	defaultApiEndpoint = "https://app.databend.com"
)

func NewApiClient() *APIClient {
	accessToken, refreshToken := config.GetAuthToken()
	return &APIClient{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		CurrentOrgSlug:   config.GetOrg(),
		CurrentWarehouse: config.GetWarehouse(),
	}
}

func (c *APIClient) DoRequest(method, path string, headers http.Header, req interface{}, resp interface{}) error {
	var err error

	reqBody := []byte{}
	if req != nil {
		reqBody, err = json.Marshal(req)
		if err != nil {
			panic(err)
		}
	}

	url := c.makeURL(path)
	httpReq, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	if headers != nil {
		httpReq.Header = headers.Clone()
	}
	httpReq.Header.Set(contentType, jsonContentType)
	httpReq.Header.Set(accept, jsonContentType)
	if len(c.AccessToken) > 0 {
		httpReq.Header.Set(authorization, "Bearer "+c.AccessToken)
	}

	httpClient := &http.Client{}
	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed http do request: %w", err)
	}
	defer httpResp.Body.Close()

	httpRespBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("io read error: %w", err)
	}

	if httpResp.StatusCode == http.StatusUnauthorized {
		return apierrors.New("please use `bendctl auth login` to login your account.", httpResp.StatusCode, httpRespBody)
	} else if httpResp.StatusCode >= 500 {
		return apierrors.New("please retry again later.", httpResp.StatusCode, httpRespBody)
	} else if httpResp.StatusCode >= 400 {
		return apierrors.New("please check your arguments.", httpResp.StatusCode, httpRespBody)
	}

	if resp != nil {
		if err := json.Unmarshal(httpRespBody, &resp); err != nil {
			return err
		}
	}

	//if method != "GET" {
	//	respBody := string(httpRespBody)
	//	if strings.Contains(respBody, "Token") || strings.Contains(respBody, "secret") {
	//		log.Printf("webapi.doRequest %s url=%s req=nodata resp=nodata", method, url)
	//	}
	//}
	return nil
}

func (c *APIClient) makeURL(path string, args ...interface{}) string {
	apiEndpoint := os.Getenv("BENDCTL_API_ENDPOINT")
	if apiEndpoint == "" {
		apiEndpoint = defaultApiEndpoint
	}
	format := apiEndpoint + path
	return fmt.Sprintf(format, args...)
}
