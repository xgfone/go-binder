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
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/xgfone/go-defaults"
	"github.com/xgfone/go-defaults/assists"
	"github.com/xgfone/go-structs"
	"github.com/xgfone/go-validation"
)

// Predefine some decoders to decode a value,
// such as body, query and header of the http request.
var (
	// It only supports to decode *http.Request with the tag "query" by default.
	DefaultQueryDecoder Decoder = DecoderFunc(func(dst, src interface{}) error {
		if req, ok := src.(*http.Request); ok {
			return BindStructToURLValues(dst, "query", req.URL.Query())
		}
		return fmt.Errorf("binder.DefaultQueryDecoder: unsupport to decode %T", src)
	})

	// It only supports to decode *http.Request with the tag "header" by default.
	DefaultHeaderDecoder Decoder = DecoderFunc(func(dst, src interface{}) error {
		if req, ok := src.(*http.Request); ok {
			return BindStructToHTTPHeader(dst, "header", req.Header)
		}
		return fmt.Errorf("binder.DefaultHeaderDecoder: unsupport to decode %T", src)
	})

	// By default, during initializing the package, it will register
	// some decoders for the http request with the content-types:
	//   - "application/xml"
	//   - "application/json"
	//   - "multipart/form-data"
	//   - "application/x-www-form-urlencoded"
	// For the http request, it can be used like
	//   DefaultMuxDecoder.Decode(dst, httpRequest).
	DefaultMuxDecoder = NewMuxDecoder()

	// It will use defaults.ValidateStruct to validate the struct value by default.
	DefaultStructValidationDecoder Decoder = StructValidationDecoder(nil)

	// Some encapsulated http decoders, which can be used directly.
	BodyDecoder   Decoder = ComposeDecoders(DefaultMuxDecoder, DefaultStructValidationDecoder)
	QueryDecoder  Decoder = ComposeDecoders(DefaultQueryDecoder, DefaultStructValidationDecoder)
	HeaderDecoder Decoder = ComposeDecoders(DefaultHeaderDecoder, DefaultStructValidationDecoder)
)

func init() {
	if defaults.RuleValidator.Get() == nil {
		defaults.RuleValidator.Set(assists.RuleValidateFunc(validation.Validate))
	}
	if defaults.StructValidator.Get() == nil {
		defaults.StructValidator.Set(assists.StructValidateFunc(structs.Reflect))
	}

	DefaultMuxDecoder.Add("application/json", DecoderFunc(func(dst, src interface{}) error {
		if req := src.(*http.Request); req.ContentLength > 0 {
			return json.NewDecoder(req.Body).Decode(dst)
		}
		return nil
	}))
}

func init() {
	DefaultMuxDecoder.Add("application/xml", DecoderFunc(func(dst, src interface{}) error {
		if req := src.(*http.Request); req.ContentLength > 0 {
			return xml.NewDecoder(req.Body).Decode(dst)
		}
		return nil
	}))
}

func init() {
	registerFormDecoder("multipart/form-data")
	registerFormDecoder("application/x-www-form-urlencoded")
}

func registerFormDecoder(ct string) {
	const maxMemory = 10 << 20
	DefaultMuxDecoder.Add(ct, DecoderFunc(func(dst, src interface{}) (err error) {
		req := src.(*http.Request)
		switch ct := getContentType(req.Header); ct {
		case "multipart/form-data":
			err = req.ParseMultipartForm(maxMemory)

		case "application/x-www-form-urlencoded":
			err = req.ParseForm()

		default:
			return fmt.Errorf("unsupported Content-Type '%s'", ct)
		}

		if err != nil {
			return
		}

		err = BindStructToURLValues(dst, "form", req.Form)
		if err == nil && len(req.MultipartForm.File) > 0 {
			err = BindStructToMultipartFileHeaders(dst, "form", req.MultipartForm.File)
		}

		return
	}))
}
