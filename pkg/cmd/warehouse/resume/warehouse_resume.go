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
	"time"

	"github.com/avast/retry-go"

	"github.com/MakeNowJust/heredoc"
	"github.com/databendcloud/bendsql/internal/config"
	"github.com/databendcloud/bendsql/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdWarehouseResume(f *cmdutil.Factory) *cobra.Command {

	var wait bool
	cmd := &cobra.Command{
		Use:   "resume warehouseName --wait",
		Short: "Resume a warehouse",
		Long:  "Resume a warehouse",
		Example: heredoc.Doc(`
			# resume a warehouse and return until the warehouse running
			$ bendsql warehouse resume WAREHOUSENAME --wait
			
			# resume a warehouse and return 
			$ bendsql warehouse resume WAREHOUSENAME
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				fmt.Printf("Wrong params, example: bendsql warehouse resume WAREHOUSENAME \n")
				return
			}
			if len(args) == 0 {
				args = append(args, config.GetWarehouse())
			}
			err := resumeWarehouse(f, args[0], wait)
			if err != nil {
				fmt.Printf("resume warehouse %s failed,err: %v", args[0], err)
			}
		},
	}
	cmd.Flags().BoolVar(&wait, "wait", false, "Wait until resume warehouse success")
	return cmd
}

func resumeWarehouse(f *cmdutil.Factory, warehouseName string, wait bool) error {
	apiClient, err := f.ApiClient()
	if err != nil {
		return err
	}
	if wait {
		err = retry.Do(
			func() (err error) {
				err = apiClient.ResumeWarehouse(warehouseName)
				if err != nil {
					panic(err)
				}
				status, err := apiClient.ViewWarehouse(warehouseName)
				if err != nil {
					panic(err)
				}
				if status.State != "Running" {
					return fmt.Errorf("resume warehouse %s timeout, state is %s", warehouseName, status.State)
				}
				fmt.Printf("resume warehouse %s success", warehouseName)
				return nil
			},
			retry.Delay(1*time.Second),
			retry.Attempts(20),
		)
	}
	err = apiClient.ResumeWarehouse(warehouseName)
	fmt.Printf("resume warehouse %s done please use `bendsql warehouse status WAREHOUSENAME to check`", warehouseName)

	return err
}
