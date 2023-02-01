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

package benchmark

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/MakeNowJust/heredoc"
	dc "github.com/databendcloud/databend-go"
	"github.com/montanaflynn/stats"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/databendcloud/bendsql/internal/config"
	"github.com/databendcloud/bendsql/pkg/cmdutil"
)

type benchmarkOptions struct {
	WarmCount    int
	TestCount    int
	TestDir      string
	OutputFormat string
	OutputDir    string
	Tags         []string
}

func NewCmdBenchmark(f *cmdutil.Factory) *cobra.Command {
	opts := &benchmarkOptions{}

	cmd := &cobra.Command{
		Use:   "benchmark",
		Short: "Run benchmark test",
		Long:  "Run benchmark test",
		Example: heredoc.Doc(`
			# run benchmark test

			$ bendsql benchmark
		`),
		Annotations: map[string]string{
			"IsCore": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.GetConfig()
			if err != nil {
				return errors.Wrap(err, "failed to get config")
			}
			dsn, err := cfg.GetDSN()
			if err != nil {
				return errors.Wrap(err, "failed to get dsn")
			}
			dcConfig, err := dc.ParseDSN(dsn)
			if err != nil {
				return errors.Wrap(err, "failed to parse dsn")
			}
			cli := dc.NewAPIClientFromConfig(dcConfig)

			fmt.Printf("Running benchmark with options: %+v\n", opts)
			targets, err := ReadTargetFiles(opts.TestDir)
			if err != nil {
				return errors.Wrap(err, "ReadTargetFiles")
			}
			for _, target := range targets {
				err := runTarget(target, cli, opts)
				if err != nil {
					return errors.Wrapf(err, "runTarget(%+v)", target)
				}
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&opts.WarmCount, "warm", "w", 3, "warm up count for each benchmark")
	cmd.Flags().IntVarP(&opts.TestCount, "count", "c", 10, "test count for each benchmark")
	cmd.Flags().StringVarP(&opts.TestDir, "test-dir", "d", "./testdata", "test directory")
	cmd.Flags().StringVarP(&opts.OutputFormat, "output-format", "f", "json", "output format such as json, yaml")
	cmd.Flags().StringVarP(&opts.OutputDir, "output-dir", "o", "./target", "output directory to store tests")
	cmd.Flags().StringSliceVarP(&opts.Tags, "tags", "t", []string{}, "tags for the test")

	return cmd
}

func runQuery(ctx context.Context, cli *dc.APIClient, query string) (*dc.QueryStats, error) {
	r0, err := cli.DoQuery(ctx, query, []driver.Value{})
	if err != nil {
		return nil, errors.Wrap(err, "DoQuery")
	}
	if r0.Error != nil {
		return nil, errors.Wrapf(err, "DoQuery: %s", r0.Error)
	}
	s := r0.Stats
	nextURI := r0.NextURI
	for len(nextURI) != 0 {
		p, err := cli.QueryPage(nextURI)
		if err != nil {
			return nil, errors.Wrap(err, "QueryPage")
		}
		if p.Error != nil {
			return nil, fmt.Errorf("query has error: %s", p.Error)
		}
		nextURI = p.NextURI
		if p.Stats.RunningTimeMS > 0 {
			s = p.Stats
		}
	}
	return &s, err
}

func runTarget(target *InputQueryFile, cli *dc.APIClient, opts *benchmarkOptions) error {
	ctx := context.Background()

	output := &OutputFile{}
	// output.MetaData.Tag = cfg.WarehouseTag
	output.MetaData.Table = target.MetaData.Table
	output.Schema = make([]OutputSchema, 0)

	for _, i := range target.Statements {
		fmt.Printf("\nstart to run query %s : %s\n", i.Name, i.Query)

		o := OutputSchema{}
		o.Error = make([]string, 0)
		o.Time = make([]float64, 0)
		o.Name = i.Name
		o.SQL = i.Query

		realQuery, err := RenderQueryStatment(i.Query)
		if err != nil {
			return errors.Wrapf(err, "RenderQueryStatment(%s)", i.Query)
		}
		for j := 0; j < opts.WarmCount; j++ {
			_, _ = runQuery(ctx, cli, realQuery)
		}
		fmt.Printf("%s finished warm up %d times: %s\n", i.Name, opts.WarmCount, i.Query)

		testOK := false
		for j := 0; j < opts.TestCount; j++ {
			fmt.Printf("%s[%d] running...\n", i.Name, j)

			if s, err := runQuery(ctx, cli, realQuery); err != nil {
				fmt.Printf("%s[%d] result has error %s, stats: %+v\n", i.Name, j, err.Error(), s)
				o.Error = append(o.Error, err.Error())
			} else {
				testOK = true
				fmt.Printf("%s[%d] result in raw: %+v\n", i.Name, j, s)
				o.ReadRow = s.ScanProgress.Rows
				o.ReadByte = s.ScanProgress.Bytes
				ms := s.RunningTimeMS
				t := float64(time.Duration(ms)*time.Millisecond) / float64(time.Second)
				o.Time = append(o.Time, t)
			}
			if len(o.Time) > 0 {
				o.Min, _ = stats.Min(o.Time)
				o.Max, _ = stats.Max(o.Time)
				o.Median, _ = stats.Median(o.Time)
				o.Mean, _ = stats.GeometricMean(o.Time)
				o.StdDev, _ = stats.StandardDeviation(o.Time)
			}
		}
		if !testOK {
			return errors.Wrapf(err, "%s failed %d times", i.Name, opts.TestCount)
		}
		output.Schema = append(output.Schema, o)
	}
	return generateOutput(opts, output)
}

func generateOutput(opts *benchmarkOptions, output *OutputFile) error {
	switch opts.OutputFormat {
	case "json":
		b, err := json.Marshal(output)
		if err != nil {
			return errors.Wrap(err, "failed to marshal json")
		}
		outFile := filepath.Join(opts.OutputDir, output.MetaData.Table+".json")
		err = os.WriteFile(outFile, b, 0644)
		if err != nil {
			return errors.Wrapf(err, "cannot write %s", outFile)
		}
	case "yaml":
		b, err := yaml.Marshal(output)
		if err != nil {
			fmt.Printf("failed to marshal yaml : %+v\n", err)
		}
		outFile := filepath.Join(opts.OutputDir, output.MetaData.Table+".yaml")
		err = os.WriteFile(outFile, b, 0644)
		if err != nil {
			return errors.Wrapf(err, "cannot write %s", outFile)
		}
	default:
		return errors.Errorf("unsupported output type %s", opts.OutputFormat)
	}
	return nil
}
