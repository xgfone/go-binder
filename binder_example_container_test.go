// Copyright 2023 xgfone
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

package binder

import (
	"fmt"
	"net/url"
)

func ExampleBinder_Container() {
	type Ints []int
	var S struct {
		Maps    map[string]interface{} `json:"maps"`
		Slices  []string               `json:"slices"`
		Structs []struct {
			Ints  Ints       `json:"ints"`
			Query url.Values `json:"query"`
		} `json:"structs"`
	}

	maps := map[string]interface{}{
		"maps":   map[string]string{"k11": "v11", "k12": "v12"},
		"slices": []interface{}{"a", "b", "c"},
		"structs": []map[string]interface{}{
			{
				"ints": []string{"21", "22"},
				"query": map[string][]string{
					"k20": {"v21", "v22"},
					"k30": {"v31", "v32"},
				},
			},
			{
				"ints": []int{31, 32},
				"query": map[string][]string{
					"k40": {"v40"},
				},
			},
		},
	}

	err := Bind(&S, maps)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Maps: %v\n", S.Maps)
	fmt.Printf("Slices: %v\n", S.Slices)
	for i, s := range S.Structs {
		fmt.Printf("Structs[%d]: Ints=%v, Query=%v\n", i, s.Ints, s.Query)
	}

	// Output:
	// Maps: map[k11:v11 k12:v12]
	// Slices: [a b c]
	// Structs[0]: Ints=[21 22], Query=map[k20:[v21 v22] k30:[v31 v32]]
	// Structs[1]: Ints=[31 32], Query=map[k40:[v40]]
}
