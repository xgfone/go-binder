// Copyright 2026 xgfone
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
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestBinder_Basic(t *testing.T) {
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
	)

	testCases := []struct {
		name     string
		dst      interface{}
		src      interface{}
		expected interface{}
	}{
		{"Bool", &Bool, "true", true},
		{"Int", &Int, time.Second, int(1000)},
		{"Int8", &Int8, 11.0, int8(11)},
		{"Int16", &Int16, "12", int16(12)},
		{"Int32", &Int32, true, int32(1)},
		{"Int64", &Int64, time.Unix(1672531200, 0), int64(1672531200)},
		{"Uint", &Uint, 20, uint(20)},
		{"Uint8", &Uint8, 21.0, uint8(21)},
		{"Uint16", &Uint16, "22", uint16(22)},
		{"Uint32", &Uint32, true, uint32(1)},
		{"Uint64", &Uint64, 23, uint64(23)},
		{"Float32", &Float32, "1.2", float32(1.2)},
		{"Float64", &Float64, 30, float64(30)},
		{"String", &String, 40, "40"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Bind(tc.dst, tc.src)
			if err != nil {
				t.Errorf("Bind failed: %v", err)
				return
			}

			// 获取实际值
			actual := reflect.ValueOf(tc.dst).Elem().Interface()
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, actual)
			}
		})
	}
}

func TestBinder_Struct(t *testing.T) {
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
	}

	err := Bind(&S, maps)
	if err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	// 验证所有字段
	if S.Bool != true {
		t.Errorf("Expected Bool=true, got %v", S.Bool)
	}
	if S.Int != 10 {
		t.Errorf("Expected Int=10, got %v", S.Int)
	}
	if S.Int8 != 11 {
		t.Errorf("Expected Int8=11, got %v", S.Int8)
	}
	if S.Int16 != 12 {
		t.Errorf("Expected Int16=12, got %v", S.Int16)
	}
	if S.Int32 != 13 {
		t.Errorf("Expected Int32=13, got %v", S.Int32)
	}
	if S.Int64 != 14 {
		t.Errorf("Expected Int64=14, got %v", S.Int64)
	}
	if S.Uint != 20 {
		t.Errorf("Expected Uint=20, got %v", S.Uint)
	}
	if S.Uint8 != 21 {
		t.Errorf("Expected Uint8=21, got %v", S.Uint8)
	}
	if S.Uint16 != 22 {
		t.Errorf("Expected Uint16=22, got %v", S.Uint16)
	}
	if S.Uint32 != 23 {
		t.Errorf("Expected Uint32=23, got %v", S.Uint32)
	}
	if S.Uint64 != 24 {
		t.Errorf("Expected Uint64=24, got %v", S.Uint64)
	}
	if S.Float32 != 30 {
		t.Errorf("Expected Float32=30, got %v", S.Float32)
	}
	if S.Float64 != 31 {
		t.Errorf("Expected Float64=31, got %v", S.Float64)
	}
	if S.String != "abc" {
		t.Errorf("Expected String=abc, got %v", S.String)
	}
	if S.Duration1 != time.Second {
		t.Errorf("Expected Duration1=1s, got %v", S.Duration1)
	}
	if S.Duration2 != 2*time.Second {
		t.Errorf("Expected Duration2=2s, got %v", S.Duration2)
	}
	if S.Duration3 != 3*time.Second {
		t.Errorf("Expected Duration3=3s, got %v", S.Duration3)
	}
	expectedTime1 := time.Unix(1672531200, 0)
	if !S.Time1.Equal(expectedTime1) {
		t.Errorf("Expected Time1=%v, got %v", expectedTime1.Format(time.RFC3339), S.Time1.Format(time.RFC3339))
	}
	expectedTime2, _ := time.Parse(time.RFC3339, "2023-02-01T00:00:00Z")
	if !S.Time2.Equal(expectedTime2) {
		t.Errorf("Expected Time2=%v, got %v", expectedTime2.Format(time.RFC3339), S.Time2.Format(time.RFC3339))
	}
	if S.Embed.Int1 != 41 {
		t.Errorf("Expected Embed.Int1=41, got %v", S.Embed.Int1)
	}
	if S.Embed.Int2 != 42 {
		t.Errorf("Expected Embed.Int2=42, got %v", S.Embed.Int2)
	}
	if S.Embed.Uint1 != 43 {
		t.Errorf("Expected Embed.Uint1=43, got %v", S.Embed.Uint1)
	}
	if S.Embed.Uint2 != 44 {
		t.Errorf("Expected Embed.Uint2=44, got %v", S.Embed.Uint2)
	}
	if S.Embed.String1 != "47" {
		t.Errorf("Expected Embed.String1=47, got %v", S.Embed.String1)
	}
	if S.Embed.String2 != "48" {
		t.Errorf("Expected Embed.String2=48, got %v", S.Embed.String2)
	}
	if S.Embed.Float1 != 45 {
		t.Errorf("Expected Embed.Float1=45, got %v", S.Embed.Float1)
	}
	if S.Embed.Float2 != 46 {
		t.Errorf("Expected Embed.Float2=46, got %v", S.Embed.Float2)
	}
}

func TestBinder_Container(t *testing.T) {
	type Ints []int
	var S struct {
		Maps    map[string]any `json:"maps"`
		Slices  []string       `json:"slices"`
		Structs []struct {
			Ints  Ints       `json:"ints"`
			Query url.Values `json:"query"`
		} `json:"structs"`
	}

	maps := map[string]any{
		"maps":   map[string]string{"k11": "v11", "k12": "v12"},
		"slices": []any{"a", "b", "c"},
		"structs": []map[string]any{
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
		t.Fatalf("Bind failed: %v", err)
	}

	// 验证 Maps
	expectedMaps := map[string]any{"k11": "v11", "k12": "v12"}
	if !reflect.DeepEqual(S.Maps, expectedMaps) {
		t.Errorf("Expected Maps: %v, got %v", expectedMaps, S.Maps)
	}

	// 验证 Slices
	expectedSlices := []string{"a", "b", "c"}
	if !reflect.DeepEqual(S.Slices, expectedSlices) {
		t.Errorf("Expected Slices: %v, got %v", expectedSlices, S.Slices)
	}

	// 验证 Structs
	if len(S.Structs) != 2 {
		t.Fatalf("Expected 2 structs, got %d", len(S.Structs))
	}

	// 第一个结构体
	if !reflect.DeepEqual(S.Structs[0].Ints, Ints{21, 22}) {
		t.Errorf("Expected Structs[0].Ints=[21 22], got %v", S.Structs[0].Ints)
	}
	expectedQuery1 := url.Values{
		"k20": {"v21", "v22"},
		"k30": {"v31", "v32"},
	}
	if !reflect.DeepEqual(S.Structs[0].Query, expectedQuery1) {
		t.Errorf("Expected Structs[0].Query=%v, got %v", expectedQuery1, S.Structs[0].Query)
	}

	// 第二个结构体
	if !reflect.DeepEqual(S.Structs[1].Ints, Ints{31, 32}) {
		t.Errorf("Expected Structs[1].Ints=[31 32], got %v", S.Structs[1].Ints)
	}
	expectedQuery2 := url.Values{"k40": {"v40"}}
	if !reflect.DeepEqual(S.Structs[1].Query, expectedQuery2) {
		t.Errorf("Expected Structs[1].Query=%v, got %v", expectedQuery2, S.Structs[1].Query)
	}
}

func TestBinder_Hook(t *testing.T) {
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
	multiparthook := func(dst reflect.Value, src any) (any, error) {
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

	err := Binder{Hook: multiparthook}.Bind(&dst, src)
	if err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	// 验证 File
	if dst.File == nil || dst.File.Filename != "file" {
		t.Errorf("Expected File.Filename=file, got %v", dst.File)
	}

	// 验证 Files
	if len(dst.Files) != 2 {
		t.Fatalf("Expected 2 files, got %d", len(dst.Files))
	}
	if dst.Files[0].Filename != "file1" {
		t.Errorf("Expected Files[0].Filename=file1, got %s", dst.Files[0].Filename)
	}
	if dst.Files[1].Filename != "file2" {
		t.Errorf("Expected Files[1].Filename=file2, got %s", dst.Files[1].Filename)
	}
}

// Int is the customized int.
type testInt int

// Set implements the interface Setter.
func (i *testInt) Set(src any) (err error) {
	switch v := src.(type) {
	case int:
		*i = testInt(v)
	case string:
		var _v int64
		_v, err = strconv.ParseInt(v, 10, 64)
		if err == nil {
			*i = testInt(_v)
		}
	default:
		err = fmt.Errorf("unsupport to convert %T to Int", src)
	}
	return
}

func (i testInt) String() string {
	return fmt.Sprint(int64(i))
}

// Struct is the customized struct.
type testStruct struct {
	Name string
	Age  testInt
}

// UnmarshalBind implements the interface Unmarshaler.
func (s *testStruct) UnmarshalBind(src any) (err error) {
	switch v := src.(type) {
	case string:
		items := strings.Split(src.(string), ";")

		var age int64
		age, err = strconv.ParseInt(items[1], 10, 64)

		s.Age = testInt(age)
		s.Name = items[0]

	case map[string]any:
		s.Name, _ = v["Name"].(string)
		err = s.Age.Set(v["Age"])

	default:
		err = fmt.Errorf("unsupport to convert %T to a struct", src)
	}

	return
}

func (s testStruct) String() string {
	return fmt.Sprintf("Name=%s, Age=%d", s.Name, s.Age)
}

func TestBinder_Interface(t *testing.T) {
	var iface1 testInt
	var iface2 testStruct
	var S = struct {
		Interface1 Setter
		Interface2 Unmarshaler

		Interface3 error
		Interface4 *error

		Interface5 any // Use to store any type value.
		// Unmarshaler         // Do not use the anonymous interface.

		Interface6 testStruct
	}{
		Interface1: &iface1, // For interface, must be set to a pointer
		Interface2: &iface2, //  to an implementation.
	}

	iface3 := errors.New("test1")
	iface4 := errors.New("test2")
	maps := map[string]any{
		"Interface1": "123",
		"Interface2": map[string]any{"Name": "Aaron", "Age": 18},
		"Interface3": iface3,
		"Interface4": iface4,
		"Interface5": "any",
		"Interface6": "Xgfone;20",
	}

	err := Bind(&S, maps)
	if err != nil {
		t.Fatalf("Bind failed: %v", err)
	}

	// 验证 Interface1
	if iface1 != 123 {
		t.Errorf("Expected Interface1=123, got %v", iface1)
	}

	// 验证 Interface2
	if iface2.Name != "Aaron" || iface2.Age != 18 {
		t.Errorf("Expected Interface2: Name=Aaron, Age=18, got %v", iface2)
	}

	// 验证 Interface3
	if S.Interface3 == nil || S.Interface3.Error() != "test1" {
		t.Errorf("Expected Interface3=test1, got %v", S.Interface3)
	}

	// 验证 Interface4
	if S.Interface4 == nil || (*S.Interface4).Error() != "test2" {
		t.Errorf("Expected Interface4=test2, got %v", S.Interface4)
	}

	// 验证 Interface5
	if S.Interface5 != "any" {
		t.Errorf("Expected Interface5=any, got %v", S.Interface5)
	}

	// 验证 Interface6
	if S.Interface6.Name != "Xgfone" || S.Interface6.Age != 20 {
		t.Errorf("Expected Interface6: Name=Xgfone, Age=20, got %v", S.Interface6)
	}
}

func TestBindStructToStringMap(t *testing.T) {
	src := map[string]string{
		"Int": "123",
		"Str": "456",
	}

	var dst struct {
		Int  int `tag:"-"`
		Int1 int `tag:"Int"`
		Int2 int `tag:"Str"`
	}

	err := BindStructToStringMap(&dst, "tag", src)
	if err != nil {
		t.Fatalf("BindStructToStringMap failed: %v", err)
	}

	if dst.Int != 0 {
		t.Errorf("Expected Int=0, got %d", dst.Int)
	}
	if dst.Int1 != 123 {
		t.Errorf("Expected Int1=123, got %d", dst.Int1)
	}
	if dst.Int2 != 456 {
		t.Errorf("Expected Int2=456, got %d", dst.Int2)
	}
}

func TestBindStructToHTTPHeader(t *testing.T) {
	src := http.Header{
		"X-Int":  []string{"1", "2"},
		"X-Ints": []string{"3", "4"},
		"X-Str":  []string{"a", "b"},
		"X-Strs": []string{"c", "d"},
	}

	var dst struct {
		unexported string   `header:"-"`
		Other      string   `header:"Other"`
		Int        int      `header:"x-int"`
		Ints       []int    `header:"x-ints"`
		Str        string   `header:"x-str"`
		Strs       []string `header:"x-strs"`
	}

	err := BindStructToHTTPHeader(&dst, "header", src)
	if err != nil {
		t.Fatalf("BindStructToHTTPHeader failed: %v", err)
	}

	if dst.unexported != "" {
		t.Errorf("Expected unexported=, got %s", dst.unexported)
	}
	if dst.Other != "" {
		t.Errorf("Expected Other=, got %s", dst.Other)
	}
	if dst.Int != 1 {
		t.Errorf("Expected Int=1, got %d", dst.Int)
	}
	if !reflect.DeepEqual(dst.Ints, []int{3, 4}) {
		t.Errorf("Expected Ints=[3 4], got %v", dst.Ints)
	}
	if dst.Str != "a" {
		t.Errorf("Expected Str=a, got %s", dst.Str)
	}
	if !reflect.DeepEqual(dst.Strs, []string{"c", "d"}) {
		t.Errorf("Expected Strs=[c d], got %v", dst.Strs)
	}
}

func TestBindStructToURLValues(t *testing.T) {
	src := url.Values{
		"int":  []string{"1", "2"},
		"ints": []string{"3", "4"},
		"str":  []string{"a", "b"},
		"strs": []string{"c", "d"},
	}

	var dst struct {
		unexported string   `qeury:"-"`
		Other      string   `query:"Other"`
		Int        int      `query:"int"`
		Ints       []int    `query:"ints"`
		Str        string   `query:"str"`
		Strs       []string `query:"strs"`
	}

	err := BindStructToURLValues(&dst, "query", src)
	if err != nil {
		t.Fatalf("BindStructToURLValues failed: %v", err)
	}

	if dst.unexported != "" {
		t.Errorf("Expected unexported=, got %s", dst.unexported)
	}
	if dst.Other != "" {
		t.Errorf("Expected Other=, got %s", dst.Other)
	}
	if dst.Int != 1 {
		t.Errorf("Expected Int=1, got %d", dst.Int)
	}
	if !reflect.DeepEqual(dst.Ints, []int{3, 4}) {
		t.Errorf("Expected Ints=[3 4], got %v", dst.Ints)
	}
	if dst.Str != "a" {
		t.Errorf("Expected Str=a, got %s", dst.Str)
	}
	if !reflect.DeepEqual(dst.Strs, []string{"c", "d"}) {
		t.Errorf("Expected Strs=[c d], got %v", dst.Strs)
	}
}

func TestBindStructToMultipartFileHeaders(t *testing.T) {
	src := map[string][]*multipart.FileHeader{
		"file":  {{Filename: "file"}},
		"files": {{Filename: "file1"}, {Filename: "file2"}},
		"_file": {{Filename: "file3"}},
	}

	var dst struct {
		Other       string                  `form:"Other"`
		_File       *multipart.FileHeader   `form:"_file"` // unexported, so ignored
		FileHeader  *multipart.FileHeader   `form:"file"`
		FileHeaders []*multipart.FileHeader `form:"files"`
	}

	err := BindStructToMultipartFileHeaders(&dst, "form", src)
	if err != nil {
		t.Fatalf("BindStructToMultipartFileHeaders failed: %v", err)
	}

	// 验证 FileHeader
	if dst.FileHeader == nil || dst.FileHeader.Filename != "file" {
		t.Errorf("Expected FileHeader.Filename=file, got %v", dst.FileHeader)
	}

	// 验证 _File 应该是 nil（未导出字段）
	if dst._File != nil {
		t.Errorf("Expected _File=nil, got %v", dst._File)
	}

	// 验证 FileHeaders
	if len(dst.FileHeaders) != 2 {
		t.Fatalf("Expected 2 FileHeaders, got %d", len(dst.FileHeaders))
	}
	if dst.FileHeaders[0].Filename != "file1" {
		t.Errorf("Expected FileHeaders[0].Filename=file1, got %s", dst.FileHeaders[0].Filename)
	}
	if dst.FileHeaders[1].Filename != "file2" {
		t.Errorf("Expected FileHeaders[1].Filename=file2, got %s", dst.FileHeaders[1].Filename)
	}
}
