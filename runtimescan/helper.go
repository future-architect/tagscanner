package runtimescan

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// BasicTag is for convenience.
type BasicTag struct {
	Name     string
	Tag      string
	Path     string
	ElemType reflect.Type
}

// BasicParseTag is a helper function to make tagscanner easy.
//
// Both Encoder and Decoder should implement
func BasicParseTag(name, tagStr, pathStr string, elemType reflect.Type) (tag interface{}, err error) {
	if tagStr == "" {
		tagStr = strings.ToLower(name)
	}
	return &BasicTag{
		Name: name,
		Tag: tagStr,
		Path: pathStr,
		ElemType: elemType,
	}, nil
}

func Str2PrimitiveValue(value string, elemType reflect.Type) (interface{}, error) {
	switch elemType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}
		return v, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil, err
		}
		return v, nil
	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, err
		}
		return v, nil
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return nil, err
		}
		return v, nil
	case reflect.String:
		return value, nil
	default:
		return nil, fmt.Errorf("can't convert to %s", elemType.String())
	}
}

func IsPointerOfSliceOfStruct(i interface{}) bool {
	v := reflect.ValueOf(i)
	if v.Type().Kind() != reflect.Ptr {
		return false
	}
	if v.Type().Elem().Kind() != reflect.Slice {
		return false
	}
	return v.Type().Elem().Elem().Kind() == reflect.Struct
}

func IsPointerOfSliceOfPointerOfStruct(i interface{}) bool {
	v := reflect.ValueOf(i)
	if v.Type().Kind() != reflect.Ptr {
		return false
	}
	if v.Type().Elem().Kind() != reflect.Slice {
		return false
	}
	if v.Type().Elem().Elem().Kind() != reflect.Ptr {
		return false
	}
	return v.Type().Elem().Elem().Elem().Kind() == reflect.Struct
}

func IsPointerOfStruct(i interface{}) bool {
	v := reflect.ValueOf(i)
	if v.Type().Kind() != reflect.Ptr {
		return false
	}
	return v.Type().Elem().Kind() == reflect.Struct
}

func NewStructInstance(i interface{}) (interface{}, error) {
	if IsPointerOfStruct(i) {
		v := reflect.New(reflect.TypeOf(i).Elem())
		return v.Interface(), nil
	} else if IsPointerOfSliceOfStruct(i) {
		v := reflect.New(reflect.TypeOf(i).Elem().Elem())
		return v.Interface(), nil
	} else if IsPointerOfSliceOfPointerOfStruct(i) {
		v := reflect.New(reflect.TypeOf(i).Elem().Elem().Elem())
		return v.Interface(), nil
	}
	return nil, errors.New("input is not *Struct or *[]Struct")
}