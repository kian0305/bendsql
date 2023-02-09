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

	"github.com/MakeNowJust/heredoc"
	"github.com/avast/retry-go"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/databendcloud/bendsql/api"
	"github.com/databendcloud/bendsql/pkg/cmdutil"
)

func NewCmdWarehouseResume(f *cmdutil.Factory) *cobra.Command {
	var (
		wait    bool
		timeout time.Duration
	)

	cmd := &cobra.Command{
		Use:   "resume",
		Short: "Resume a warehouse",
		Long:  "Resume a warehouse",
		Args:  cobra.MaximumNArgs(1),
		Example: heredoc.Doc(`
			# resume a warehouse and return until the warehouse running
			$ bendsql cloud warehouse resume [WAREHOUSE] --wait

			# resume a warehouse and return
			$ bendsql cloud warehouse resume [WAREHOUSE]
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			apiClient, err := api.NewClient()
			if err != nil {
				return errors.Wrap(err, "get api client failed")
			}
			var warehouse string
			switch len(args) {
			case 0:
				warehouse = apiClient.CurrentWarehouse()
			case 1:
				warehouse = args[0]
			default:
				return errors.New("wrong params")
			}
			err = resumeWarehouse(apiClient, warehouse, wait, timeout)
			if err != nil {
				return errors.Wrapf(err, "resume warehouse %s failed", warehouse)
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&wait, "wait", "w", false, "Wait until resume warehouse success")
	cmd.Flags().DurationVarP(&timeout, "timeout", "t", 2*time.Minute, "Timeout for resume warehouse")
	return cmd
}

func resumeWarehouse(apiClient *api.Client, warehouseName string, wait bool, timeout time.Duration) error {
	err := apiClient.ResumeWarehouse(warehouseName)
	if err != nil {
		return errors.Wrap(err, "resume warehouse failed")
	}

	if wait {
		err = retry.Do(
			func() (err error) {
				status, err := apiClient.ViewWarehouse(warehouseName)
				if err != nil {
					panic(err)
				}
				if status.State != "Running" {
					return fmt.Errorf("resume warehouse %s timeout, state is %s", warehouseName, status.State)
				}
				fmt.Printf("Resume warehouse %s succeed.\n", warehouseName)
				return nil
			},
			retry.Delay(1*time.Second),
			retry.Attempts(uint(timeout/time.Second)),
		)
		if err != nil {
			return errors.Wrap(err, "wait for resume warehouse failed")
		}
	}
	fmt.Printf("Resume warehouse %s done, please check with `bendsql cloud warehouse status [WAREHOUSE]`\n", warehouseName)
	return nil
}
