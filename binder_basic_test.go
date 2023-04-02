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
	"reflect"
	"time"
)

func ExampleBinder_Basic() {
	type (
		IntT    int
		UintT   uint
		FloatT  float64
		StringT string
	)

	var (
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

		intT    IntT
		uintT   UintT
		floatT  FloatT
		stringT StringT
	)

	println := func(dst, src interface{}) {
		err := Bind(dst, src)
		fmt.Println(reflect.ValueOf(dst).Elem().Interface(), err)
	}

	println(&Bool, "true")
	println(&Int, time.Second)
	println(&Int8, 11.0)
	println(&Int16, "12")
	println(&Int32, true)
	println(&Int64, time.Unix(1672531200, 0))
	println(&Uint, 20)
	println(&Uint8, 21.0)
	println(&Uint16, "22")
	println(&Uint32, true)
	println(&Uint64, 23)
	println(&Float32, "1.2")
	println(&Float64, 30)
	println(&String, 40)

	println(&intT, 50.0)
	println(&uintT, IntT(60))
	println(&floatT, StringT("70"))
	println(&stringT, "test")

	// Output:
	// true <nil>
	// 1000 <nil>
	// 11 <nil>
	// 12 <nil>
	// 1 <nil>
	// 1672531200 <nil>
	// 20 <nil>
	// 21 <nil>
	// 22 <nil>
	// 1 <nil>
	// 23 <nil>
	// 1.2 <nil>
	// 30 <nil>
	// 40 <nil>
	// 50 <nil>
	// 60 <nil>
	// 70 <nil>
	// test <nil>
}
