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

package query

import (
	"fmt"

	"github.com/datafuselabs/bendcloud-cli/internal/config"

	"github.com/datafuselabs/bendcloud-cli/api"
	"github.com/datafuselabs/bendcloud-cli/pkg/iostreams"

	"github.com/MakeNowJust/heredoc"
	"github.com/datafuselabs/bendcloud-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

type querySQLOptions struct {
	IO        *iostreams.IOStreams
	ApiClient func() (*api.APIClient, error)
	Warehouse string
	QuerySQL  string
}

func NewCmdQuerySQL(f *cmdutil.Factory) *cobra.Command {
	opts := &querySQLOptions{
		IO:        f.IOStreams,
		ApiClient: f.ApiClient,
	}
	var warehouse, querySQL string

	cmd := &cobra.Command{
		Use:   "sql",
		Short: "Exec query SQL using warehouse",
		Long:  "Exec query SQL using warehouse",
		Example: heredoc.Doc(`
			# exec SQL using warehouse 
			$ bendctl query --sql "YOURSQL" --warehouse [WAREHOUSENAME]
		`),
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.NewConfig()
			if err != nil {
				panic(err)
			}

			if warehouse == "" {
				// TODO: check the warehouse whether in warehouse list
				warehouse, err = cfg.Get(config.Warehouse)
				if warehouse == "" || err != nil {
					fmt.Printf("get default warehouse failed, please your default warehouse in $HOME/.config/bendctl/bendctl.ini")
					return
				}
			}
			opts.Warehouse = warehouse
			opts.QuerySQL = args[0]
			err = execQuery(opts)
			if err != nil {
				fmt.Printf("exec query failed, err: %v", err)
				return
			}
		},
	}

	cmd.Flags().StringVar(&warehouse, "warehouse", "", "warehouse")
	cmd.Flags().StringVar(&querySQL, "sql", "", "querysql")

	return cmd
}

func execQuery(opts *querySQLOptions) error {
	apiClient, err := opts.ApiClient()
	if err != nil {
		return err
	}
	queryResp, err := apiClient.Query(opts.Warehouse, opts.QuerySQL)
	if err != nil {
		return err
	}
	fmt.Println(queryResp)

	return nil
}
