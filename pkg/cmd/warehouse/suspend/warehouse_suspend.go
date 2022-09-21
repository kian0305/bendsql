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

package warehouse

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/databendcloud/bendsql/internal/config"
	"github.com/databendcloud/bendsql/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdWarehouseSuspend(f *cmdutil.Factory) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "suspend [warehouseName]",
		Short: "Suspend a warehouse",
		Long:  "Suspend a warehouse",
		Example: heredoc.Doc(`
			# suspend a warehouse 
			$ bendsql warehouse suspend [WAREHOUSENAME]
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				fmt.Printf("Wrong params, example: bendsql warehouse suspend [WAREHOUSENAME] \n")
				return
			}
			if len(args) == 0 {
				args = append(args, config.GetWarehouse())
			}
			err := suspendWarehouse(f, args[0])
			if err != nil {
				fmt.Printf("suspend warehouse %s failed,err: %v", args[0], err)
			}
		},
	}
	return cmd
}

func suspendWarehouse(f *cmdutil.Factory, warehouseName string) error {
	apiClient, err := f.ApiClient()
	if err != nil {
		return err
	}
	err = apiClient.SuspendWarehouse(warehouseName)
	fmt.Printf("suspend warehouse %s success you can use `bendsql warehouse status WAREHOUSENAME to check`", warehouseName)

	return err
}
