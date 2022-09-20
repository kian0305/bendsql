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

package auth

import (
	"fmt"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/datafuselabs/bendcloud-cli/api"
	"github.com/datafuselabs/bendcloud-cli/internal/config"
	"github.com/datafuselabs/bendcloud-cli/pkg/cmdutil"
	"github.com/datafuselabs/bendcloud-cli/pkg/iostreams"
	"github.com/datafuselabs/bendcloud-cli/pkg/prompt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type LoginType int

const (
	UserPasswordLogin LoginType = 0
	AccessTokenLogin  LoginType = 1
)

type LoginOptions struct {
	IO             *iostreams.IOStreams
	ApiClient      func() (*api.APIClient, error)
	Interactive    bool
	MainExecutable string
	Config         config.Config
	Email          string
	Password       string
	Org            string
}

func NewCmdLogin(f *cmdutil.Factory, runF func(*LoginOptions) error) *cobra.Command {
	opts := &LoginOptions{
		IO:        f.IOStreams,
		ApiClient: f.ApiClient,
	}

	var email, password, org string
	cmd := &cobra.Command{
		Use:   "login",
		Args:  cobra.ExactArgs(0),
		Short: "Authenticate with a Databend Cloud host",
		Long: heredoc.Docf(`
			Authenticate with a Databend Cloud host.

			The default authentication mode is a user-password flow. After completion, an
			authentication token will be stored internally.

			Alternatively, bendcli will use the authentication token found in environment variables.
		`, "`"),
		Example: heredoc.Doc(`
			# start interactive setup
			$ bendsql auth login

			# authenticate by reading the token from a file
			$ bendsql auth login --email EMAIL --password PASSWORD [--org ORG]
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.IO.CanPrompt() && (email == "" || password == "") {
				// default use interactive tty
				opts.Interactive = true
			}

			opts.Email = email
			opts.Password = password
			opts.Org = org
			opts.MainExecutable = f.Executable()
			if runF != nil {
				return runF(opts)
			}

			return loginRun(opts)
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "org")
	cmd.Flags().StringVar(&email, "email", "", "email")
	cmd.Flags().StringVar(&password, "password", "", "password")
	return cmd
}

func loginRun(opts *LoginOptions) error {
	cfg := opts.Config
	apiClient, err := opts.ApiClient()
	if err != nil {
		return err
	}

	// interactive login
	if opts.Email == "" {
		err = prompt.SurveyAskOne(
			&survey.Input{
				Message: "Paste your user email:",
			}, &opts.Email, survey.WithValidator(survey.Required))
		if err != nil {
			return fmt.Errorf("could not prompt: %w", err)
		}
	}
	if opts.Password == "" {
		err = prompt.SurveyAskOne(&survey.Password{
			Message: "Paste your password:",
		}, &opts.Password, survey.WithValidator(survey.Required))
		if err != nil {
			return fmt.Errorf("could not prompt: %w", err)
		}
	}

	apiClient.UserEmail = opts.Email
	apiClient.Password = opts.Password
	err = apiClient.Login()
	if err != nil {
		return err
	}

	// get current account info
	currentAccountInfo, err := apiClient.GetCurrentAccountInfo()
	if err != nil {
		return fmt.Errorf("get current account failed: %w", err)
	}
	// TODO: new apiClient in a func: sjhan
	if cfg.UserEmail == "" {
		cfg.UserEmail = currentAccountInfo.Email
	}

	if opts.Org == "" {
		err = prompt.SurveyAskOne(&survey.Input{
			Message: "Paste your org slug:",
			Default: currentAccountInfo.DefaultOrgSlug,
		}, &opts.Org, survey.WithValidator(survey.Required))
		if err != nil {
			return fmt.Errorf("could not prompt: %w", err)
		}
	}
	apiClient.CurrentOrgSlug = opts.Org

	warehouses, err := apiClient.ListWarehouses()
	if err != nil || len(warehouses) == 0 {
		logrus.Warnf("you have no warehouse in %s", cfg.Org)
	} else {
		cfg.Warehouse = warehouses[0].Name
	}

	cfg.Org = apiClient.CurrentOrgSlug
	cfg.AccessToken = apiClient.AccessToken
	cfg.RefreshToken = apiClient.RefreshToken
	err = cfg.Write()
	if err != nil {
		return fmt.Errorf("save config failed:%w", err)
	}

	logrus.Infof("%s logged in %s of Databend Cloud successfully.", cfg.UserEmail, cfg.Org)
	return nil
}

func mustString(accountID uint64) string {
	return strconv.FormatInt(int64(accountID), 10)
}
