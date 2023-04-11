# Go Binder [![Build Status](https://github.com/xgfone/go-binder/actions/workflows/go.yml/badge.svg)](https://github.com/xgfone/go-binder/actions/workflows/go.yml) [![GoDoc](https://pkg.go.dev/badge/github.com/xgfone/go-binder)](https://pkg.go.dev/github.com/xgfone/go-binder) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://raw.githubusercontent.com/xgfone/go-binder/master/LICENSE)

Provide a common binder to bind a value to any, for example, binding a struct to a map.

For the struct, the package registers `github.com/xgfone/go-structs.Reflect` to reflect the fields of struct to validate the struct value and `github.com/xgfone/go-validation.Validate` to validate the struct field based on the built rule.


## Install
```shell
$ go get -u github.com/xgfone/go-binder
```


## Example

### Bind the basic types
```go
package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/xgfone/go-binder"
)

func main() {
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
		err := binder.Bind(dst, src)
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
```

### Bind the struct
```go
package main

import (
	"fmt"
	"time"

	"github.com/xgfone/go-binder"
)

func main() {
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

		Ignore string `json:"-"`
		Squash struct {
			Field1 int
			Field2 int
		} `json:",squash"`
	}

	maps := map[string]interface{}{
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

		"Embed": map[string]interface{}{
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

	err := binder.Bind(&S, maps)
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
```

### Bind the containers
```go
package main

import (
	"fmt"
	"net/url"

	"github.com/xgfone/go-binder"
)

func main() {
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

	err := binder.Bind(&S, maps)
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
```

### Bind the interfaces
```go
package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/xgfone/go-binder"
)

// Int is the customized int.
type Int int

// Set implements the interface binder.Setter.
func (i *Int) Set(src interface{}) (err error) {
	switch v := src.(type) {
	case int:
		*i = Int(v)
	case string:
		var _v int64
		_v, err = strconv.ParseInt(v, 10, 64)
		if err == nil {
			*i = Int(_v)
		}
	default:
		err = fmt.Errorf("unsupport to convert %T to Int", src)
	}
	return
}

func (i Int) String() string {
	return fmt.Sprint(int64(i))
}

// Struct is the customized struct.
type Struct struct {
	Name string
	Age  Int
}

// UnmarshalBind implements the interface binder.Unmarshaler.
func (s *Struct) UnmarshalBind(src interface{}) (err error) {
	if maps, ok := src.(map[string]interface{}); ok {
		s.Name, _ = maps["Name"].(string)
		err = s.Age.Set(maps["Age"])
		return
	}
	return fmt.Errorf("unsupport to convert %T to a struct", src)
}

func (s Struct) String() string {
	return fmt.Sprintf("Name=%s, Age=%d", s.Name, s.Age)
}

func main() {
	var iface1 Int
	var iface2 Struct
	var S = struct {
		Interface1 binder.Setter
		Interface2 binder.Unmarshaler

		Interface3 error
		Interface4 *error

		Interface5 interface{} // Use to store any type value.
		// binder.Unmarshaler  // Do not use the anonymous interface.
	}{
		Interface1: &iface1, // For interface, must be set to a pointer
		Interface2: &iface2, //  to an implementation.
	}

	iface3 := errors.New("test1")
	iface4 := errors.New("test2")
	maps := map[string]interface{}{
		"Interface1": "123",
		"Interface2": map[string]interface{}{"Name": "Aaron", "Age": 18},
		"Interface3": iface3,
		"Interface4": iface4,
		"Interface5": "any",
	}

	err := binder.Bind(&S, maps)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Interface1: %v\n", S.Interface1)
	fmt.Printf("Interface2: %v\n", S.Interface2)
	fmt.Printf("Interface3: %v\n", S.Interface3)
	fmt.Printf("Interface4: %v\n", *S.Interface4)
	fmt.Printf("Interface5: %v\n", S.Interface5)

	// Output:
	// Interface1: 123
	// Interface2: Name=Aaron, Age=18
	// Interface3: test1
	// Interface4: test2
	// Interface5: any
}
```

### Hook
```go
package main

import (
	"fmt"
	"mime/multipart"
	"reflect"

	"github.com/xgfone/go-binder"
)

func main() {
	src := map[string][]*multipart.FileHeader{
		"file":  {{Filename: "file"}},
		"files": {{Filename: "file1"}, {Filename: "file2"}},
	}

	var dst struct {
		File  *multipart.FileHeader   `json:"file"`
		Files []*multipart.FileHeader `json:"files"`
	}

	// (xgf) By default, the binder cannot bind *multipart.FileHeader
	// to []*multipart.FileHeader. However, we can use hook to do it.
	// Here, there are two ways to finish it:
	//   1. We just convert []*multipart.FileHeader to *multipart.FileHeader,
	//      then let the binder continue to finish binding.
	//   2. We finish the binding in the hook.
	//   3. Set ConvertSliceToSingle to true to enable the auto-conversion.
	// In the exmaple, we use the first.
	multiparthook := func(dst reflect.Value, src interface{}) (interface{}, error) {
		if _, ok := dst.Interface().(*multipart.FileHeader); !ok {
			return src, nil // Let the binder continue to handle it.
		}

		srcfiles, ok := src.([]*multipart.FileHeader)
		if !ok {
			return src, nil // Let the binder continue to handle it.
		} else if len(srcfiles) == 0 {
			return nil, nil // FileHeader is empty, we tell the binder not to do it.
		}
		return srcfiles[0], nil
	}

	err := binder.Binder{Hook: multiparthook}.Bind(&dst, src)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("File.Filename=%s\n", dst.File.Filename)
	for i, file := range dst.Files {
		fmt.Printf("Files[%d].Filename=%s\n", i, file.Filename)
	}

	// Output:
	// File.Filename=file
	// Files[0].Filename=file1
	// Files[1].Filename=file2
}
```

### Bind HTTP Request Body
```go
func main() {
	http.HandleFunc("/path", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Field int `json:"field"`
			// ...
		}
		err := binder.BodyDecoder.Decode(&body, r)
		// ...
	})
}
```

### Bind HTTP Request Query
```go
func main() {
	http.HandleFunc("/path", func(w http.ResponseWriter, r *http.Request) {
		var query struct {
			Field int `query:"field"`
			// ...
		}
		err := binder.QueryDecoder.Decode(&query, r)
		// ...
	})
}
```

### Bind HTTP Request Header
```go
func main() {
	http.HandleFunc("/path", func(w http.ResponseWriter, r *http.Request) {
		var header struct {
			Field int `header:"x-field"` // or "X-Field"
			// ...
		}
		err := binder.HeaderDecoder.Decode(&header, r)
		// ...
	})
}
```
