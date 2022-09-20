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

	"github.com/AlecAivazis/survey/v2"
	"github.com/datafuselabs/bendcloud-cli/internal/config"
	"github.com/datafuselabs/bendcloud-cli/pkg/prompt"

	"github.com/datafuselabs/bendcloud-cli/api"
	"github.com/datafuselabs/bendcloud-cli/pkg/iostreams"

	"github.com/MakeNowJust/heredoc"
	"github.com/datafuselabs/bendcloud-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

type ConfigureOptions struct {
	IO        *iostreams.IOStreams
	ApiClient func() (*api.APIClient, error)
	Config    config.Config
	Org       string
	Warehouse string
}

func NewCmdConfigure(f *cmdutil.Factory) *cobra.Command {
	opts := &ConfigureOptions{
		IO:        f.IOStreams,
		ApiClient: f.ApiClient,
	}
	var org, warehouse string

	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Set your default org and using warehouse",
		Long:  "Set your default org and using warehouse",
		Example: heredoc.Doc(`
			# Set your default org and using warehouse with flag
			# NOTE: Using flag is faster than interactive shell
			$ bendsql auth configure --org ORG --warehouse WAREHOUSENAME
			
			# Set with interactive shell
			$ bendsql auth configure
		`),
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.NewConfig()
			if err != nil {
				panic(err)
			}
			if org != "" {
				opts.Org = org
				err = cfg.Set(config.Org, org)
				if err != nil {
					panic(err)
				}
			}
			if warehouse != "" {
				opts.Warehouse = warehouse
				// TODO: check the warehouse whether in warehouse list
				err = cfg.Set(config.Warehouse, warehouse)
				if err != nil {
					panic(err)
				}
			}
			if org == "" || warehouse == "" {
				err = configureRunInteractive(opts)
				if err != nil {
					fmt.Printf("configure failed, err: %v", err)
					return
				}
			}
			fmt.Printf("configure success, current org is %s, current warehosue is %s", opts.Org, opts.Warehouse)
		},
	}

	cmd.Flags().StringVar(&org, "org", "", "org")
	cmd.Flags().StringVar(&warehouse, "warehouse", "", "warehouse")

	return cmd
}

func configureRunInteractive(opts *ConfigureOptions) error {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	apiClient, err := opts.ApiClient()
	if err != nil {
		return err
	}
	orgs, err := apiClient.ListOrgs()
	if err != nil {
		return err
	}
	warehouseDtos, err := apiClient.ListWarehouses()
	if err != nil {
		return err
	}
	var warehouses []string
	for i := range warehouseDtos {
		warehouses = append(warehouses, warehouseDtos[i].Name)
	}

	err = prompt.SurveyAskOne(
		&survey.Select{
			Message: "Select your working org:",
			Options: orgs,
			Default: orgs[0],
		}, &opts.Org, survey.WithValidator(survey.Required))
	if err != nil {
		return fmt.Errorf("could not prompt: %w", err)
	}
	err = cfg.Set(config.Org, opts.Org)
	if err != nil {
		panic(err)
	}

	err = prompt.SurveyAskOne(
		&survey.Select{
			Message: "Select your working warehouse:",
			Options: warehouses,
			Default: warehouses[0],
		}, &opts.Warehouse, survey.WithValidator(survey.Required))
	if err != nil {
		return fmt.Errorf("could not prompt: %w", err)
	}

	err = cfg.Set(config.Warehouse, opts.Warehouse)
	if err != nil {
		panic(err)
	}

	return nil
}
