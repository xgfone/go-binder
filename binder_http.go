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
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"reflect"

	"github.com/xgfone/go-structs/field"
)

// BindStructToMap binds the struct to map[string]any.
//
// For the key name, it is case-sensitive.
func BindStructToMap(structptr any, tag string, data map[string]any) (err error) {
	return BindWithTag(structptr, data, tag)
}

// BindStructToStringMap binds the struct to map[string]string.
//
// For the key name, it is case-sensitive.
func BindStructToStringMap(structptr any, tag string, data map[string]string) (err error) {
	return BindWithTag(structptr, data, tag)
}

// BindStructToURLValues binds the struct to url.Values.
//
// For the key name, it is case-sensitive.
func BindStructToURLValues(structptr any, tag string, data url.Values) error {
	return BindWithTag(structptr, data, tag)
}

// BindStructToHTTPHeader binds the struct to http.Header.
//
// For the key name, it will use textproto.CanonicalMIMEHeaderKey(s) to normalize it.
func BindStructToHTTPHeader(structptr any, tag string, data http.Header) error {
	binder := NewBinder()
	binder.GetFieldName = func(sf reflect.StructField) (name, arg string) {
		switch name, arg = field.GetTag(sf, tag); name {
		case "":
			name = textproto.CanonicalMIMEHeaderKey(sf.Name)
		case "-":
			name = ""
		default:
			name = textproto.CanonicalMIMEHeaderKey(name)
		}
		return
	}
	return binder.Bind(structptr, data)
}

// BindStructToMultipartFileHeaders binds the struct to the multipart form file headers.
//
// For the key name, it is case-sensitive.
func BindStructToMultipartFileHeaders(structptr any, tag string, fhs map[string][]*multipart.FileHeader) error {
	return BindWithTag(structptr, fhs, tag)
}
