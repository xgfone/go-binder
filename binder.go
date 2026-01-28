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

// Package binder provides a common binder to bind a value to any,
// for example, binding a struct to a map.
package binder

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/xgfone/go-cast"
	"github.com/xgfone/go-structs/field"
)

// DefaultBinder is the default binder.
var DefaultBinder = NewBinder()

// Unmarshaler is an interface to unmarshal itself from the parameter.
type Unmarshaler interface {
	UnmarshalBind(any) error
}

// Setter is an interface to set itself to the parameter.
type Setter interface {
	Set(any) error
}

// Bind uses DefaultBinder to bind dstptr to src.
func Bind(dstptr, src any) error {
	return DefaultBinder.Bind(dstptr, src)
}

// BindWithTag is used to bind dstptr to src,
// which uses the given tag to try to get the field name.
func BindWithTag(dstptr, src any, tag string) error {
	binder := NewBinder()
	binder.GetFieldName = func(sf reflect.StructField) (string, string) {
		return getStructFieldNameWithTag(sf, tag)
	}
	return binder.Bind(dstptr, src)
}

// Hook is used to intercept the binding operation.
type Hook func(dst reflect.Value, src any) (newsrc any, err error)

// Binder is a common binder to bind a value to any.
//
// In general, Binder is used to transform a value between different types.
type Binder struct {
	// If true, convert src from slice/array, that's the first element,
	// to a single value on demand by the bound value.
	ConvertSliceToSingle bool

	// if true, convert src from a single value to slice/array on demand
	// by the bound value.
	ConvertSingleToSlice bool

	// GetFieldName is used to get the name and arg of the given field.
	//
	// If ignoring the field, return the empty string for the field name.
	// For the tag value, it maybe contain the argument, just like
	//   type S struct {
	//       OnlyName    int `json:"fieldname"`
	//       OnlyArg     int `json:",fieldarg"`
	//       NameAndArg  int `json:"fieldname,fieldarg"`
	//       NameAndArgs int `json:"fieldname,fieldarg1,fieldarg2"`
	//       Ignore1     int `json:"-"`
	//       Ignore2     int `json:"-,"`
	//   }
	//
	// For the field argument, it only supports "squash" to squash
	// all the fields of the struct, just like the anonymous field.
	GetFieldName func(reflect.StructField) (name, arg string)

	// Hook is used to intercept the binding operation if set.
	//
	// If newsrc is not nil, the engine will continue to handle it.
	// Or, ignore it and go on to bind the next value.
	// So, if the hook has bound the value, return (nil, nil).
	//
	// Default: nil
	Hook Hook
}

// NewBinder returns a default binder.
func NewBinder() Binder { return NewBinderWithHook(nil) }

// NewBinderWithHook returns a default binder with the hook.
func NewBinderWithHook(hook Hook) Binder {
	return Binder{
		ConvertSliceToSingle: true,
		ConvertSingleToSlice: true,
		Hook:                 hook,
	}
}

// Bind is used to bind the value dstptr to src.
//
// In general, dstptr is a pointer to a contain variable.
// Moreover, dstptr may be a reflect.Value, but it must can be set
// or a pointer that the element can be set.
//
// Support the types of the struct fields as follow:
//
//   - ~bool
//   - ~int
//   - ~int8
//   - ~int16
//   - ~int32
//   - ~int64
//   - ~uint
//   - ~uint8
//   - ~uint16
//   - ~uint32
//   - ~uint64
//   - ~string
//   - ~float32
//   - ~float64
//   - ~Array[E]
//   - ~Slice[E]
//   - ~Map[E]V
//   - time.Time
//   - time.Duration
//   - Struct
//
// And any pointer to the types above, and the interfaces Unmarshaler and Setter.
func (b Binder) Bind(dstptr, src any) error {
	return binder{b.fieldNameGetter(), b}.Bind(dstptr, src)
}

func (b Binder) fieldNameGetter() func(reflect.StructField) (string, string) {
	if b.GetFieldName != nil {
		return b.GetFieldName
	}
	return getStructFieldName
}

func getStructFieldName(sf reflect.StructField) (name string, arg string) {
	return getStructFieldNameWithTag(sf, "json")
}

func getStructFieldNameWithTag(sf reflect.StructField, tag string) (name string, arg string) {
	name, arg = field.GetTag(sf, tag)
	switch name {
	case "":
		name = sf.Name

	case "-":
		name = ""
	}
	return
}

type binder struct {
	getFieldName func(reflect.StructField) (name, arg string)
	Binder
}

func (b binder) Bind(dst, src any) error {
	dstValue, ok := dst.(reflect.Value)
	if !ok {
		dstValue = reflect.ValueOf(dst)
	}

	switch {
	case dstValue.CanSet():
	case dstValue.Kind() == reflect.Pointer:
		if dstValue = dstValue.Elem(); !dstValue.CanSet() {
			return fmt.Errorf("%T must be a pointer to a value that can be set", dst)
		}
	default:
		return fmt.Errorf("Binder.Bind: %T must be canset or a pointer", dst)
	}

	return b.bind(dstValue.Kind(), dstValue, src)
}

func (b binder) bind(kind reflect.Kind, value reflect.Value, src any) (err error) {
	if src == nil {
		return
	}

	if !value.CanSet() {
		switch kind {
		case reflect.Pointer, reflect.Interface:
			if !value.Elem().CanAddr() {
				return
			}
		default:
			return
		}
	}

	if b.Hook != nil {
		src, err = b.Hook(value, src)
		if err != nil || src == nil {
			return err
		}
	}

	if b.ConvertSliceToSingle && kind != reflect.Array && kind != reflect.Slice {
		switch srcValue := reflect.ValueOf(src); srcValue.Kind() {
		case reflect.Slice, reflect.Array:
			if srcValue.Len() == 0 {
				return
			}
			src = srcValue.Index(0).Interface()
		}
	}

	ptrvalue := value
	if kind != reflect.Pointer {
		ptrvalue = value.Addr()
	}
	switch t := ptrvalue.Interface().(type) {
	case Unmarshaler:
		return t.UnmarshalBind(src)
	case Setter:
		return t.Set(src)
	}

	if reflect.TypeOf(src).AssignableTo(value.Type()) {
		value.Set(reflect.ValueOf(src))
		return
	}

	switch kind {
	case reflect.Bool:
		err = b.bindBool(value, src)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		err = b.bindInt(value, src)
	case reflect.Int64:
		err = b.bindInt64(value, src)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		err = b.bindUint(value, src)
	case reflect.Float32, reflect.Float64:
		err = b.bindFloat(value, src)
	case reflect.String:
		err = b.bindString(value, src)
	case reflect.Pointer:
		err = b.bindPointer(value, src)
	case reflect.Interface:
		err = b.bindInterface(value, src)
	case reflect.Struct:
		err = b.bindStruct(value, src)
	case reflect.Array:
		err = b.bindArray(value, src)
	case reflect.Slice:
		err = b.bindSlice(value, src)
	case reflect.Map:
		err = b.bindMap(value, src)

	// case reflect.Chan:
	// case reflect.Func:
	// case reflect.Complex64:
	// case reflect.Complex128:
	// case reflect.UnsafePointer:
	default:
		err = fmt.Errorf("unsupport to bind %T to a value", value.Interface())
	}

	return
}

func (b binder) bindBool(dstValue reflect.Value, src any) (err error) {
	v, err := cast.ToBool(src)
	if err == nil {
		dstValue.SetBool(v)
	}
	return
}

func (b binder) bindInt(dstValue reflect.Value, src any) (err error) {
	v, err := cast.ToInt64(src)
	if err == nil {
		dstValue.SetInt(v)
	}
	return
}

func (b binder) bindInt64(dstValue reflect.Value, src any) (err error) {
	if _, ok := dstValue.Interface().(time.Duration); !ok {
		return b.bindInt(dstValue, src)
	}

	v, err := cast.ToDuration(src)
	if err == nil {
		dstValue.SetInt(int64(v))
	}
	return
}

func (b binder) bindUint(dstValue reflect.Value, src any) (err error) {
	v, err := cast.ToUint64(src)
	if err == nil {
		dstValue.SetUint(v)
	}
	return
}

func (b binder) bindFloat(dstValue reflect.Value, src any) (err error) {
	v, err := cast.ToFloat64(src)
	if err == nil {
		dstValue.SetFloat(v)
	}
	return
}

func (b binder) bindString(dstValue reflect.Value, src any) (err error) {
	v, err := cast.ToString(src)
	if err == nil {
		dstValue.SetString(v)
	}
	return
}

func (b binder) bindPointer(dstValue reflect.Value, src any) (err error) {
	if dstValue.IsNil() {
		dstValue.Set(reflect.New(dstValue.Type().Elem()))
	}
	dstValue = dstValue.Elem()
	return b.bind(dstValue.Kind(), dstValue, src)
}

func (b binder) bindInterface(dstValue reflect.Value, src any) (err error) {
	if dstValue.IsValid() && dstValue.Elem().IsValid() { // Interface is set to a specific value.
		elem := dstValue.Elem()
		bindElem := elem

		// If we can't address this element, then its not writable. Instead,
		// we make a copy of the value (which is a pointer and therefore
		// writable), decode into that, and replace the whole value.
		var copied bool
		if !elem.CanAddr() {
			if elem.Kind() == reflect.Pointer {
				// (xgf) If it is a pointer and the element is addressable,
				// we should not new one, and still use the old.
				if elem.Elem().CanAddr() {
					// (xgf) We use the old pointer to check
					// whether it has implemented the interface
					// Unmarshaler or Setter.
					bindElem = elem
				} else {
					copied = true
				}
			}
		}
		if copied {
			bindElem = reflect.New(elem.Type()) // v = new(T)
			bindElem.Elem().Set(elem)           // *v = elem
		}

		err = b.bind(bindElem.Kind(), elem, src)
		if err != nil || !copied {
			return
		}

		dstValue.Set(elem.Elem()) // elem is copied.
		return
	}

	srcValue := reflect.ValueOf(src)
	dstType := dstValue.Type()

	// If the input data is a pointer, and the assigned type is the dereference
	// of that exact pointer, then indirect it so that we can assign it.
	// Example: *string to string
	if srcValue.Kind() == reflect.Pointer && srcValue.Type().Elem() == dstType {
		srcValue = reflect.Indirect(srcValue)
	}

	if !srcValue.IsValid() {
		srcValue = reflect.Zero(dstType)
	}

	srcType := srcValue.Type()
	if !srcType.AssignableTo(dstType) {
		return fmt.Errorf("cannot assign %s to %s", srcType.String(), dstType.String())
	}

	dstValue.Set(srcValue)
	return
}

func (b binder) bindArray(dstValue reflect.Value, src any) (err error) {
	return b._bindList(dstValue, src, true)
}

func (b binder) bindSlice(dstValue reflect.Value, src any) (err error) {
	return b._bindList(dstValue, src, false)
}

func (b binder) _bindList(dstValue reflect.Value, src any, isArray bool) (err error) {
	dstType := dstValue.Type()
	ekind := dstType.Elem().Kind()

	var _len int
	var bind func(reflect.Value, int) error
	switch vs := src.(type) {
	case []any:
		_len = len(vs)
		bind = func(v reflect.Value, i int) error { return b.bind(ekind, v, vs[i]) }

	case []string:
		_len = len(vs)
		bind = func(v reflect.Value, i int) error { return b.bind(ekind, v, vs[i]) }

	default:
		srcValue := reflect.ValueOf(src)
		switch srcValue.Kind() {
		case reflect.Array, reflect.Slice:
			_len = srcValue.Len()
			bind = func(v reflect.Value, i int) error {
				return b.bind(ekind, v, srcValue.Index(i).Interface())
			}
		default:

			return errors.New("cannot bind a slice type to a non-array/slice type")
		}
	}

	elems := dstValue
	if isArray {
		dstlen := dstValue.Len()
		if dstlen == 0 {
			return
		}
		if _len < dstlen {
			_len = dstlen
		}
	} else {
		elems = reflect.MakeSlice(dstType, _len, _len)
	}

	for i := 0; i < _len; i++ {
		if err = bind(elems.Index(i), i); err != nil {
			return
		}
	}

	if !isArray {
		dstValue.Set(elems)
	}
	return
}

func (b binder) bindMap(dstValue reflect.Value, src any) (err error) {
	dstType := dstValue.Type()
	keyType := dstType.Key()
	valueType := dstType.Elem()

	var dstmaps reflect.Value
	switch srcmaps := src.(type) {
	case map[string]any:
		dstmaps = reflect.MakeMapWithSize(dstType, len(srcmaps))
		for key, value := range srcmaps {
			err = b._bindMapIndex(dstmaps, keyType, valueType, key, value)
			if err != nil {
				return
			}
		}

	case map[string]string:
		dstmaps = reflect.MakeMapWithSize(dstType, len(srcmaps))
		for key, value := range srcmaps {
			err = b._bindMapIndex(dstmaps, keyType, valueType, key, value)
			if err != nil {
				return
			}
		}

	default:
		srcValue := reflect.ValueOf(src)
		if srcValue.Kind() != reflect.Map {
			return errors.New("cannot bind a map type to a non-map type")
		}

		dstmaps = reflect.MakeMapWithSize(dstType, srcValue.Len())
		for iter := srcValue.MapRange(); iter.Next(); {
			key, value := iter.Key().Interface(), iter.Value().Interface()
			err = b._bindMapIndex(dstmaps, keyType, valueType, key, value)
			if err != nil {
				return
			}
		}
	}

	dstValue.Set(dstmaps)
	return
}

func (b binder) _bindMapIndex(dstmap reflect.Value, keyType, valueType reflect.Type, key, value any) (err error) {
	srckey := reflect.New(keyType)
	err = b.bind(keyType.Kind(), srckey.Elem(), key)
	if err != nil {
		return
	}

	dstvalue := reflect.New(valueType)
	err = b.bind(valueType.Kind(), dstvalue.Elem(), value)
	if err != nil {
		return
	}

	dstmap.SetMapIndex(srckey.Elem(), dstvalue.Elem())
	return
}

func (b binder) bindStruct(dstStructValue reflect.Value, src any) (err error) {
	if _, ok := dstStructValue.Interface().(time.Time); ok {
		var v time.Time
		if v, err = cast.ToTime(src); err == nil {
			dstStructValue.Set(reflect.ValueOf(v))
		}
		return
	}

	fields := field.GetAllFields(dstStructValue.Type())
	for index, field := range fields {
		err = b.bindField(dstStructValue.Field(index), field, src)
		if err != nil {
			return
		}
	}
	return
}

func (b binder) bindField(fieldValue reflect.Value, fieldType reflect.StructField, src any) (err error) {
	if !fieldValue.CanSet() {
		return
	}

	name, arg := b.getFieldName(fieldType)
	if name == "" {
		return
	}

	fieldKind := fieldValue.Kind()
	if fieldKind == reflect.Struct && (fieldType.Anonymous || arg == "squash") {
		return b.bindStruct(fieldValue, src)
	}

	srcValue := reflect.ValueOf(src)
	if srcValue.Kind() != reflect.Map {
		return fmt.Errorf("unsupport to bind a struct to %T", src)
	} else if srcValue.Len() == 0 {
		return
	}

	if value := srcValue.MapIndex(reflect.ValueOf(name)); value.IsValid() {
		err = b.bind(fieldKind, fieldValue, value.Interface())
	}

	return
}
