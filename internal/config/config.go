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

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	dc "github.com/databendcloud/databend-go"
	"github.com/pkg/errors"
)

const (
	TARGET_COMMUNITY = "community"
	TARGET_CLOUD     = "cloud"
)

var (
	configFile     = ""
	cloudTokenFile = ""
)

func init() {
	if a := os.Getenv("BENDSQL_CONFIG"); a != "" {
		configFile = a
	} else {
		d, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		configFile = filepath.Join(d, ".config", "bendsql", "config.toml")
	}

	if t := os.Getenv("BENDSQL_CLOUD_TOKEN"); t != "" {
		cloudTokenFile = t
	} else {
		d, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		cloudTokenFile = filepath.Join(d, ".config", "bendsql", "cloud-token")
	}

	if !exists(configFile) {
		fmt.Printf("config file %s not found, creating a new one\n", configFile)
		if !exists(filepath.Dir(configFile)) {
			err := os.MkdirAll(filepath.Dir(configFile), 0755)
			if err != nil {
				panic(err)
			}
		}
		f, err := os.Create(configFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()
	}
}

type Config struct {
	Target    string           `toml:"target"`
	Cloud     *CloudConfig     `toml:"cloud,omitempty"`
	Community *CommunityConfig `toml:"community,omitempty"`
}

func (c *Config) Clone() Config {
	cfg := Config{
		Target: c.Target,
	}
	if c.Cloud != nil {
		cfg.Cloud = c.Cloud.Clone()
	}
	if c.Community != nil {
		cfg.Community = c.Community.Clone()
	}
	return cfg
}

func (c *Config) GetDSN(opts RuntimeOptions) (string, error) {
	switch c.Target {
	case TARGET_COMMUNITY:
		if c.Community == nil {
			return "", errors.New("please use `bendsql connect` to connect to your instance first")
		}
		return c.Community.GetDSN(opts)
	case TARGET_CLOUD:
		if c.Cloud == nil {
			return "", errors.New("please use `bendsql cloud login` to connect to your account first")
		}
		return c.Cloud.GetDSN(opts)
	default:
		return "", errors.New("please use `bendsql connect` or `bendsql cloud login` to connect to your instance first")
	}
}

type CommunityConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
	SSL      bool   `toml:"ssl"`

	Options map[string]string `toml:"options"`
}

func (c *CommunityConfig) Clone() *CommunityConfig {
	cfg := *c
	cfg.Options = map[string]string{}
	for k, v := range c.Options {
		cfg.Options[k] = v
	}
	return &cfg
}

func (c *CommunityConfig) GetDSN(opts RuntimeOptions) (string, error) {
	cfg := dc.NewConfig()
	cfg.Host = fmt.Sprintf("%s:%d", c.Host, c.Port)
	if opts.Username != "" {
		cfg.User = opts.Username
	} else {
		cfg.User = c.User
	}
	if opts.Password != "" {
		cfg.Password = opts.Password
	} else {
		cfg.Password = c.Password
	}
	if opts.Database != "" {
		cfg.Database = opts.Database
	} else {
		cfg.Database = c.Database
	}

	if !c.SSL {
		cfg.SSLMode = dc.SSL_MODE_DISABLE
	}
	cfg.AddParams(c.Options)

	dsn := cfg.FormatDSN()
	return dsn, nil
}

type CloudConfig struct {
	Org       string `toml:"org"`
	Tenant    string `toml:"tenant"`
	Warehouse string `toml:"warehouse"`
	Gateway   string `toml:"gateway"`
	Endpoint  string `toml:"endpoint"`

	Token *CloudToken `toml:"token,omitempty"`
}

func (c *CloudConfig) Clone() *CloudConfig {
	cfg := *c
	if c.Token != nil {
		t := *c.Token
		cfg.Token = &t
	}
	return &cfg
}

func (c *CloudConfig) GetDSN(opts RuntimeOptions) (string, error) {
	if c.Token == nil {
		return "", errors.New("please use `bendsql cloud login` to login your account first")
	}
	if c.Gateway == "" || c.Tenant == "" || c.Warehouse == "" {
		return "", errors.New("please use `bendsql cloud configure` to select organization and warehouse first")
	}

	cfg := dc.NewConfig()
	if strings.HasPrefix(c.Endpoint, "http://") {
		cfg.SSLMode = dc.SSL_MODE_DISABLE
	}
	cfg.Host = c.Gateway
	cfg.Tenant = c.Tenant
	cfg.Warehouse = c.Warehouse

	if opts.Username != "" {
		cfg.User = opts.Username
		cfg.Password = opts.Password
	} else {
		cfg.AccessToken = c.Token.AccessToken
	}
	if opts.Database != "" {
		cfg.Database = opts.Database
	}

	cfg.AccessToken = c.Token.AccessToken

	dsn := cfg.FormatDSN()
	return dsn, nil
}

type CloudToken struct {
	AccessToken  string    `toml:"access_token"`
	RefreshToken string    `toml:"refresh_token"`
	ExpiresAt    time.Time `toml:"expires_at"`
}

type RuntimeOptions struct {
	Username string
	Password string
	Database string
}

func LoadConfig() (*Config, error) {
	content, err := os.ReadFile(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "read config file")
	}

	var cfg Config
	_, err = toml.Decode(string(content), &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal config file")
	}

	// if target is cloud, try load the auth token
	var cloudToken CloudToken
	if exists(cloudTokenFile) {
		_, err = toml.DecodeFile(cloudTokenFile, &cloudToken)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshal auth token file")
		}

		cfg.Cloud.Token = &cloudToken
	}
	return &cfg, nil
}

func FlushCloudToken(token *CloudToken) error {
	file, err := os.OpenFile(cloudTokenFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrapf(err, "open cloud token file: %s", cloudTokenFile)
	}
	err = toml.NewEncoder(file).Encode(token)
	if err != nil {
		return errors.Wrapf(err, "encode cloud token file: %s", cloudTokenFile)
	}
	return nil
}

func FlushConfig(config *Config) error {
	cfg := config.Clone()

	// save the cloud token file seperately
	if cfg.Cloud != nil && cfg.Cloud.Token != nil {
		FlushCloudToken(cfg.Cloud.Token)
		// do not save the cloud token in the config file
		cfg.Cloud.Token = nil
	}

	file, err := os.OpenFile(configFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrap(err, "open config file")
	}
	err = toml.NewEncoder(file).Encode(cfg)
	if err != nil {
		return errors.Wrap(err, "encode config file")
	}
	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}
