package warehouse

import (
	"fmt"

	"github.com/datafuselabs/bendcloud-cli/internal/config"

	"github.com/MakeNowJust/heredoc"
	"github.com/datafuselabs/bendcloud-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdWarehouseStatus(f *cmdutil.Factory) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "status warehouseName",
		Short: "show warehouse status",
		Long:  "show warehouse status",
		Example: heredoc.Doc(`
			# show warehouse status
			$ bendctl warehouse status WAREHOUSENAME
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				fmt.Printf("Wrong params, example: bendctl warehouse status WAREHOUSENAME \n")
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
