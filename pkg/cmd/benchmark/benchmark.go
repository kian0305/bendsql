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
	"strings"
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
	Tag          string
	Size         string
}

func NewCmdBenchmark(f *cmdutil.Factory) *cobra.Command {
	opts := &benchmarkOptions{}

	cmd := &cobra.Command{
		Use:   "benchmark",
		Short: "Run benchmark",
		Long:  "Run benchmark",
		Example: heredoc.Doc(`
			# run benchmark

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
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&opts.WarmCount, "warm", "w", 3, "warm up count for each benchmark")
	cmd.Flags().IntVarP(&opts.TestCount, "count", "c", 10, "test count for each benchmark")
	cmd.Flags().StringVarP(&opts.TestDir, "test-dir", "d", "./testdata", "test directory")
	cmd.Flags().StringVarP(&opts.OutputFormat, "output-format", "f", "json", "comma separated format: json, yaml, md")
	cmd.Flags().StringVarP(&opts.OutputDir, "output-dir", "o", "./target", "output directory to store tests")
	cmd.Flags().StringVarP(&opts.Tag, "tags", "t", "", "tag for the test")
	cmd.Flags().StringVarP(&opts.Size, "size", "s", "", "size of the test")

	return cmd
}

func runQuery(ctx context.Context, cli *dc.APIClient, query string) (*dc.QueryStats, error) {
	r0, err := cli.DoQuery(ctx, query, []driver.Value{})
	if err != nil {
		return nil, errors.Wrap(err, "DoQuery")
	}
	if r0.Error != nil {
		return nil, fmt.Errorf("query has error: %s", r0.Error)
	}
	s := r0.Stats
	nextURI := r0.NextURI
	for nextURI != "" {
		p, err := cli.QueryPage(nextURI)
		if err != nil {
			return nil, errors.Wrap(err, "QueryPage")
		}
		if p.Error != nil {
			return nil, fmt.Errorf("query page has error: %s", p.Error)
		}
		nextURI = p.NextURI
		if p.Stats.RunningTimeMS > 0 {
			s = p.Stats
		}
	}
	return &s, nil
}

func runTarget(target *InputQueryFile, cli *dc.APIClient, opts *benchmarkOptions) error {
	ctx := context.Background()

	output := &OutputFile{}
	output.MetaData.Tag = opts.Tag
	output.MetaData.Size = opts.Size
	output.MetaData.Table = target.MetaData.Table
	output.Schema = make([]OutputSchema, 0)

	for _, i := range target.Statements {
		fmt.Printf("\nstart to run query %s : %s\n", i.Name, i.Query)

		o := OutputSchema{}
		o.Error = make([]string, 0)
		o.Time = make([]float64, 0)
		o.Name = i.Name
		o.SQL = i.Query

		for j := 0; j < opts.WarmCount; j++ {
			_, _ = runQuery(ctx, cli, i.Query)
		}
		fmt.Printf("%s finished warm up %d times\n", i.Name, opts.WarmCount)

		testOK := false
		for j := 0; j < opts.TestCount; j++ {
			fmt.Printf("%s[%d] running...\n", i.Name, j)

			if s, err := runQuery(ctx, cli, i.Query); err != nil {
				fmt.Printf("%s[%d] result has error: %s\n", i.Name, j, err.Error())
				o.Error = append(o.Error, err.Error())
			} else {
				testOK = true
				fmt.Printf("%s[%d] result stats: %.2f ms, %d bytes, %d rows\n",
					i.Name, j, s.RunningTimeMS, s.ScanProgress.Bytes, s.ScanProgress.Rows)
				o.ReadRow = s.ScanProgress.Rows
				o.ReadByte = s.ScanProgress.Bytes
				ms := s.RunningTimeMS
				t := float64(time.Duration(ms)*time.Millisecond) / float64(time.Second)
				o.Time = append(o.Time, t)
			}
		}
		if len(o.Time) > 0 {
			o.Min, _ = stats.Min(o.Time)
			o.Max, _ = stats.Max(o.Time)
			o.Median, _ = stats.Median(o.Time)
			o.Mean, _ = stats.GeometricMean(o.Time)
			o.StdDev, _ = stats.StandardDeviation(o.Time)
		}
		if !testOK {
			return fmt.Errorf("test failed for %s", i.Name)
		}
		output.Schema = append(output.Schema, o)
	}
	return generateOutput(opts, output)
}

func generateOutput(opts *benchmarkOptions, output *OutputFile) error {
	for _, format := range strings.Split(opts.OutputFormat, ",") {
		var data []byte
		switch format {
		case "json":
			b, err := json.Marshal(output)
			if err != nil {
				return errors.Wrap(err, "failed to marshal json")
			}
			data = b
		case "yaml":
			b, err := yaml.Marshal(output)
			if err != nil {
				fmt.Printf("failed to marshal yaml : %+v\n", err)
			}
			data = b
		case "markdown", "md":
			text := fmt.Sprintf("## Benchmark for %s with `%s`\n\n", output.MetaData.Table, output.MetaData.Size)
			if output.MetaData.Tag != "" {
				text += fmt.Sprintf("tag: `%s`\n\n", output.MetaData.Tag)
			}
			text += "|Name|Min|Max|Median|Mean|StdDev|ReadRow|ReadByte|\n"
			text += "|----|---|---|------|----|------|-------|--------|\n"
			for _, o := range output.Schema {
				text += fmt.Sprintf("|%s|%.2f|%.2f|%.2f|%.2f|%.2f|%d|%d|\n",
					o.Name, o.Min, o.Max, o.Median, o.Mean, o.StdDev, o.ReadRow, o.ReadByte)
			}
			data = []byte(text)
		default:
			return errors.Errorf("unsupported output type %s", format)
		}
		outFile := filepath.Join(opts.OutputDir, output.MetaData.Table+"."+format)
		err := os.WriteFile(outFile, data, 0644)
		if err != nil {
			return errors.Wrapf(err, "cannot write %s", outFile)
		}
	}
	return nil
}
