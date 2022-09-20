package warehouse

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/datafuselabs/bendcloud-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdWarehouseDelete(f *cmdutil.Factory) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete warehouseName",
		Short: "Delete a warehouse",
		Long:  "Delete a warehouse",
		Example: heredoc.Doc(`
			# delete a warehouse
			$ bendctl warehouse delete WAREHOUSENAME
		`),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 1 {
				fmt.Printf("Wrong params, example: bendctl warehouse delete WAREHOUSENAME \n")
				return
			}
			if len(args) == 0 {
				fmt.Printf("No warehouseName, example: bendctl warehouse delete WAREHOUSENAME \n")
				return
			}
			err := deleteWarehouse(f, args[0])
			if err != nil {
				fmt.Printf("delete warehouse %s failed, err: %v", args[0], err)
				return
			}
			fmt.Printf("warehouse %s deleted", args[0])
		},
	}
	return cmd
}

func deleteWarehouse(f *cmdutil.Factory, warehouseName string) error {
	apiClient, err := f.ApiClient()
	if err != nil {
		return err
	}
	err = apiClient.DeleteWarehouse(warehouseName)
	if err != nil {
		return err
	}
	return nil
}
