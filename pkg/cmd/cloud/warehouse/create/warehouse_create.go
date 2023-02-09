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
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/databendcloud/bendsql/api"
	"github.com/databendcloud/bendsql/pkg/cmdutil"
)

func NewCmdWarehouseCreate(f *cmdutil.Factory) *cobra.Command {
	var size, tag string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a warehouse",
		Long:  "Create a warehouse",
		Example: heredoc.Doc(`
			# create a warehouse, the size has XSmall, Small, Medium, Large, default is Small
			$ bendsql cloud warehouse create [WAREHOUSE] --size Small
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return errors.New("wrong params")
			}
			if len(args) == 0 {
				return errors.New("warehouse name is required")
			}
			err := createWarehouse(f, args[0], size, tag)
			if err != nil {
				return errors.Errorf("create warehouse %s failed, err: %v", args[0], err)
			}
			fmt.Printf("warehouse %s created, size is %s\n", args[0], size)
			return nil
		},
	}
	cmd.Flags().StringVarP(&size, "size", "", "Small", "Warehouse size")
	cmd.Flags().StringVarP(&tag, "tag", "", "", "Databend query image tag, default to cloud stable.\nNOT RECOMMENDED to use this option, use it at your own risk")

	return cmd
}

func createWarehouse(f *cmdutil.Factory, warehouseName, size, tag string) error {
	fmt.Printf("warehouse %s is creating, please wait...\n", warehouseName)
	apiClient, err := api.NewClient()
	if err != nil {
		return err
	}
	err = apiClient.CreateWarehouse(warehouseName, size, tag)
	if err != nil {
		return err
	}
	return nil
}
