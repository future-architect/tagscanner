package runtimescan

import (
	"reflect"
)

// Encode convert from some source into struct by using tag information.
func Encode(src interface{}, tags []string, encoder Encoder) error {
	v, err := getParser(src, tags, encoder)
	if err != nil {
		return err
	}
	return encode(encoder, v, src)
}

func encode(encoder Encoder, v *parser, src interface{}) error {
	current := reflect.ValueOf(src).Elem()
	stack := []reflect.Value{current}
	var errors []error
	for i, op := range v.fieldOps {
		index := v.fieldIndexes[i]
		field := v.fields[i]
		switch op {
		case visitFieldOp:
			fv := current.Field(index)
			var value interface{}
			if field.isPtr {
				if fv.IsNil() {
					value = nil
				} else {
					value = fv.Elem().Interface()
				}
			} else {
				value = fv.Interface()
			}
			err := encoder.VisitField(field.tag, value)
			if err == Skip {
				continue
			} else if err != nil {
				errors = append(errors, err)
				continue
			}
		case visitChildOp:
			// todo: call EnterChild
			current = current.Field(index)
			stack = append(stack, current)
		case leaveChildOp:
			// todo: call LeaveChild
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
