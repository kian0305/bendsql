package query

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/datafuselabs/bendcloud-cli/api"
	"github.com/datafuselabs/bendcloud-cli/internal/config"
	"github.com/datafuselabs/bendcloud-cli/pkg/cmdutil"
	"github.com/datafuselabs/bendcloud-cli/pkg/iostreams"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type querySQLOptions struct {
	IO        *iostreams.IOStreams
	ApiClient func() (*api.APIClient, error)
	Warehouse string
	QuerySQL  string
	Verbose   bool
}

func NewCmdQuerySQL(f *cmdutil.Factory) *cobra.Command {
	opts := &querySQLOptions{
		IO:        f.IOStreams,
		ApiClient: f.ApiClient,
	}
	var warehouse, querySQL string
	var sqlStdin bool
	var verbose bool

	cmd := &cobra.Command{
		Use:   "query",
		Short: "Exec query SQL using warehouse",
		Long:  "Exec query SQL using warehouse",
		Example: heredoc.Doc(`
			# exec SQL using warehouse 
			# use sql
			$ bendctl query --sql "YOURSQL" --warehouse [WAREHOUSENAME] [--verbose]
			
			# use stdin
			$ echo "select * from YOURTABLE limit 10" | bendctl query 
		`),
		Run: func(cmd *cobra.Command, args []string) {
			opts.Warehouse = warehouse
			opts.QuerySQL = querySQL
			opts.Verbose = verbose

			cfg, err := config.NewConfig()
			if err != nil {
				panic(err)
			}
			if len(querySQL) == 0 {
				sqlStdin = true
			}
			if sqlStdin {
				defer opts.IO.In.Close()
				sql, err := io.ReadAll(opts.IO.In)
				if err != nil {
					fmt.Printf("failed to read sql from standard input: %v", err)
					return
				}
				opts.QuerySQL = strings.TrimSpace(string(sql))
			}

			if warehouse == "" {
				// TODO: check the warehouse whether in warehouse list
				warehouse, err = cfg.Get(config.Warehouse)
				if warehouse == "" || err != nil {
					fmt.Printf("get default warehouse failed, please your default warehouse in $HOME/.config/bendctl/bendctl.ini")
					return
				}
				opts.Warehouse = warehouse
			}
			err = execQuery(opts)
			if err != nil {
				fmt.Printf("exec query failed, err: %v", err)
				return
			}
		},
	}

	cmd.Flags().StringVar(&warehouse, "warehouse", "", "warehouse")
	cmd.Flags().StringVar(&querySQL, "sql", "", "querysql")
	cmd.Flags().BoolVar(&verbose, "verbose", false, "display progress info across paginated results")

	return cmd
}

func execQuery(opts *querySQLOptions) error {
	apiClient, err := opts.ApiClient()
	if err != nil {
		return err
	}
	respCh := make(chan api.QueryResponse)
	errCh := make(chan error)
	logrus.Infof("start query in warehouse %s: %s", opts.Warehouse, opts.QuerySQL)
	go func() {
		err := apiClient.QuerySync(opts.Warehouse, opts.QuerySQL, respCh)
		errCh <- err
	}()

	for {
		select {
		case err := <-errCh:
			if err != nil {
				logrus.Errorf("error on query: %s", err)
				os.Exit(1)
			} else {
				os.Exit(0)
			}
		case resp := <-respCh:
			if opts.Verbose {
				logrus.WithFields(logrus.Fields{
					"queryID":        resp.Id,
					"runningSeconds": (resp.Stats.RunningTimeMS / 1000),
					"scanBytes":      resp.Stats.ScanProgress.Bytes,
					"scanRows":       resp.Stats.ScanProgress.Rows,
				}).Info("query progress")
			}
			var schemaStr string
			for i := range resp.Schema.Fields {
				schemaStr += fmt.Sprintf("| %v", resp.Schema.Fields[i].Name)
			}

			for i := range resp.Data {
				var a string
				for j := range resp.Data[i] {
					a += fmt.Sprintf("| %v ", resp.Data[i][j])
				}
				a += "| \n"
				fmt.Println(a)
			}
		}
	}
}
