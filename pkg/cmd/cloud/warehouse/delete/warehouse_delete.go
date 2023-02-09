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

func NewCmdWarehouseDelete(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a warehouse",
		Long:  "Delete a warehouse",
		Example: heredoc.Doc(`
			# delete a warehouse
			$ bendsql cloud warehouse delete [WAREHOUSE]
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return errors.New("wrong params")
			}
			if len(args) == 0 {
				return errors.New("warehouse name is required")
			}
			err := deleteWarehouse(f, args[0])
			if err != nil {
				return errors.Errorf("Delete warehouse %s failed, err: %v", args[0], err)
			}
			fmt.Printf("Warehouse %s deleted.\n", args[0])
			return nil
		},
	}
	return cmd
}

func deleteWarehouse(f *cmdutil.Factory, warehouseName string) error {
	apiClient, err := api.NewClient()
	if err != nil {
		return err
	}
	err = apiClient.DeleteWarehouse(warehouseName)
	if err != nil {
		return err
	}
	return nil
}
