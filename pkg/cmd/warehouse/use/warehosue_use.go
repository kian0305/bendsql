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

func NewCmdWarehouseUse(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use",
		Short: "select working warehouse",
		Long:  "select working warehouse",
		Example: heredoc.Doc(`
			# "select working warehouse",
			$ bendsql warehouse use WORKINGWAREHOUSE
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				fmt.Printf("Wrong params, example: bendsql warehouse use [WAREHOUSENAME] \n")
				return
			}
			warehouse := args[0]
			if !isWarehouseExist(f, warehouse) {
				fmt.Printf("warehouse %s not exist", warehouse)
			}
			err := config.SetUsingWarehouse(warehouse)
			if err != nil {
				fmt.Printf("set working warehouse %s failed", warehouse)
			}

			fmt.Printf("Now using warehouse %s", warehouse)
		},
	}

	return cmd
}

func isWarehouseExist(f *cmdutil.Factory, warehouse string) bool {
	apiClient, err := f.ApiClient()
	if err != nil {
		return false
	}
	warehouseList, err := apiClient.ListWarehouses()
	if err != nil {
		return false
	}
	for i := range warehouseList {
		if warehouse == warehouseList[i].Name {
			return true
		}
	}
	return false
}
