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

package token

import (
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"github.com/databendcloud/bendsql/api"
	"github.com/databendcloud/bendsql/internal/config"
	"github.com/databendcloud/bendsql/pkg/cmdutil"
	"github.com/databendcloud/bendsql/pkg/iostreams"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type TokenOptions struct {
	IO             *iostreams.IOStreams
	ApiClient      func() (*api.APIClient, error)
	MainExecutable string
	AccessToken    string
	RefreshToken   string
	Org            string
	Config         config.Config
}

func NewCmdAuthToken(f *cmdutil.Factory, runF func(*TokenOptions) error) *cobra.Command {
	opts := &TokenOptions{
		IO:        f.IOStreams,
		ApiClient: f.ApiClient,
	}

	var accessToken, refreshToken, org string
	cmd := &cobra.Command{
		Use:   "token",
		Args:  cobra.ExactArgs(0),
		Short: "Authenticate by access & refresh token",
		Long: heredoc.Docf(`
			Authenticate by access & refresh token from Databend Cloud.
		`, "`"),
		Example: heredoc.Doc(`
			# authenticate by tokens
			$ bendsql auth token --accessToken ACCESSTOKEN --refreshToken REFRESHTOKEN [--org ORG]
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if accessToken == "" || refreshToken == "" {
				cmd.Help()
				os.Exit(0)
			}

			opts.AccessToken = accessToken
			opts.RefreshToken = refreshToken
			opts.Org = org

			opts.MainExecutable = f.Executable()
			if runF != nil {
				return runF(opts)
			}

			return runAuthByToken(opts)
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "org")
	cmd.Flags().StringVar(&accessToken, "accessToken", "", "accessToken")
	cmd.Flags().StringVar(&refreshToken, "refreshToken", "", "refreshToken")
	return cmd
}

func runAuthByToken(opts *TokenOptions) error {
	cfg := opts.Config
	apiClient, err := opts.ApiClient()
	if err != nil {
		return err
	}

	apiClient.AccessToken = opts.AccessToken
	apiClient.RefreshToken = opts.RefreshToken
	// get current account info
	currentAccountInfo, err := apiClient.GetCurrentAccountInfo()
	if err != nil {
		return fmt.Errorf("error validating token: %w", err)
	}
	warehouses, err := apiClient.ListWarehouses()
	// TODO: error type sjh
	if err != nil || len(warehouses) == 0 {
		return fmt.Errorf("you have no warehouse in your account, please create one first")
	}
	apiClient.CurrentOrgSlug = currentAccountInfo.DefaultOrgSlug
	apiClient.CurrentWarehouse = warehouses[0].Name
	if opts.Org == "" {
		// TODO: check org slug exists
		cfg.Org = currentAccountInfo.DefaultOrgSlug
	}
	cfg.AccessToken = opts.AccessToken
	cfg.RefreshToken = opts.RefreshToken
	cfg.Warehouse = warehouses[0].Name
	err = cfg.Write()
	if err != nil {
		return fmt.Errorf("save bendsql config failed: %w", err)
	}

	logrus.Infof("%s logged in %s of Databend Cloud successfully.", currentAccountInfo.Email, cfg.Org)
	return nil
}

func mustString(accountID uint64) string {
	return strconv.FormatInt(int64(accountID), 10)
}
