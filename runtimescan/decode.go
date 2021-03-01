package runtimescan

import (
	"reflect"
)

// Decode convert from some source into struct by using tag information.
func Decode(dest interface{}, tag string, decoder Decoder) error {
	v, err := getParser(dest, tag, decoder)
	if err != nil {
		return err
	}
	return decode(dest, v, decoder)
}

func decode(dest interface{}, v *parser, decoder Decoder) error {
	current := reflect.ValueOf(dest).Elem()
	stack := []reflect.Value{current}
	var errors []error
	for i, op := range v.fieldOps {
		index := v.fieldIndexes[i]
		field := v.fields[i]
		switch op {
		case visitFieldOp:
			fv := current.Field(index)
			value, err := decoder.ExtractValue(field.tag)
			if err == Skip {
				continue
			} else if err != nil {
				errors = append(errors, err)
				continue
			}
			err = fuzzyAssign(fv, field.eType, field.eKind, field.isPtr, value)
		case visitChildOp:
			current := current.Field(index)
			stack = append(stack, current)
		case leaveChildOp:
			stack = stack[:len(stack)-1]
			current = stack[len(stack)-1]
		}
	}
	if len(errors) > 0 {
		return &Errors{
			Errors: errors,
		}
	}
	return nil
}
