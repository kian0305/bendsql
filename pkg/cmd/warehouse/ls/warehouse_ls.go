package warehouse

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/datafuselabs/bendcloud-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdWarehouseList(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "show warehouse list",
		Long:  "show warehouse list",
		Example: heredoc.Doc(`
			# show warehouse list
			$ bendctl warehouse ls
		`),
		Run: func(cmd *cobra.Command, args []string) {
			warehouseList, err := showWarehouseList(f)
			if err != nil {
				fmt.Printf("show warehouse list failed, err: %v", err)
			}
			fmt.Println(warehouseList)
		},
	}

	return cmd
}

func showWarehouseList(f *cmdutil.Factory) (string, error) {
	var warehouseListStr string
	apiClient, err := f.ApiClient()
	if err != nil {
		return "", err
	}
	warehouseList, err := apiClient.ListWarehouses()
	if err != nil {
		return "", err
	}
	for i := range warehouseList {
		warehouseListStr = warehouseListStr + warehouseList[i].Name + "\n"
	}

	return warehouseListStr, nil
}
