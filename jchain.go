package jchain

import (
	"fmt"
	"math"
	"unicode/utf8"
)

func getKind(val any) Kind {
	switch val.(type) {
	case map[string]any:
		return Object
	case []any:
		return Array
	case string:
		return String
	case int, int64, uint64:
		return Int
	case float64:
		return Float
	case bool:
		return Bool
	case nil:
		return Null
	default:
		return Invalid
	}
}

func Parse(json string) *Value {
	res, err := parseJSON(json, 1000)
	if err != nil {
		return &Value{err: err}
	}
	return &Value{kind: getKind(res), data: res}
}

func ParseUnlimited(json string) *Value {
	res, err := parseJSON(json, 0)
	if err != nil {
		return &Value{err: err}
	}
	return &Value{kind: getKind(res), data: res}
}

func ParseWithLimit(json string, maxDepth int) *Value {
	res, err := parseJSON(json, maxDepth)
	if err != nil {
		return &Value{err: err}
	}
	return &Value{kind: getKind(res), data: res}
}

type Kind int

const (
	Invalid Kind = iota
	Object
	Array
	String
	Int
	Float
	Bool
	Null
)

type Value struct {
	kind Kind
	data any
	err  error
}

func (v *Value) Index(i int) *Value {
	if v.err != nil {
		return &Value{err: v.err}
	}

	if v.kind != Array {
		return &Value{err: fmt.Errorf("not array")}
	}

	arrSlice, ok := v.data.([]any)
	if ok {
		if i < 0 || i >= len(arrSlice) {
			return &Value{err: fmt.Errorf("index out of range")}
		}
		val := arrSlice[i]
		return &Value{kind: getKind(val), data: val}
	} else {
		return &Value{err: fmt.Errorf("invalid array structure")}
	}
}

func (v *Value) Slice(start int, end int) *Value {
	if v.err != nil {
		return &Value{err: v.err}
	}

	if v.kind != Array {
		return &Value{err: fmt.Errorf("not array")}
	}

	arr, ok := v.data.([]any)
	if ok {
		if start >= len(arr) || end > len(arr) {
			return &Value{err: fmt.Errorf("slice index out of range")}
		} else if start > end {
			return &Value{err: fmt.Errorf("invalid index")}
		}

		slice := arr[start:end]
		return &Value{kind: Array, data: slice}
	} else {
		return &Value{err: fmt.Errorf("invalid array structure")}
	}
}

func (v *Value) Get(key string) *Value {
	if v.err != nil {
		return &Value{err: v.err}
	}

	if v.kind != Object {
		return &Value{err: fmt.Errorf("not object")}
	}

	obj, ok := v.data.(map[string]any)
	if ok {
		hasKey := false
		for k := range obj {
			if k == key {
				hasKey = true
				break
			}
		}
		if hasKey {
			val := obj[key]
			return &Value{kind: getKind(val), data: val}
		} else {
			return &Value{err: fmt.Errorf("object doesnt have that key")}
		}
	} else {
		return &Value{err: fmt.Errorf("invalid map structure")}
	}
}

func (v *Value) Int64() (int64, error) {
	if v.err != nil {
		return 0, v.err
	}

	if v.kind != Int {
		return 0, fmt.Errorf("not int")
	}

	switch val := v.data.(type) {
	case int64:
		return val, nil
	case uint64:
		if val > math.MaxInt64 {
			return 0, fmt.Errorf("out of range")
		}
		return int64(val), nil
	default:
		return 0, fmt.Errorf("invalid number")
	}
}

func (v *Value) Int32() (int32, error) {
	if v.err != nil {
		return 0, v.err
	}

	if v.kind != Int {
		return 0, fmt.Errorf("not int")
	}

	switch val := v.data.(type) {
	case int64:
		if val > math.MaxInt32 || val < math.MinInt32 {
			return 0, fmt.Errorf("out of range")
		}
		return int32(val), nil
	case uint64:
		if val > math.MaxInt32 {
			return 0, fmt.Errorf("out of range")
		}
		return int32(val), nil
	default:
		return 0, fmt.Errorf("not int")
	}
}

func (v *Value) Int16() (int16, error) {
	if v.err != nil {
		return 0, v.err
	}

	if v.kind != Int {
		return 0, fmt.Errorf("not int")
	}

	switch val := v.data.(type) {
	case int64:
		if val > math.MaxInt16 || val < math.MinInt16 {
			return 0, fmt.Errorf("out of range")
		}
		return int16(val), nil
	case uint64:
		if val > math.MaxInt16 {
			return 0, fmt.Errorf("out of range")
		}
		return int16(val), nil
	default:
		return 0, fmt.Errorf("not int")
	}
}

func (v *Value) Int8() (int8, error) {
	if v.err != nil {
		return 0, v.err
	}

	if v.kind != Int {
		return 0, fmt.Errorf("not int")
	}

	switch val := v.data.(type) {
	case int64:
		if val > math.MaxInt8 || val < math.MinInt8 {
			return 0, fmt.Errorf("out of range")
		}
		return int8(val), nil
	case uint64:
		if val > math.MaxInt8 {
			return 0, fmt.Errorf("out of range")
		}
		return int8(val), nil
	default:
		return 0, fmt.Errorf("not int")
	}
}

func (v *Value) Int() (int, error) {
	if v.err != nil {
		return 0, v.err
	}

	if v.kind != Int {
		return 0, fmt.Errorf("not int")
	}

	switch val := v.data.(type) {
	case int64:
		if val > math.MaxInt || val < math.MinInt {
			return 0, fmt.Errorf("out of range")
		}
		return int(val), nil
	case uint64:
		if val > math.MaxInt {
			return 0, fmt.Errorf("out of range")
		}
		return int(val), nil
	default:
		return 0, fmt.Errorf("not int")
	}
}

func (v *Value) Uint64() (uint64, error) {
	if v.err != nil {
		return 0, v.err
	}

	if v.kind != Int {
		return 0, fmt.Errorf("not int")
	}

	switch val := v.data.(type) {
	case int64:
		if val < 0 {
			return 0, fmt.Errorf("out of range")
		}
		return uint64(val), nil
	case uint64:
		return val, nil
	default:
		return 0, fmt.Errorf("not int")
	}
}

func (v *Value) Uint32() (uint32, error) {
	if v.err != nil {
		return 0, v.err
	}

	if v.kind != Int {
		return 0, fmt.Errorf("not int")
	}

	switch val := v.data.(type) {
	case int64:
		if val < 0 || val > math.MaxUint32 {
			return 0, fmt.Errorf("out of range")
		}
		return uint32(val), nil
	case uint64:
		if val > math.MaxUint32 {
			return 0, fmt.Errorf("out of range")
		}
		return uint32(val), nil
	default:
		return 0, fmt.Errorf("not int")
	}
}

func (v *Value) Uint16() (uint16, error) {
	if v.err != nil {
		return 0, v.err
	}

	if v.kind != Int {
		return 0, fmt.Errorf("not int")
	}

	switch val := v.data.(type) {
	case int64:
		if val < 0 || val > math.MaxUint16 {
			return 0, fmt.Errorf("out of range")
		}
		return uint16(val), nil
	case uint64:
		if val > math.MaxUint16 {
			return 0, fmt.Errorf("out of range")
		}
		return uint16(val), nil
	default:
		return 0, fmt.Errorf("not int")
	}
}

func (v *Value) Uint8() (uint8, error) {
	if v.err != nil {
		return 0, v.err
	}

	if v.kind != Int {
		return 0, fmt.Errorf("not int")
	}

	switch val := v.data.(type) {
	case int64:
		if val < 0 || val > math.MaxUint8 {
			return 0, fmt.Errorf("out of range")
		}
		return uint8(val), nil
	case uint64:
		if val > math.MaxUint8 {
			return 0, fmt.Errorf("out of range")
		}
		return uint8(val), nil
	default:
		return 0, fmt.Errorf("not int")
	}
}

func (v *Value) Uint() (uint, error) {
	if v.err != nil {
		return 0, v.err
	}

	if v.kind != Int {
		return 0, fmt.Errorf("not int")
	}

	switch val := v.data.(type) {
	case int64:
		if val < 0 || uint(val) > math.MaxUint {
			return 0, fmt.Errorf("out of range")
		}
		return uint(val), nil
	case uint64:
		if val > math.MaxUint {
			return 0, fmt.Errorf("out of range")
		}
		return uint(val), nil
	default:
		return 0, fmt.Errorf("not int")
	}
}

func (v *Value) Float64() (float64, error) {
	if v.err != nil {
		return 0, v.err
	}

	if v.kind != Float {
		return 0, fmt.Errorf("not float")
	}

	res, ok := v.data.(float64)
	if !ok {
		return 0, fmt.Errorf("not float")
	}
	return res, nil
}

func (v *Value) Float32() (float32, error) {
	if v.err != nil {
		return 0, v.err
	}

	if v.kind != Float {
		return 0, fmt.Errorf("not float")
	}

	res, ok := v.data.(float64)
	if !ok {
		return 0, fmt.Errorf("not float")
	}
	if res > math.MaxFloat32 || res < -math.MaxFloat32 {
		return 0, fmt.Errorf("out of range")
	}
	return float32(res), nil
}

func (v *Value) String() (string, error) {
	if v.err != nil {
		return "", v.err
	}

	if v.kind != String {
		return "", fmt.Errorf("not string")
	}

	res, ok := v.data.(string)
	if !ok {
		return "", fmt.Errorf("not string")
	}
	return res, nil
}

func (v *Value) Rune() (rune, error) {
	if v.err != nil {
		return 0, v.err
	}

	if v.kind != String {
		return 0, fmt.Errorf("not string")
	}

	res, ok := v.data.(string)
	if !ok {
		return 0, fmt.Errorf("not string")
	}
	if len(res) == 0 {
		return 0, fmt.Errorf("empty string")
	}
	if utf8.RuneCountInString(res) > 1 {
		return 0, fmt.Errorf("string length is greater than 1")
	}
	if !utf8.ValidString(res) {
		return 0, fmt.Errorf("not valid string")
	}
	decodedRune, _ := utf8.DecodeRuneInString(res)
	return decodedRune, nil
}

func (v *Value) Bool() (bool, error) {
	if v.err != nil {
		return false, v.err
	}

	if v.kind != Bool {
		return false, fmt.Errorf("not boolean")
	}

	res, ok := v.data.(bool)
	if !ok {
		return false, fmt.Errorf("not boolean")
	}
	return res, nil
}

func (v *Value) Array() ([]any, error) {
	if v.err != nil {
		return nil, v.err
	}

	if v.kind != Array {
		return nil, fmt.Errorf("not array")
	}
	arr, ok := v.data.([]any)
	if !ok {
		return nil, fmt.Errorf("not array")
	}
	return arr, nil
}

func (v *Value) Map() (map[string]any, error) {
	if v.err != nil {
		return nil, v.err
	}

	if v.kind != Object {
		return nil, fmt.Errorf("not object")
	}
	obj, ok := v.data.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("not map")
	}
	return obj, nil
}

func (v *Value) Nil() (any, error) {
	if v.err != nil {
		return nil, v.err
	}

	if v.kind != Null {
		return nil, fmt.Errorf("not null")
	}
	return nil, nil
}

func (v *Value) Any() (any, error) {
	if v.err != nil {
		return nil, v.err
	}
	return v.data, nil
}

func (v *Value) Error() error {
	if v.err != nil {
		return v.err
	}
	return nil
}

func (v *Value) Kind() Kind {
	if v.err != nil {
		return Invalid
	}
	return v.kind
}
