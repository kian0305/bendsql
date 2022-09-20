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
	"sync"

	"github.com/sirupsen/logrus"

	"gopkg.in/ini.v1"
)

var (
	once sync.Once
)

const (
	UserEmail    string = "user_email"
	AccessToken  string = "access_token"
	RefreshToken string = "refresh_token"
	Warehouse    string = "warehouse"
	Org          string = "org"
)

const (
	bendsqlConfigDir  = "BENDSQL_CONFIG_DIR"
	bendsqlCinfigFile = "bendsql.ini"
)

type Config struct {
	UserEmail    string `ini:"user_email"`
	AccessToken  string `ini:"access_token"`
	RefreshToken string `ini:"refresh_token"`
	Warehouse    string `ini:"warehouse"`
	Org          string `ini:"org"`
}

type Configer interface {
	AuthToken() (string, string)
	Get(string) (string, error)
	Set(string, string) error
}

func NewConfig() (Configer, error) {
	c, err := Read()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) Write() error {
	if Exists(filepath.Join(ConfigDir(), bendsqlCinfigFile)) {
		err := os.RemoveAll(ConfigDir())
		if err != nil {
			return err
		}
	}
	if !Exists(ConfigDir()) {
		err := os.MkdirAll(ConfigDir(), os.ModePerm)
		if err != nil {
			return err
		}
	}
	if !Exists(filepath.Join(ConfigDir(), bendsqlCinfigFile)) {
		_, err := os.Create(filepath.Join(ConfigDir(), bendsqlCinfigFile))
		if err != nil {
			return err
		}
	}
	cg := ini.Empty()
	defaultSection := cg.Section("")
	defaultSection.NewKey(AccessToken, c.AccessToken)
	defaultSection.NewKey(RefreshToken, c.RefreshToken)
	defaultSection.NewKey(Warehouse, c.Warehouse)
	defaultSection.NewKey(Org, c.Org)
	defaultSection.NewKey(UserEmail, c.UserEmail)

	return cg.SaveTo(filepath.Join(ConfigDir(), bendsqlCinfigFile))
}

func (c *Config) AuthToken() (string, string) {
	accessToken, err := c.Get(AccessToken)
	refreshToken, err := c.Get(RefreshToken)
	if err != nil {
		panic(err)
	}

	return accessToken, refreshToken
}

// Get a string value from a ConfigFile.
func (c *Config) Get(key string) (string, error) {
	if !Exists(filepath.Join(ConfigDir(), bendsqlCinfigFile)) {
		return "", nil
	}
	log := logrus.WithField("bendsql", "get")
	cfg, err := ini.Load(filepath.Join(ConfigDir(), bendsqlCinfigFile))
	if err != nil {
		log.Errorf("Fail to read file: %v", err)
		return "", err
	}
	return cfg.Section("").Key(key).String(), nil
}

func (c *Config) Set(key, value string) error {
	log := logrus.WithField("bendsql", "set")
	cfg, err := ini.Load(filepath.Join(ConfigDir(), bendsqlCinfigFile))
	if err != nil {
		log.Errorf("Fail to read file: %v", err)
		return err
	}
	cfg.Section("").Key(key).SetValue(value)
	err = cfg.SaveTo(filepath.Join(ConfigDir(), bendsqlCinfigFile))
	if err != nil {
		log.Errorf("Fail to save file: %v", err)
		return err
	}
	return nil
}

func RenewTokens(accessToken, refreshToken string) error {
	if !Exists(filepath.Join(ConfigDir(), bendsqlCinfigFile)) {
		return os.ErrNotExist
	}
	cfg, err := NewConfig()
	if err != nil {
		return fmt.Errorf("config failed: %w", err)
	}
	err = cfg.Set(AccessToken, accessToken)
	err = cfg.Set(RefreshToken, refreshToken)
	if err != nil {
		return fmt.Errorf("renew tokens failed %w", err)
	}
	return nil
}

func GetAuthToken() (string, string) {
	if !Exists(filepath.Join(ConfigDir(), bendsqlCinfigFile)) {
		return "", ""
	}
	cfg, err := NewConfig()
	if err != nil {
		logrus.Errorf("read config failed %v", err)
		return "", ""
	}
	return cfg.AuthToken()
}

func GetWarehouse() string {
	if !Exists(filepath.Join(ConfigDir(), bendsqlCinfigFile)) {
		return ""
	}
	cfg, err := NewConfig()
	if err != nil {
		logrus.Errorf("read config failed %v", err)
		return ""
	}
	warehouse, err := cfg.Get(Warehouse)
	if err != nil {
		logrus.Errorf("get warehouse failed %v", err)
		return ""
	}
	return warehouse
}

func GetUserEmail() string {
	if !Exists(filepath.Join(ConfigDir(), bendsqlCinfigFile)) {
		return ""
	}
	cfg, err := NewConfig()
	if err != nil {
		logrus.Errorf("read config failed %v", err)
		return ""
	}
	userEmail, err := cfg.Get(UserEmail)
	if err != nil {
		logrus.Errorf("get userEmail failed %v", err)
		return ""
	}
	return userEmail
}

func GetOrg() string {
	if !Exists(filepath.Join(ConfigDir(), bendsqlCinfigFile)) {
		return ""
	}
	cfg, err := NewConfig()
	if err != nil {
		logrus.Errorf("read config failed %v", err)
		return ""
	}
	warehouse, err := cfg.Get(Org)
	if err != nil {
		logrus.Errorf("get org failed %v", err)
		return ""
	}
	return warehouse
}

func ConfigDir() string {
	var path string
	if a := os.Getenv(bendsqlConfigDir); a != "" {
		path = a
	} else {
		d, _ := os.UserHomeDir()
		path = filepath.Join(d, ".config", "bendsql")
	}
	return path
}

// Read bendsql configuration files from the local file system and
// return a Config.
var Read = func() (*Config, error) {
	var err error
	var iniCfg *ini.File
	cfg := &Config{}
	once.Do(func() {
		iniCfg, err = ini.Load(filepath.Join(ConfigDir(), bendsqlCinfigFile))
		err = iniCfg.MapTo(cfg)
	})
	return cfg, err
}

func Exists(path string) bool {
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
