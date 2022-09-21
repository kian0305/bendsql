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

	"github.com/databendcloud/bendsql/internal/config"

	"github.com/MakeNowJust/heredoc"
	"github.com/databendcloud/bendsql/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdWarehouseStatus(f *cmdutil.Factory) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "status warehouseName",
		Short: "show warehouse status",
		Long:  "show warehouse status",
		Example: heredoc.Doc(`
			# show warehouse status
			$ bendsql warehouse status WAREHOUSENAME
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				fmt.Printf("Wrong params, example: bendsql warehouse status WAREHOUSENAME \n")
				return
			}
			if len(args) == 0 {
				args = append(args, config.GetWarehouse())
			}
			warehouseStatus, err := showWarehouseStatus(f, args[0])
			if err != nil {
				fmt.Printf("show warehouse %s status failed, err: %v", args[0], err)
			}
			fmt.Println(warehouseStatus)
		},
	}

	return cmd
}

func showWarehouseStatus(f *cmdutil.Factory, warehouseName string) (string, error) {
	apiClient, err := f.ApiClient()
	if err != nil {
		return "", err
	}
	warehouseStatus, err := apiClient.ViewWarehouse(warehouseName)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("warehouse %s status is %s, size is %s, readyInstance is %d, totalInstance is %d",
		warehouseName, warehouseStatus.State, warehouseStatus.Size, warehouseStatus.ReadyInstances, warehouseStatus.TotalInstances), nil
}
