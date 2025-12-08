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
	case int, int64:
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
	res, err := parseJSON(json, 0)
	if err != nil {
		return &Value{Err: err}
	}
	return &Value{Kind: getKind(res), data: res}
}

func ParseWithLimit(json string, maxDepth int) *Value {
	res, err := parseJSON(json, maxDepth)
	if err != nil {
		return &Value{Err: err}
	}
	return &Value{Kind: getKind(res), data: res}
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
	Kind Kind
	data any
	Err  error
}

func (v *Value) Index(i int) *Value {
	if v.Err != nil {
		return &Value{Err: v.Err}
	}

	if v.Kind != Array {
		return &Value{Err: fmt.Errorf("not array")}
	}

	arrSlice, ok := v.data.([]any)
	if ok {
		if i < 0 || i >= len(arrSlice) {
			return &Value{Err: fmt.Errorf("index out of range")}
		}
		val := arrSlice[i]
		return &Value{Kind: getKind(val), data: val}
	} else {
		return &Value{Err: fmt.Errorf("invalid array structure")}
	}
}

func (v *Value) Slice(start int, end int) *Value {
	if v.Err != nil {
		return &Value{Err: v.Err}
	}

	if v.Kind != Array {
		return &Value{Err: fmt.Errorf("not array")}
	}

	arr, ok := v.data.([]any)
	if ok {
		if start >= len(arr) || end > len(arr) {
			return &Value{Err: fmt.Errorf("slice index out of range")}
		} else if start > end {
			return &Value{Err: fmt.Errorf("invalid index")}
		}

		slice := arr[start:end]
		return &Value{Kind: Array, data: slice}
	} else {
		return &Value{Err: fmt.Errorf("invalid array structure")}
	}
}

func (v *Value) Get(key string) *Value {
	if v.Err != nil {
		return &Value{Err: v.Err}
	}

	if v.Kind != Object {
		return &Value{Err: fmt.Errorf("not object")}
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
			return &Value{Kind: getKind(val), data: val}
		} else {
			return &Value{Err: fmt.Errorf("object doesnt have that key")}
		}
	} else {
		return &Value{Err: fmt.Errorf("invalid map structure")}
	}
}

func (v *Value) Int64() (int64, error) {
	if v.Err != nil {
		return 0, v.Err
	}

	if v.Kind != Int {
		return 0, fmt.Errorf("not int")
	}

	res, ok := v.data.(int64)
	if !ok {
		return 0, fmt.Errorf("invalid number")
	}
	return res, nil
}

func (v *Value) Int32() (int32, error) {
	if v.Err != nil {
		return 0, v.Err
	}

	if v.Kind != Int {
		return 0, fmt.Errorf("not int")
	}

	res, ok := v.data.(int64)
	if !ok {
		return 0, fmt.Errorf("not int")
	}
	if res > math.MaxInt32 || res < math.MinInt32 {
		return 0, fmt.Errorf("out of range")
	}
	return int32(res), nil
}

func (v *Value) Int16() (int16, error) {
	if v.Err != nil {
		return 0, v.Err
	}

	if v.Kind != Int {
		return 0, fmt.Errorf("not int")
	}

	res, ok := v.data.(int64)
	if !ok {
		return 0, fmt.Errorf("not int")
	}
	if res > math.MaxInt16 || res < math.MinInt16 {
		return 0, fmt.Errorf("out of range")
	}
	return int16(res), nil
}

func (v *Value) Int8() (int8, error) {
	if v.Err != nil {
		return 0, v.Err
	}

	if v.Kind != Int {
		return 0, fmt.Errorf("not int")
	}

	res, ok := v.data.(int64)
	if !ok {
		return 0, fmt.Errorf("not int")
	}
	if res > math.MaxInt8 || res < math.MinInt8 {
		return 0, fmt.Errorf("out of range")
	}
	return int8(res), nil
}

func (v *Value) Int() (int, error) {
	if v.Err != nil {
		return 0, v.Err
	}

	if v.Kind != Int {
		return 0, fmt.Errorf("not int")
	}

	res, ok := v.data.(int64)
	if !ok {
		return 0, fmt.Errorf("not int")
	}
	if res > math.MaxInt || res < math.MinInt {
		return 0, fmt.Errorf("out of range")
	}
	return int(res), nil
}

func (v *Value) Uint64() (uint64, error) {
	if v.Err != nil {
		return 0, v.Err
	}

	if v.Kind != Int {
		return 0, fmt.Errorf("not int")
	}

	res, ok := v.data.(int64)
	if !ok {
		return 0, fmt.Errorf("not int")
	}
	if res < 0 {
		return 0, fmt.Errorf("out of range")
	}
	return uint64(res), nil
}

func (v *Value) Uint32() (uint32, error) {
	if v.Err != nil {
		return 0, v.Err
	}

	if v.Kind != Int {
		return 0, fmt.Errorf("not int")
	}

	res, ok := v.data.(int64)
	if !ok {
		return 0, fmt.Errorf("not int")
	}
	if res < 0 || res > math.MaxUint32 {
		return 0, fmt.Errorf("out of range")
	}
	return uint32(res), nil
}

func (v *Value) Uint16() (uint16, error) {
	if v.Err != nil {
		return 0, v.Err
	}

	if v.Kind != Int {
		return 0, fmt.Errorf("not int")
	}

	res, ok := v.data.(int64)
	if !ok {
		return 0, fmt.Errorf("not int")
	}
	if res < 0 || res > math.MaxUint16 {
		return 0, fmt.Errorf("out of range")
	}
	return uint16(res), nil
}

func (v *Value) Uint8() (uint8, error) {
	if v.Err != nil {
		return 0, v.Err
	}

	if v.Kind != Int {
		return 0, fmt.Errorf("not int")
	}

	res, ok := v.data.(int64)
	if !ok {
		return 0, fmt.Errorf("not int")
	}
	if res < 0 || res > math.MaxUint8 {
		return 0, fmt.Errorf("out of range")
	}
	return uint8(res), nil
}

func (v *Value) Uint() (uint, error) {
	if v.Err != nil {
		return 0, v.Err
	}

	if v.Kind != Int {
		return 0, fmt.Errorf("not int")
	}

	res, ok := v.data.(int64)
	if !ok {
		return 0, fmt.Errorf("not int")
	}
	if res < 0 || uint(res) > math.MaxUint {
		return 0, fmt.Errorf("out of range")
	}
	return uint(res), nil
}

func (v *Value) Float64() (float64, error) {
	if v.Err != nil {
		return 0, v.Err
	}

	if v.Kind != Float {
		return 0, fmt.Errorf("not float")
	}

	res, ok := v.data.(float64)
	if !ok {
		return 0, fmt.Errorf("not float")
	}
	return res, nil
}

func (v *Value) Float32() (float32, error) {
	if v.Err != nil {
		return 0, v.Err
	}

	if v.Kind != Float {
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
	if v.Err != nil {
		return "", v.Err
	}

	if v.Kind != String {
		return "", fmt.Errorf("not string")
	}

	res, ok := v.data.(string)
	if !ok {
		return "", fmt.Errorf("not string")
	}
	return res, nil
}

func (v *Value) Rune() (rune, error) {
	if v.Err != nil {
		return 0, v.Err
	}

	if v.Kind != String {
		return 0, fmt.Errorf("not string")
	}

	res, ok := v.data.(string)
	if !ok {
		return 0, fmt.Errorf("not string")
	}
	if len(res) == 0 {
		return 0, fmt.Errorf("empty string")
	}
	if len(res) > 1 {
		return 0, fmt.Errorf("string length is greater than 1")
	}
	if !utf8.ValidString(res) {
		return 0, fmt.Errorf("not valid string")
	}
	decodedRune, _ := utf8.DecodeRuneInString(res)
	return decodedRune, nil
}

func (v *Value) Bool() (bool, error) {
	if v.Err != nil {
		return false, v.Err
	}

	if v.Kind != Bool {
		return false, fmt.Errorf("not boolean")
	}

	res, ok := v.data.(bool)
	if !ok {
		return false, fmt.Errorf("not boolean")
	}
	return res, nil
}

func (v *Value) Array() ([]any, error) {
	if v.Err != nil {
		return nil, v.Err
	}

	if v.Kind != Array {
		return nil, fmt.Errorf("not array")
	}
	arr, ok := v.data.([]any)
	if !ok {
		return nil, fmt.Errorf("not array")
	}
	return arr, nil
}

func (v *Value) Object() (map[string]any, error) {
	if v.Err != nil {
		return nil, v.Err
	}

	if v.Kind != Object {
		return nil, fmt.Errorf("not object")
	}
	obj, ok := v.data.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("not object")
	}
	return obj, nil
}

func (v *Value) Null() (any, error) {
	if v.Err != nil {
		return nil, v.Err
	}

	if v.Kind != Null {
		return nil, fmt.Errorf("not null")
	}
	return nil, nil
}

func (v *Value) Any() (any, error) {
	if v.Err != nil {
		return nil, v.Err
	}
	return v.data, nil
}

func (v *Value) Error() error {
	return v.Err
}
