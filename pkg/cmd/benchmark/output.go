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

type OutputMetaData struct {
	Table string `json:"table"`
	Tag   string `json:"tag"`
	Size  string `json:"size"`
}

type OutputSchema struct {
	Name     string    `json:"name"`
	SQL      string    `json:"sql"`
	Min      float64   `json:"min"`
	Max      float64   `json:"max"`
	Median   float64   `json:"median"`
	StdDev   float64   `json:"std_dev"`
	ReadRow  uint64    `json:"read_row"`
	ReadByte uint64    `json:"read_byte"`
	Time     []float64 `json:"time"`
	Error    []string  `json:"error"`
	Mean     float64   `json:"mean"`
}

type OutputFile struct {
	MetaData OutputMetaData `json:"metadata"`
	Schema   []OutputSchema `json:"schema"`
}
