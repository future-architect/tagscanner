package restmap

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"unicode"

	"github.com/future-architect/tagscanner/runtimescan"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type FieldType string

const (
	MethodField  FieldType = "method"
	PathField    FieldType = "path"
	HeaderField  FieldType = "header"
	CookieField  FieldType = "cookie"
	QueryField   FieldType = "query"
	BodyField    FieldType = "body"
	ContextField FieldType = "context"
	IgnoreField  FieldType = "-"
)

type defaultFieldType int

const (
	lowerCase defaultFieldType = iota + 1
	hyphenatedPascalCase
	hyphenatedLowerCase
	noField
)

func (t defaultFieldType) Convert(source string) string {
	switch t {
	case lowerCase:
		return strings.ToLower(source)
	case hyphenatedLowerCase:
		words := splitWord(source)
		var result []string
		for _, word := range words {
			result = append(result, strings.ToLower(word))
		}
		return strings.Join(result, "-")
	case hyphenatedPascalCase:
		words := splitWord(source)
		var result []string
		for _, word := range words {
			result = append(result, strings.Title(word))
		}
		return strings.Join(result, "-")
	case noField:
		return ""
	}
	return ""
}

func splitWord(source string) []string {
	var word []rune
	var result []string
	for _, c := range source {
		if unicode.IsUpper(c) {
			if len(word) != 0 {
				result = append(result, string(word))
			}
			word = []rune{c}
		} else {
			word = append(word, c)
		}
	}
	if len(word) > 0 {
		result = append(result, string(word))
	}
	return result
}

var fieldTypes = map[FieldType]struct {
	Type     defaultFieldType
	Optional bool
}{
	MethodField: {
		Type:     noField,
		Optional: false,
	},
	PathField: {
		Type:     hyphenatedLowerCase,
		Optional: false,
	},
	QueryField: {
		Type:     lowerCase,
		Optional: true,
	},
	HeaderField: {
		Type:     hyphenatedPascalCase,
		Optional: true,
	},
	CookieField: {
		Type:     hyphenatedLowerCase,
		Optional: true,
	},
	BodyField: {
		Type:     hyphenatedLowerCase,
		Optional: true,
	},
	ContextField: {
		Type:     noField,
		Optional: true,
	},
}

func nearest(tagName string) string {
	dmp := diffmatchpatch.New()
	var nearest string
	var dist = math.MaxInt64
	for fieldType := range fieldTypes {
		diffs := dmp.DiffMain(tagName, string(fieldType), false)
		lv := dmp.DiffLevenshtein(diffs)
		if lv < dist {
			dist = lv
			nearest = string(fieldType)
		}
	}
	return nearest
}

type RestTag struct {
	Type     FieldType
	EType    reflect.Type
	Name     string
	Default  string
	Optional bool
	Base     bool
}

func ParseRestTag(fieldName, tagSource, fullPath string, eType reflect.Type) (*RestTag, error) {
	sources := strings.Split(tagSource, ",")
	if len(sources) > 1 {
		tagSource = sources[0]
	}
	tags := strings.Split(tagSource, ":")
	t := RestTag{
		Type:  IgnoreField,
		EType: eType,
	}
	var hasDefault bool
	switch len(tags) {
	case 1:
		t.Type = FieldType(tags[0])
		if f, ok := fieldTypes[t.Type]; ok {
			hasDefault = f.Optional
			t.Name = f.Type.Convert(fieldName)
		} else {
			return nil, fmt.Errorf("tag name '%s' of fieldName '%s' is invalid. did you mean '%s'?: %w", tags[0], fullPath, nearest(tags[0]), runtimescan.ErrParseTag)
		}
	case 2:
		t.Type = FieldType(tags[0])
		if f, ok := fieldTypes[t.Type]; ok {
			hasDefault = f.Optional
			if f.Type != noField {
				t.Name = tags[1]
			} else {
				return nil, fmt.Errorf("tag type '%s' of fieldName '%s' can't have extra information '%s': %w", tags[0], fullPath, tags[1], runtimescan.ErrParseTag)
			}
		} else {
			return nil, fmt.Errorf("tag type '%s' of fieldName '%s' is invalid. did you mean '%s'?: %w", tags[0], fullPath, nearest(tags[0]), runtimescan.ErrParseTag)
		}
	default:
		return nil, fmt.Errorf("tag '%s' of fieldName '%s' is invalid. Zero or one colon ':' should be included: %w", tagSource, fullPath, runtimescan.ErrParseTag)
	}
	for i := 1; i < len(sources); i++ {
		fragment := sources[i]
		if strings.HasPrefix(fragment, "default:") {
			if hasDefault {
				t.Default = strings.TrimPrefix(fragment, "default:")
			} else {
				return nil, fmt.Errorf("tag type '%s' of fieldName '%s' can't have default value: %w", tags[0], fullPath, runtimescan.ErrParseTag)
			}
		} else if fragment == "optional" {
			if hasDefault {
				t.Optional = true
			} else {
				return nil, fmt.Errorf("tag type '%s' of fieldName '%s' can't be optional: %w", tags[0], fullPath, runtimescan.ErrParseTag)
			}
		} else if fragment == "base" {
			if hasDefault {
				t.Base = true
			} else {
				return nil, fmt.Errorf("tag type '%s' of fieldName '%s' can't be base: %w", tags[0], fullPath, runtimescan.ErrParseTag)
			}
		}
	}
	return &t, nil
}
