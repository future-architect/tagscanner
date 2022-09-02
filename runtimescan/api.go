package runtimescan

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Parser is a part of Decoder and Encoder. The user of this library will implement this method.
type Parser interface {
	// ParseTag is called when parsing struct. This is called once per each struct definition
	// inside Decode() or Encode()
	// Returned value of this method is passed to ExtractValue() method call.
	//
	// name: Field name
	// tagKey: Tag key
	// tagStr: Tag string
	// pathStr: Field name for error message (it contains nested struct names)
	// elemType: Field type
	ParseTag(name, tagKey, tagStr, pathStr string, elemType reflect.Type) (tag any, err error)
}

// Decoder is an interface that extracts value from some type and assign to struct instance.
//
// ParseTag() is used when parsing struct tag.
type Decoder interface {
	Parser
	// ExtractValue is called for each field of each struct instance inside Decode().
	ExtractValue(tag any) (value any, err error)
}

// Encoder is an interface that receive values from some type and create other files.
//
// ParseTag() is used when parsing struct tag.
//
// VisitField() is called when getting value from source struct.
// EnterChild() and LeaveChild)() are called when traversing struct structure.
type Encoder interface {
	Parser
	VisitField(tag, value any) (err error)
	EnterChild(tag any) (err error)
	LeaveChild(tag any) (err error)
}

type Errors struct {
	Errors []error
}

func (e Errors) Error() string {
	var errors []string
	for _, e := range e.Errors {
		errors = append(errors, e.Error())
	}
	return fmt.Sprintf("%d errors: \n* %s", len(e.Errors), strings.Join(errors, "\n  "))
}

var (
	// SkipTraverse is returned by Parser interface's ParseTag() method to notify to skip traversing child struct.
	SkipTraverse = errors.New("skip traverse")
	ErrParseTag  = errors.New("tag parse error")
	// Skip is a flag to skip. This is returned by Parser interface's ParseTag() method to notify to add skip tag
	// and ExtractValue() method of Decoder interface.
	Skip = errors.New("skip")
)
