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
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/databendcloud/bendsql/api"
	"github.com/databendcloud/bendsql/pkg/cmdutil"
	"github.com/databendcloud/bendsql/pkg/prompt"
)

type LoginOptions struct {
	Email    string
	Password string
	Org      string
	Endpoint string
}

func NewCmdLogin(f *cmdutil.Factory) *cobra.Command {
	opts := &LoginOptions{}
	cmd := &cobra.Command{
		Use:   "login",
		Args:  cobra.ExactArgs(0),
		Short: "Authenticate with a Databend Cloud host",
		Long: heredoc.Docf(`
			Authenticate with a Databend Cloud host.

			The default authentication mode is a user-password flow. After completion, an
			authentication token will be stored internally.

			Alternatively, bendcli will use the authentication token found in environment variables.`,
		),
		Example: heredoc.Doc(`
			# start interactive setup
			$ bendsql cloud login

			# authenticate by reading the token from a file
			$ bendsql cloud login --email EMAIL --password PASSWORD [--org ORG]
		`),
		Annotations: map[string]string{
			"IsCore": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return loginRun(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Org, "org", "", "", "org")
	cmd.Flags().StringVarP(&opts.Email, "email", "", "", "email")
	cmd.Flags().StringVarP(&opts.Password, "password", "", "", "password")
	cmd.Flags().StringVarP(&opts.Endpoint, "endpoint", "", "", "endpoint")
	return cmd
}

func loginRun(opts *LoginOptions) error {
	apiClient, err := api.NewClient()
	if err != nil {
		return errors.Wrap(err, "could not create api client")
	}

	if endpoint := os.Getenv("BENDSQL_API_ENDPOINT"); endpoint != "" {
		opts.Endpoint = endpoint
	}
	// interactive select endpoint
	if opts.Endpoint == "" {
		err = prompt.SurveyAskOne(
			&survey.Select{
				Message: "Select your login endpoint:",
				Options: []string{api.EndpointGlobal, api.EndpointCN},
				Default: api.EndpointGlobal,
				Description: func(value string, index int) string {
					switch value {
					case api.EndpointGlobal:
						return "Global"
					case api.EndpointCN:
						return "China"
					default:
						return ""
					}
				},
			}, &opts.Endpoint, survey.WithValidator(survey.Required))
		if err != nil {
			return errors.Wrap(err, "could not prompt")
		}
	}

	// interactive login
	if opts.Email == "" {
		err = prompt.SurveyAskOne(
			&survey.Input{
				Message: "Paste your user email:",
			}, &opts.Email, survey.WithValidator(survey.Required))
		if err != nil {
			return errors.Wrap(err, "could not prompt")
		}
	}
	if opts.Password == "" {
		err = prompt.SurveyAskOne(&survey.Password{
			Message: "Paste your password:",
		}, &opts.Password, survey.WithValidator(survey.Required))
		if err != nil {
			return errors.Wrap(err, "could not prompt")
		}
	}

	apiClient.SetEndpoint(opts.Endpoint)
	err = apiClient.Login(opts.Email, opts.Password)
	if err != nil {
		return err
	}

	orgDtos, err := apiClient.ListOrgs()
	if err != nil {
		return errors.Wrap(err, "list orgs failed")
	}

	var currentOrg *api.OrgMembershipDTO
	switch len(orgDtos) {
	case 0:
		return fmt.Errorf("no orgs found, please create one first")
	case 1:
		currentOrg = &orgDtos[0]
	default:
		if opts.Org == "" {
			var orgs []string
			for i := range orgDtos {
				orgs = append(orgs, orgDtos[i].OrgSlug)
			}
			err = prompt.SurveyAskOne(
				&survey.Select{
					Message: "Select your working org:",
					Options: orgs,
					Default: orgs[0],
					Description: func(value string, index int) string {
						return orgDtos[index].Description()
					},
				}, &opts.Org, survey.WithValidator(survey.Required))
			if err != nil {
				return errors.Wrap(err, "could not prompt")
			}
		}
		for _, org := range orgDtos {
			if org.OrgSlug == opts.Org {
				currentOrg = &org
				break
			}
		}
	}

	if currentOrg == nil {
		return fmt.Errorf("org %s not found", opts.Org)
	}

	apiClient.SetCurrentOrg(currentOrg.OrgSlug, currentOrg.OrgTenantID, currentOrg.Gateway)
	err = apiClient.WriteConfig()
	if err != nil {
		return errors.Wrap(err, "could not write config")
	}

	logrus.Infof("logged in %s of Databend Cloud %s successfully.",
		apiClient.CurrentOrganization(), apiClient.CurrentEndpoint())
	logrus.Infoln("you can use `bendsql cloud configure` to switch to another org and warehouse")
	return nil
}
