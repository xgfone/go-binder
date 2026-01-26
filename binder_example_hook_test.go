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
	"mime/multipart"
	"reflect"
)

func ExampleBinder_Hook() {
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
