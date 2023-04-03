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
	"net/http"
	"strings"

	"github.com/xgfone/go-defaults"
	"github.com/xgfone/go-defaults/assists"
)

var errMissingContentType = errors.New("missing the header Content-Type")

// Decoder is used to decode the data src to dst.
//
// In general, Deocder is used to decode a byte stream to a type,
// such as struct or map.
type Decoder interface {
	Decode(dst, src interface{}) error
}

// DecoderFunc is a function to decode the data src to dst.
type DecoderFunc func(dst, src interface{}) error

// Decode implements the interface Decoder.
func (f DecoderFunc) Decode(dst, src interface{}) error { return f(dst, src) }

// ComposeDecoders composes a group of decoders, which will be called in turn,
// to a Decoder.
func ComposeDecoders(decoders ...Decoder) Decoder {
	if len(decoders) == 0 {
		panic("ComposeDecoders: missing decoders")
	}

	return DecoderFunc(func(dst, src interface{}) (err error) {
		for _, decoder := range decoders {
			if err = decoder.Decode(dst, src); err != nil {
				return
			}
		}
		return
	})
}

// StructValidationDecoder returns a struct validation decoder,
// which only validates whether the value dst is valid, not decodes any.
func StructValidationDecoder(validator assists.StructValidator) Decoder {
	validate := defaults.ValidateStruct
	if validator != nil {
		validate = validator.Validate
	}

	return DecoderFunc(func(dst, src interface{}) (err error) {
		return validate(dst)
	})
}

// MuxDecoder is a multiplexer for kinds of Decoders.
type MuxDecoder struct {
	// GetDecoder is used to get the deocder by the funciton get
	// with decoder that comes from src.
	//
	// If nil, use the default implementation, which inspects the decoder type
	// by the type of src, and supports the types as follow:
	//   *http.Request: => Content-Type
	//   interface{ DecodeType() string }
	//   interface{ Type() string }
	GetDecoder func(src interface{}, get func(string) Decoder) (Decoder, error)

	decoders map[string]Decoder
}

// NewMuxDecoder returns a new MuxDecoder.
func NewMuxDecoder() *MuxDecoder {
	return &MuxDecoder{decoders: make(map[string]Decoder, 8)}
}

// Add adds a decoder to decode the data of the given type.
func (md *MuxDecoder) Add(dtype string, decoder Decoder) {
	md.decoders[dtype] = decoder
}

// Del removes the corresponding decoder by the type.
func (md *MuxDecoder) Del(dtype string) { delete(md.decoders, dtype) }

// Get returns the corresponding decoder by the type.
//
// Return nil if not found.
func (md *MuxDecoder) Get(dtype string) Decoder { return md.decoders[dtype] }

// Decode implements the interface Decoder.
func (md *MuxDecoder) Decode(dst, src interface{}) (err error) {
	var decoder Decoder
	if md.GetDecoder != nil {
		decoder, err = md.GetDecoder(src, md.Get)
	} else {
		decoder, err = md.getDecoder(src, md.Get)
	}
	if err == nil {
		err = decoder.Decode(dst, src)
	}
	return
}

func (md *MuxDecoder) getDecoder(src interface{}, get func(string) Decoder) (Decoder, error) {
	switch req := src.(type) {
	case *http.Request:
		ct := getContentType(req.Header)
		if ct == "" {
			return nil, errMissingContentType
		}
		if decoder := get(ct); decoder != nil {
			return decoder, nil
		}
		return nil, fmt.Errorf("unsupported Content-Type '%s'", ct)

	case interface{ DecodeType() string }:
		dtype := req.DecodeType()
		if decoder := get(dtype); decoder != nil {
			return decoder, nil
		}
		return nil, fmt.Errorf("unsupported request data type '%s'", dtype)

	case interface{ Type() string }:
		dtype := req.Type()
		if decoder := get(dtype); decoder != nil {
			return decoder, nil
		}
		return nil, fmt.Errorf("unsupported request data type '%s'", dtype)

	default:
		return nil, fmt.Errorf("unknown request data type %T", src)
	}
}

func getContentType(header http.Header) string {
	ct := header.Get("Content-Type")
	if index := strings.IndexByte(ct, ';'); index > -1 {
		ct = strings.TrimSpace(ct[:index])
	}
	return ct
}
