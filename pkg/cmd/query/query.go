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

package query

import (
	"context"
	"os"
	"os/user"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/xo/usql/drivers"
	"github.com/xo/usql/env"
	"github.com/xo/usql/handler"
	"github.com/xo/usql/rline"

	"github.com/databendcloud/bendsql/internal/config"
	"github.com/databendcloud/bendsql/pkg/cmdutil"
)

type querySQLOptions struct {
	NonInteractive bool
	Format         string
	RowsOnly       bool
	Expanded       bool
	LineStyle      string
}

var (
	outputFormats = []string{"table", "unaligned", "html", "json", "csv", "vertical"}
	LineStyles    = []string{"ascii", "unicode-single", "unicode-double"}
)

func NewCmdQuery(f *cmdutil.Factory) *cobra.Command {
	opts := &querySQLOptions{}
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Run query SQL using warehouse",
		Long:  "Run query SQL using warehouse or use interactive mode",
		Annotations: map[string]string{
			"IsCore": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return err
			}
			dsn, err := cfg.GetDSN()
			if err != nil {
				return errors.Wrap(err, "failed to get dsn")
			}

			// register databend driver
			drivers.Register("databend", drivers.Driver{
				UseColumnTypes: true,
			})

			// load current user
			cur, err := user.Current()
			if err != nil {
				return errors.Wrap(err, "failed to get current user")
			}
			wd, err := os.Getwd()
			if err != nil {
				return errors.Wrap(err, "failed to get current working directory")
			}

			// create input/output
			l, err := rline.New(opts.NonInteractive, "", env.HistoryFile(cur))
			if err != nil {
				return errors.Wrap(err, "failed to create readline")
			}
			defer l.Close()

			if !l.Interactive() {
				env.Set("QUIET", "on")
				env.Pset("format", opts.Format)
			}
			if opts.RowsOnly {
				env.Pset("tuples_only", "on")
				env.Pset("border", "0")
			} else {
				env.Pset("border", "2")
			}
			if opts.Expanded {
				env.Pset("expanded", "on")
			}
			switch opts.LineStyle {
			case "ascii":
				env.Pset("linestyle", "ascii")
			case "unicode", "unicode-single":
				env.Pset("linestyle", "unicode")
				env.Pset("border", "2")
			case "unicode-double":
				env.Pset("linestyle", "unicode")
				env.Pset("unicode_border_linestyle", "double")
				env.Pset("border", "2")
			default:
				return errors.Errorf("invalid line style: %q", opts.LineStyle)
			}

			// create handler
			h := handler.New(l, cur, wd, true)
			// open dsn
			if err = h.Open(context.Background(), dsn); err != nil {
				return errors.Wrap(err, "failed to open dsn")
			}
			return h.Run()
		},
	}
	cmd.Flags().BoolVarP(&opts.NonInteractive, "non-interactive", "n", false, "Do not use interactive mode")
	cmd.Flags().StringVarP(&opts.Format, "format", "f", "table",
		"Output format, one of: "+strings.Join(outputFormats, ", "))
	cmd.Flags().BoolVarP(&opts.RowsOnly, "rows-only", "t", false, "Output print rows only")
	cmd.Flags().BoolVarP(&opts.Expanded, "expanded", "x", false, "Table output rurn on expanded mode")
	cmd.Flags().StringVarP(&opts.LineStyle, "line-style", "l", "ascii",
		"Table output line style, one of: "+strings.Join(LineStyles, ", "))

	return cmd
}
