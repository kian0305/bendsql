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
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type QueryStatement struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

type MetaData struct {
	Table string `json:"table"`
}

type InputQueryFile struct {
	Statements []QueryStatement `json:"statements"`
	MetaData   MetaData         `json:"metadata"`
}

func (f *InputQueryFile) Decode(b []byte) error {
	err := yaml.Unmarshal(b, f)
	if err != nil {
		return err
	}
	return nil
}

func ReadTargetFiles(directory string) ([]*InputQueryFile, error) {
	res := make([]*InputQueryFile, 0)
	fileInfo, err := os.Stat(directory)
	if err != nil {
		// error handling
		return nil, err
	}

	if !fileInfo.IsDir() {
		// is a directory
		// handle file there
		if strings.HasSuffix(directory, ".yaml") {
			b, err := os.ReadFile(directory)
			if err != nil {
				return nil, err
			}
			var r = &InputQueryFile{}
			err = r.Decode([]byte(b))
			if err != nil {
				log.Printf("error during parsing file %s, %+v", directory, err)
				return nil, err
			}
			res = append(res, r)
		}
		return res, nil
	}
	items, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.IsDir() {
			continue
		} else {
			// handle file there
			if strings.HasSuffix(item.Name(), ".yaml") {
				b, err := os.ReadFile(directory + "/" + item.Name())
				if err != nil {
					return nil, err
				}
				var r = &InputQueryFile{}
				if err := r.Decode(b); err != nil {
					log.Printf("error during parsing file %s, %+v", item.Name(), err)
					return nil, err
				}
				res = append(res, r)
			}
		}
	}
	return res, nil
}
