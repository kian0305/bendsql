package warehouse

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/datafuselabs/bendcloud-cli/internal/config"
	"github.com/datafuselabs/bendcloud-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdWarehouseSuspend(f *cmdutil.Factory) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "suspend [warehouseName]",
		Short: "Suspend a warehouse",
		Long:  "Suspend a warehouse",
		Example: heredoc.Doc(`
			# suspend a warehouse 
			$ bendctl warehouse suspend [WAREHOUSENAME]
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				fmt.Printf("Wrong params, example: bendctl warehouse suspend [WAREHOUSENAME] \n")
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
	fmt.Printf("suspend warehouse %s success you can use `bendctl warehouse status WAREHOUSENAME to check`", warehouseName)

	return err
}
