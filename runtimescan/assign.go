package runtimescan

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// Epsilon is a threshold value that specify input floating point value is true or false
// when converting from float64 to bool by using FuzzyAssign().
var Epsilon = 0.001

// ErrAssignError is a base error that is happens in FuzzyAssign().
var ErrAssignError = errors.New("assign error")

var supportedKind = map[reflect.Kind]bool{
	reflect.Int:       true,
	reflect.Uint:      true,
	reflect.Float64:   true,
	reflect.Bool:      true,
	reflect.String:    true,
	reflect.Struct:    true,
	reflect.Interface: true,
}

// FuzzyAssign assigns value to variable. It converts data format to meet variable type as much as possible.
func FuzzyAssign(dest, value interface{}) error {
	dv := reflect.ValueOf(dest)
	if dv.Kind() != reflect.Ptr {
		return fmt.Errorf("dest should be pointer: %w", ErrAssignError)
	}
	de := dv.Elem()
	isPtr := false
	eKind := de.Kind()
	eType := de.Type()
	if de.Type().Kind() == reflect.Ptr {
		isPtr = true
		eType = de.Type().Elem()
		eKind = de.Type().Elem().Kind()
		de = de.Elem()
	}
	if de.Kind() == reflect.Ptr {
		return fmt.Errorf("dest should be pointer or pointer of pointer: %w", ErrAssignError)
	}
	if !supportedKind[eKind] {
		return fmt.Errorf("dest type(%s) is not supported: %w", de.String(), ErrAssignError)
	}
	return fuzzyAssign(dv.Elem(), eType, eKind, isPtr, value)
}

func fuzzyAssign(dest reflect.Value, eType reflect.Type, eKind reflect.Kind, isPtr bool, value interface{}) error {
	if isPtr {
		switch eKind {
		case reflect.String:
			switch v := value.(type) {
			case int:
				s := strconv.Itoa(v)
				dest.Set(reflect.ValueOf(&s))
			case *int:
				s := strconv.Itoa(*v)
				dest.Set(reflect.ValueOf(&s))
			case uint:
				s := strconv.FormatUint(uint64(v), 10)
				dest.Set(reflect.ValueOf(&s))
			case *uint:
				s := strconv.FormatUint(uint64(*v), 10)
				dest.Set(reflect.ValueOf(&s))
			case float64:
				s := strconv.FormatFloat(v, 'G', -1, 64)
				dest.Set(reflect.ValueOf(&s))
			case *float64:
				s := strconv.FormatFloat(*v, 'G', -1, 64)
				dest.Set(reflect.ValueOf(&s))
			case bool:
				value := "false"
				if v {
					value = "true"
				}
				dest.Set(reflect.ValueOf(&value))
			case *bool:
				value := "false"
				if *v {
					value = "true"
				}
				dest.Set(reflect.ValueOf(&value))
			case string:
				dest.Set(reflect.ValueOf(&v))
			case *string:
				dest.Set(reflect.ValueOf(v))
			default:
				panic("please send PR to support this type (1)")
			}
		case reflect.Int:
			var intValue int
			switch v := value.(type) {
			case int:
				intValue = v
			case *int:
				intValue = *v
			case uint:
				intValue = int(v)
			case *uint:
				intValue = int(*v)
			case float64:
				intValue = int(math.Round(v))
			case *float64:
				intValue = int(math.Round(*v))
			case bool:
				if v {
					intValue = 1
				}
			case *bool:
				if *v {
					intValue = 1
				}
			case string:
				value64, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return err
				}
				intValue = int(value64)
			case *string:
				value64, err := strconv.ParseInt(*v, 10, 64)
				if err != nil {
					return err
				}
				intValue = int(value64)
			default:
				panic("please send PR to support this type (2)")
			}
			dest.Set(reflect.ValueOf(&intValue))
		case reflect.Uint:
			var uintValue uint
			switch v := value.(type) {
			case int:
				uintValue = uint(v)
			case *int:
				uintValue = uint(*v)
			case uint:
				uintValue = v
			case *uint:
				uintValue = *v
			case float64:
				uintValue = uint(math.Round(v))
			case *float64:
				uintValue = uint(math.Round(*v))
			case bool:
				if v {
					uintValue = 1.0
				}
			case *bool:
				if *v {
					uintValue = 1.0
				}
			case string:
				value64, err := strconv.ParseUint(v, 10, 64)
				if err != nil {
					return err
				}
				uintValue = uint(value64)
			case *string:
				value64, err := strconv.ParseUint(*v, 10, 64)
				if err != nil {
					return err
				}
				uintValue = uint(value64)
			default:
				panic("please send PR to support this type (3)")
			}
			dest.Set(reflect.ValueOf(&uintValue))
		case reflect.Float64:
			switch v := value.(type) {
			case int:
				value := float64(v)
				dest.Set(reflect.ValueOf(&value))
			case *int:
				value := float64(*v)
				dest.Set(reflect.ValueOf(&value))
			case uint:
				value := float64(v)
				dest.Set(reflect.ValueOf(&value))
			case *uint:
				value := float64(*v)
				dest.Set(reflect.ValueOf(&value))
			case float32:
				value := float64(v)
				dest.Set(reflect.ValueOf(&value))
			case *float32:
				value := float64(*v)
				dest.Set(reflect.ValueOf(&value))
			case float64:
				value := v
				dest.Set(reflect.ValueOf(&value))
			case *float64:
				value := *v
				dest.Set(reflect.ValueOf(&value))
			case bool:
				var value float64
				if v {
					value = 1.0
				}
				dest.Set(reflect.ValueOf(&value))
			case *bool:
				var value float64
				if *v {
					value = 1.0
				}
				dest.Set(reflect.ValueOf(&value))
			case string:
				value, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return err
				}
				dest.Set(reflect.ValueOf(&value))
			case *string:
				value, err := strconv.ParseFloat(*v, 64)
				if err != nil {
					return err
				}
				dest.Set(reflect.ValueOf(&value))
			default:
				panic("please send PR to support this type (4)")
			}
		case reflect.Bool:
			switch v := value.(type) {
			case int:
				value := v != 0
				dest.Set(reflect.ValueOf(&value))
			case *int:
				value := *v != 0
				dest.Set(reflect.ValueOf(&value))
			case uint:
				value := v != 0
				dest.Set(reflect.ValueOf(&value))
			case *uint:
				value := *v != 0
				dest.Set(reflect.ValueOf(&value))
			case float64:
				value := v < -Epsilon || Epsilon < v
				dest.Set(reflect.ValueOf(&value))
			case *float64:
				value := *v < -Epsilon || Epsilon < *v
				dest.Set(reflect.ValueOf(&value))
			case bool:
				value := v
				dest.Set(reflect.ValueOf(&value))
			case *bool:
				value := *v
				dest.Set(reflect.ValueOf(&value))
			case string:
				lv := strings.ToLower(v)
				value := lv != "false" && lv != "no" && lv != ""
				dest.Set(reflect.ValueOf(&value))
			case *string:
				lv := strings.ToLower(*v)
				value := lv != "false" && lv != "no" && lv != ""
				dest.Set(reflect.ValueOf(&value))
			default:
				panic("please send PR to support this type (5)")
			}
		case reflect.Struct:
			v := reflect.ValueOf(value)
			if v.Kind() == reflect.Ptr {
				if v.Elem().Type() == eType {
					dest.Set(v)
				} else {
					return fmt.Errorf("struct type is different (dest=%s, src=%s): %w", eType.String(), v.Elem().Type().String(), ErrAssignError)
				}
			} else {
				if v.Type() == eType {
					cp := reflect.New(eType)
					cp.Elem().Set(v)
					dest.Set(cp)
				} else {
					return fmt.Errorf("struct type is different (dest=%s, src=%s): %w", eType.String(), v.Type().String(), ErrAssignError)
				}
			}
		default:
			panic("please send PR to support this type (6)")
		}
	} else {
		switch eKind {
		case reflect.String:
			switch v := value.(type) {
			case int:
				dest.SetString(strconv.Itoa(v))
			case *int:
				dest.SetString(strconv.Itoa(*v))
			case uint:
				dest.SetString(strconv.FormatUint(uint64(v), 10))
			case *uint:
				dest.SetString(strconv.FormatUint(uint64(*v), 10))
			case float64:
				dest.SetString(strconv.FormatFloat(v, 'G', -1, 64))
			case *float64:
				dest.SetString(strconv.FormatFloat(*v, 'G', -1, 64))
			case bool:
				value := "false"
				if v {
					value = "true"
				}
				dest.SetString(value)
			case *bool:
				value := "false"
				if *v {
					value = "true"
				}
				dest.SetString(value)
			case string:
				dest.SetString(v)
			case *string:
				dest.SetString(*v)
			default:
				panic("please send PR to support this type (7)")
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			switch v := value.(type) {
			case int:
				dest.SetInt(int64(v))
			case *int:
				dest.SetInt(int64(*v))
			case uint:
				dest.SetInt(int64(v))
			case *uint:
				dest.SetInt(int64(*v))
			case float64:
				dest.SetInt(int64(math.Round(v)))
			case *float64:
				dest.SetInt(int64(math.Round(*v)))
			case bool:
				var value int64
				if v {
					value = 1
				}
				dest.SetInt(value)
			case *bool:
				var value int64
				if *v {
					value = 1
				}
				dest.SetInt(value)
			case string:
				value, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return err
				}
				dest.SetInt(value)
			case *string:
				value, err := strconv.ParseInt(*v, 10, 64)
				if err != nil {
					return err
				}
				dest.SetInt(value)
			default:
				panic("please send PR to support this type (8)")
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			switch v := value.(type) {
			case int:
				dest.SetUint(uint64(v))
			case *int:
				dest.SetUint(uint64(*v))
			case uint:
				dest.SetUint(uint64(v))
			case *uint:
				dest.SetUint(uint64(*v))
			case float64:
				dest.SetUint(uint64(math.Round(v)))
			case *float64:
				dest.SetUint(uint64(math.Round(*v)))
			case bool:
				var value uint64
				if v {
					value = 1
				}
				dest.SetUint(value)
			case *bool:
				var value uint64
				if *v {
					value = 1
				}
				dest.SetUint(value)
			case string:
				value, err := strconv.ParseUint(v, 10, 64)
				if err != nil {
					return err
				}
				dest.SetUint(value)
			case *string:
				value, err := strconv.ParseUint(*v, 10, 64)
				if err != nil {
					return err
				}
				dest.SetUint(value)
			default:
				panic("please send PR to support this type (9)")
			}
		case reflect.Float64, reflect.Float32:
			switch v := value.(type) {
			case int:
				dest.SetFloat(float64(v))
			case *int:
				dest.SetFloat(float64(*v))
			case uint:
				dest.SetFloat(float64(v))
			case *uint:
				dest.SetFloat(float64(*v))
			case float32:
				dest.SetFloat(float64(v))
			case *float32:
				dest.SetFloat(float64(*v))
			case float64:
				dest.SetFloat(v)
			case *float64:
				dest.SetFloat(*v)
			case bool:
				var value float64
				if v {
					value = 1.0
				}
				dest.SetFloat(value)
			case *bool:
				var value float64
				if *v {
					value = 1.0
				}
				dest.SetFloat(value)
			case string:
				value, err := strconv.ParseFloat(v, 64)
				if err != nil {
					return err
				}
				dest.SetFloat(value)
			case *string:
				value, err := strconv.ParseFloat(*v, 64)
				if err != nil {
					return err
				}
				dest.SetFloat(value)
			default:
				panic("please send PR to support this type (10)")
			}
		case reflect.Bool:
			switch v := value.(type) {
			case int:
				value := v != 0
				dest.SetBool(value)
			case *int:
				value := *v != 0
				dest.SetBool(value)
			case uint:
				value := v != 0
				dest.SetBool(value)
			case *uint:
				value := *v != 0
				dest.SetBool(value)
			case float64:
				value := v < -Epsilon || Epsilon < v
				dest.SetBool(value)
			case *float64:
				value := *v < -Epsilon || Epsilon < *v
				dest.SetBool(value)
			case bool:
				dest.SetBool(v)
			case *bool:
				dest.SetBool(*v)
			case string:
				lv := strings.ToLower(v)
				value := lv != "false" && lv != "no" && lv != ""
				dest.SetBool(value)
			case *string:
				lv := strings.ToLower(*v)
				value := lv != "false" && lv != "no" && lv != ""
				dest.SetBool(value)
			default:
				panic("please send PR to support this type (12)")
			}
		case reflect.Interface:
			v := reflect.ValueOf(value)
			if v.Kind() == reflect.Ptr {
				if v.Elem().Type().AssignableTo(eType) || v.Type().AssignableTo(eType) {
					dest.Set(v)
				} else {
					return fmt.Errorf("struct is not assignable to interface (dest=%s, src=%s): %w", eType.String(), v.Elem().Type().String(), ErrAssignError)
				}
			} else {
				if v.Type().AssignableTo(eType) {
					cp := reflect.New(v.Type())
					cp.Elem().Set(v)
					dest.Set(cp)
				} else {
					return fmt.Errorf("struct is not assignable to interface (dest=%s, src=%s): %w", eType.String(), v.Type().String(), ErrAssignError)
				}
			}
		case reflect.Slice:
			v := reflect.ValueOf(value)
			if v.Kind() == reflect.Ptr {
				if v.Elem().Type().AssignableTo(eType) || v.Type().AssignableTo(eType) {
					dest.Set(v)
				} else {
					return fmt.Errorf("struct is not assignable to interface (dest=%s, src=%s): %w", eType.String(), v.Elem().Type().String(), ErrAssignError)
				}
			} else {
				if v.Type().AssignableTo(eType) {
					dest.Set(v)
				} else {
					return fmt.Errorf("struct is not assignable to interface (dest=%s, src=%s): %w", eType.String(), v.Type().String(), ErrAssignError)
				}
			}
		default:
			panic("please send PR to support this type (13)")
		}
	}
	return nil
}
