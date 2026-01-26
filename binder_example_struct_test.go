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
	"time"

	"github.com/xgfone/go-defaults"
	"github.com/xgfone/go-defaults/assists"
)

func ExampleBinder_Struct() {
	type (
		Int    int
		Uint   uint
		String string
		Float  float64
	)

	var S struct {
		Bool    bool
		Int     int
		Int8    int8
		Int16   int16
		Int32   int32
		Int64   int64
		Uint    uint
		Uint8   uint8
		Uint16  uint16
		Uint32  uint32
		Uint64  uint64
		Float32 float32
		Float64 float64
		String  string

		Duration1 time.Duration
		Duration2 time.Duration
		Duration3 time.Duration
		Time1     time.Time
		Time2     time.Time

		Embed struct {
			Int1 int
			Int2 Int

			Uint1 uint
			Uint2 Uint

			String1 string
			String2 String

			Float1 float64
			Float2 Float
		}

		Ignore string `key:"-"`
		Squash struct {
			Field1 int
			Field2 int
		} `key:",squash"`
	}

	maps := map[string]any{
		"Bool": true,

		"Int":   10,
		"Int8":  11,
		"Int16": 12,
		"Int32": 13,
		"Int64": 14,

		"Uint":   20,
		"Uint8":  21,
		"Uint16": 22,
		"Uint32": 23,
		"Uint64": 24,

		"Float32": 30,
		"Float64": 31,

		"String":    "abc",
		"Duration1": "1s",                   // string => time.Duration
		"Duration2": 2000,                   // int(ms) => time.Duration
		"Duration3": 3.0,                    // float(s) => time.Duration
		"Time1":     1672531200,             // int(unix timestamp) => time.Time
		"Time2":     "2023-02-01T00:00:00Z", // string(RFC3339) => time.Time

		"Embed": map[string]any{
			"Int1":    "41", // string => int
			"Int2":    "42", // string => Int
			"Uint1":   "43", // string => uint
			"Uint2":   "44", // string => Uint
			"Float1":  "45", // string => float64
			"Float2":  "46", // string => Float
			"String1": 47,   // int => string
			"String2": 48,   // int => String
		},

		"Ignore": "xyz",
		"Field1": 51,
		"Field2": 52,
	}

	defaults.StructFieldNameFunc.Set(assists.StructFieldNameFuncWithTags("key", "json"))

	err := Bind(&S, maps)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Bool=%v\n", S.Bool)
	fmt.Printf("Int=%v\n", S.Int)
	fmt.Printf("Int8=%v\n", S.Int8)
	fmt.Printf("Int16=%v\n", S.Int16)
	fmt.Printf("Int32=%v\n", S.Int32)
	fmt.Printf("Int64=%v\n", S.Int64)
	fmt.Printf("Uint=%v\n", S.Uint)
	fmt.Printf("Uint8=%v\n", S.Uint8)
	fmt.Printf("Uint16=%v\n", S.Uint16)
	fmt.Printf("Uint32=%v\n", S.Uint32)
	fmt.Printf("Uint64=%v\n", S.Uint64)
	fmt.Printf("Float32=%v\n", S.Float32)
	fmt.Printf("Float64=%v\n", S.Float64)
	fmt.Printf("String=%v\n", S.String)
	fmt.Printf("Duration1=%v\n", S.Duration1)
	fmt.Printf("Duration2=%v\n", S.Duration2)
	fmt.Printf("Duration3=%v\n", S.Duration3)
	fmt.Printf("Time1=%v\n", S.Time1.Format(time.RFC3339))
	fmt.Printf("Time2=%v\n", S.Time2.Format(time.RFC3339))
	fmt.Printf("Embed.Int1=%v\n", S.Embed.Int1)
	fmt.Printf("Embed.Int2=%v\n", S.Embed.Int2)
	fmt.Printf("Embed.Uint1=%v\n", S.Embed.Uint1)
	fmt.Printf("Embed.Uint2=%v\n", S.Embed.Uint2)
	fmt.Printf("Embed.String1=%v\n", S.Embed.String1)
	fmt.Printf("Embed.String2=%v\n", S.Embed.String2)
	fmt.Printf("Embed.Float1=%v\n", S.Embed.Float1)
	fmt.Printf("Embed.Float2=%v\n", S.Embed.Float2)
	fmt.Printf("Squash.Field1=%v\n", S.Squash.Field1)
	fmt.Printf("Squash.Field2=%v\n", S.Squash.Field2)
	fmt.Printf("Ignore=%v\n", S.Ignore)

	// Output:
	// Bool=true
	// Int=10
	// Int8=11
	// Int16=12
	// Int32=13
	// Int64=14
	// Uint=20
	// Uint8=21
	// Uint16=22
	// Uint32=23
	// Uint64=24
	// Float32=30
	// Float64=31
	// String=abc
	// Duration1=1s
	// Duration2=2s
	// Duration3=3s
	// Time1=2023-01-01T00:00:00Z
	// Time2=2023-02-01T00:00:00Z
	// Embed.Int1=41
	// Embed.Int2=42
	// Embed.Uint1=43
	// Embed.Uint2=44
	// Embed.String1=47
	// Embed.String2=48
	// Embed.Float1=45
	// Embed.Float2=46
	// Squash.Field1=51
	// Squash.Field2=52
	// Ignore=
}
