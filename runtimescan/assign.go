package runtimescan

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Epsilon is a threshold value that specify input floating point value is true or false
// when converting from float64 to bool by using FuzzyAssign().
var Epsilon = 0.001

// ErrAssignError is a base error that is happens in FuzzyAssign().
var ErrAssignError = errors.New("assign error")

func countPointerDepth(v reflect.Value) int {
	c := 0
	for v.Kind() == reflect.Pointer {
		v = v.Elem()
		c++
	}
	return c
}

func unwrap(v reflect.Value, ptr bool) (d reflect.Value, ok bool) {
	oldV := v
	loop := false
	for v.Kind() == reflect.Pointer {
		oldV = v
		v = v.Elem()
		loop = true
	}
	if ptr {
		if !loop {
			return reflect.Value{}, false
		}
		return oldV, true
	}
	return v, true
}

func directAssign(dest, value reflect.Value) bool {
	destType := dest.Type()
	vk := value.Kind()
	if dest.Kind() == reflect.Pointer { // Pointer
		destType = destType.Elem()
		if value.CanConvert(destType) && !(destType.Kind() == reflect.String && (vk == reflect.Int) || vk == reflect.Uint) {
			if !dest.Elem().CanSet() {
				dest.Set(reflect.New(dest.Type().Elem()))
			}
			dest.Elem().Set(value.Convert(destType))
			return true
		}
		return false
	} else { // Struct Field
		if value.CanConvert(destType) && !(destType.Kind() == reflect.String && (vk == reflect.Int) || vk == reflect.Uint) {
			dest.Set(value.Convert(destType))
			return true
		}
		return false
	}
}

// FuzzyAssign assigns value to variable. It converts data format to meet variable type as much as possible.
func FuzzyAssign(dest, value any) error {
	var dv reflect.Value
	if dv2, ok := dest.(reflect.Value); ok && dv2.CanSet() {
		// struct field's value is not pointer, but CanSet() is true
		dv = dv2
	} else {
		dv3, ok := unwrap(reflect.ValueOf(dest), true)
		if !ok {
			return fmt.Errorf("dest should be pointer, but %v: %w", reflect.TypeOf(dest), ErrAssignError)
		}
		dv = dv3
	}
	vv, _ := unwrap(reflect.ValueOf(value), false)
	if directAssign(dv, vv) {
		return nil
	}
	if dv.Type().Kind() == reflect.Interface {
		vt := reflect.TypeOf(value)
		if vt.AssignableTo(dv.Type()) {
			dv.Set(reflect.ValueOf(value))
			return nil
		}
		return fmt.Errorf("%v is not assignable to interface %v", vt, dv.Type())
	} else {
		if !dv.Elem().CanSet() { // If dv points to nil, create new instance
			dv.Set(reflect.New(dv.Type().Elem()))
		}
		return fuzzyAssign(dv.Elem(), dv.Type().Elem(), value)
	}
}

func fuzzyAssign(dest reflect.Value, eType reflect.Type, value any) error {
	eKind := eType.Kind()
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
		default:
			panic("please send PR to support this type (7)")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v := value.(type) {
		case uint:
			dest.SetInt(int64(v))
		case *uint:
			dest.SetInt(int64(*v))
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
		case uint:
			dest.SetUint(uint64(v))
		case *uint:
			dest.SetUint(uint64(*v))
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
		case uint:
			dest.SetFloat(float64(v))
		case *uint:
			dest.SetFloat(float64(*v))
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
	/*case reflect.Interface:
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
		}*/
	default:
		panic("please send PR to support this type (13)")
	}
	return nil
}
