/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package warehouse

import (
	warehouseCreateCmd "github.com/datafuselabs/bendcloud-cli/pkg/cmd/warehouse/create"
	warehouseDeleteCmd "github.com/datafuselabs/bendcloud-cli/pkg/cmd/warehouse/delete"
	warehouseListCmd "github.com/datafuselabs/bendcloud-cli/pkg/cmd/warehouse/ls"
	warehouseResumeCmd "github.com/datafuselabs/bendcloud-cli/pkg/cmd/warehouse/resume"
	warehouseStatusCmd "github.com/datafuselabs/bendcloud-cli/pkg/cmd/warehouse/status"
	warehouseSuspendCmd "github.com/datafuselabs/bendcloud-cli/pkg/cmd/warehouse/suspend"
	"github.com/datafuselabs/bendcloud-cli/pkg/cmdutil"

	"github.com/spf13/cobra"
)

// NewWarehouseCmd represents the warehouse command
func NewWarehouseCmd(f *cmdutil.Factory) *cobra.Command {
	warehouseCmd := &cobra.Command{
		Use:   "warehouse cmd",
		Short: "Operate warehouse",
		Long: `Operate warehouse. For example:
            bendctl warehouse ls
            bendctl warehouse status YOUR_WAREHOUSE
            bendctl warehouse suspend YOUR_WAREHOUSE
`,
		Annotations: map[string]string{
			"IsCore": "true",
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			// refresh tokens every warehouse cmd
			apiClient, err := f.ApiClient()
			if err != nil {
				panic(err)
			}
			err = apiClient.RefreshTokens()
			if err != nil {
				panic(err)
			}
		},
	}
	warehouseCmd.AddCommand(warehouseStatusCmd.NewCmdWarehouseStatus(f))
	warehouseCmd.AddCommand(warehouseListCmd.NewCmdWarehouseList(f))
	warehouseCmd.AddCommand(warehouseResumeCmd.NewCmdWarehouseResume(f))
	warehouseCmd.AddCommand(warehouseSuspendCmd.NewCmdWarehouseSuspend(f))
	warehouseCmd.AddCommand(warehouseCreateCmd.NewCmdWarehouseCreate(f))
	warehouseCmd.AddCommand(warehouseDeleteCmd.NewCmdWarehouseDelete(f))
	return warehouseCmd
}
