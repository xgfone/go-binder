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
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Int is the customized int.
type Int int

// Set implements the interface Setter.
func (i *Int) Set(src any) (err error) {
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

// UnmarshalBind implements the interface Unmarshaler.
func (s *Struct) UnmarshalBind(src any) (err error) {
	switch v := src.(type) {
	case string:
		items := strings.Split(src.(string), ";")

		var age int64
		age, err = strconv.ParseInt(items[1], 10, 64)

		s.Age = Int(age)
		s.Name = items[0]

	case map[string]any:
		s.Name, _ = v["Name"].(string)
		err = s.Age.Set(v["Age"])

	default:
		err = fmt.Errorf("unsupport to convert %T to a struct", src)
	}

	return
}

func (s Struct) String() string {
	return fmt.Sprintf("Name=%s, Age=%d", s.Name, s.Age)
}

func ExampleBinder_Interface() {
	var iface1 Int
	var iface2 Struct
	var S = struct {
		Interface1 Setter
		Interface2 Unmarshaler

		Interface3 error
		Interface4 *error

		Interface5 any // Use to store any type value.
		// Unmarshaler         // Do not use the anonymous interface.

		Interface6 Struct
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
		fmt.Println(err)
		return
	}

	fmt.Printf("Interface1: %v\n", S.Interface1)
	fmt.Printf("Interface2: %v\n", S.Interface2)
	fmt.Printf("Interface3: %v\n", S.Interface3)
	fmt.Printf("Interface4: %v\n", *S.Interface4)
	fmt.Printf("Interface5: %v\n", S.Interface5)
	fmt.Printf("Interface6: %v\n", S.Interface6)

	// Output:
	// Interface1: 123
	// Interface2: Name=Aaron, Age=18
	// Interface3: test1
	// Interface4: test2
	// Interface5: any
	// Interface6: Name=Xgfone, Age=20
}
