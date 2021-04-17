package runtimescan

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

type visitOpType int

const (
	visitFieldOp visitOpType = iota + 1
	visitChildOp
	leaveChildOp
)

type field struct {
	tag   interface{}
	eKind reflect.Kind
	eType reflect.Type
	isPtr bool
}

type parser struct {
	errors           []error
	fields           []*field
	fieldIndexes     []int
	fieldOps         []visitOpType
	panicWhenParsing bool
}

func newParser(vi Parser, tags []string, s interface{}) (*parser, error) {
	err := shouldPointerOfStruct(s)
	if err != nil {
		return nil, err
	}

	visitor := &parser{}
	visitor.parse(vi, tags, reflect.ValueOf(s).Type().Elem())
	if len(visitor.errors) == 0 {
		return visitor, nil
	}
	return nil, &Errors{Errors: visitor.errors}
}

func shouldPointerOfStruct(s interface{}) error {
	v := reflect.ValueOf(s)
	if v.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("sample should be pointer of struct: %w", ErrParseTag)
	}
	e := v.Elem()
	if e.Type().Kind() == reflect.Ptr {
		return fmt.Errorf("param s should be pointer of struct, but it is pointer of pointer: %w", ErrParseTag)
	}
	return nil
}

func (d *parser) parse(vi Parser, tags []string, t reflect.Type) error {
	d.errors = nil
	d.fields = nil
	d.fieldIndexes = nil
	d.fieldOps = nil
	d.parseTags(vi, tags, t, nil)
	return nil
}

func isPublic(f reflect.StructField) bool {
	if f.Anonymous {
		return true
	}
	var first rune
	for _, r := range f.Name {
		first = r
		break
	}
	return unicode.IsUpper(first)
}

func (d *parser) parseTags(vi Parser, tags []string, t reflect.Type, path []string) {
	for i := 0; i < t.NumField(); i++ {
		index := i
		f := t.Field(i)
		if !isPublic(f) {
			continue
		}
		var name string
		var hasChild bool
		if f.Anonymous {
			name = "(embed)"
			hasChild = true
		} else {
			name = f.Name
			if f.Type.Kind() == reflect.Struct {
				hasChild = true
			}
		}

		currentPath := append(path, name)
		pathStr := strings.Join(currentPath, ".")
		isPtr := t.Field(i).Type.Kind() == reflect.Ptr
		var eKind reflect.Kind
		var eType reflect.Type
		if isPtr {
			eKind = t.Field(i).Type.Elem().Kind()
			eType = t.Field(i).Type.Elem()
		} else {
			eKind = t.Field(i).Type.Kind()
			eType = t.Field(i).Type
		}
		var tag string
		for _, t := range tags {
			tag = f.Tag.Get(t)
			if tag != "" {
				continue
			}
		}
		t, err := vi.ParseTag(f.Name, tag, pathStr, eType)
		var skipTraverse bool
		var skipAdd bool
		if err == Skip {
			skipAdd = true
		} else if err == SkipTraverse {
			skipTraverse = true
		} else if err != nil {
			d.errors = append(d.errors, err)
			continue
		}
		if hasChild && !skipTraverse {
			d.fieldIndexes = append(d.fieldIndexes, index)
			d.fieldOps = append(d.fieldOps, visitChildOp)
			if !skipAdd {
				d.fields = append(d.fields, &field{
					tag:   t,
					eType: eType,
					eKind: eKind,
					isPtr: isPtr,
				})
			} else {
				d.fields = append(d.fields, nil)
			}
			d.parseTags(vi, tags, f.Type, path)
			d.fieldIndexes = append(d.fieldIndexes, -1)
			d.fieldOps = append(d.fieldOps, leaveChildOp)
			d.fields = append(d.fields, nil)
		} else if !skipAdd {
			d.fieldIndexes = append(d.fieldIndexes, index)
			d.fieldOps = append(d.fieldOps, visitFieldOp)
			d.fields = append(d.fields, &field{
				tag:   t,
				eType: eType,
				eKind: eKind,
				isPtr: isPtr,
			})
		}
	}
}

type parserCacheKey struct {
	Type   reflect.Type
	Parser reflect.Type
	Tag    string
}

var parsers = make(map[parserCacheKey]*parser)

func getParser(dest interface{}, tags []string, parser Parser) (*parser, error) {
	err := shouldPointerOfStruct(dest)
	if err != nil {
		return nil, err
	}
	key := parserCacheKey{
		Type:   reflect.ValueOf(dest).Type(),
		Parser: reflect.ValueOf(parser).Elem().Type(),
		Tag:    strings.Join(tags, ":"),
	}
	v, ok := parsers[key]
	if !ok {
		v, err = newParser(parser, tags, dest)
		if err != nil {
			return nil, err
		}
		parsers[key] = v
	}
	return v, nil
}
