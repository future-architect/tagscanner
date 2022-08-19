package binarypatternmatch

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/future-architect/tagscanner/runtimescan"
)

type TagKey string

const (
	Bits   TagKey = "bits"
	Bytes  TagKey = "bytes"
	Length TagKey = "length"
)

type binaryTag struct {
	Type    TagKey
	EType   reflect.Type
	Literal []byte
	Bits    int
	Label   string
}

func maxBits(k reflect.Kind) int {
	switch k {
	case reflect.Bool:
		return 1
	case reflect.Uint8:
		return 8
	case reflect.Int8:
		return 8
	case reflect.Uint16:
		return 16
	case reflect.Int16:
		return 16
	case reflect.Uint32:
		return 32
	case reflect.Int32:
		return 32
	case reflect.Uint64:
		return 64
	case reflect.Int64:
		return 64
	case reflect.Uint:
		return 32
	case reflect.Int:
		return 32
	case reflect.Float32:
		return 32
	case reflect.Float64:
		return 64
	case reflect.String:
		return -1
	case reflect.Slice:
		return -1
	default:
		panic(k)
	}
}

func isLiteral(tagValue string) bool {
	return strings.HasPrefix(tagValue, "<<") && strings.HasSuffix(tagValue, ">>")
}

func parseLiteral(tagValue string, defaultBits, unit int) ([]byte, int, error) {
	literal := strings.TrimSuffix(strings.TrimPrefix(tagValue, "<<"), ">>")
	i := strings.LastIndex(literal, "/")
	length := defaultBits
	if i != -1 {
		l, err := strconv.ParseUint(literal[i+1:], 10, 64)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid literal(parse size option error) '%s': %w", tagValue, runtimescan.ErrParseTag)
		}
		length = int(l) * unit
		literal = literal[:i]
	}
	var base int
	var label string
	strLiteral := false
	if strings.HasPrefix(literal, "0x") {
		base = 16
		label = "hex"
	} else if strings.HasPrefix(literal, "0b") {
		base = 2
		label = "bin"
	} else {
		strLiteral = true
	}
	var buffer []byte
	if strLiteral {
		buffer = []byte(literal)
	} else {
		v, err := strconv.ParseUint(literal[2:], base, 64)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid literal(parse %s error) '%s': %w", label, tagValue, runtimescan.ErrParseTag)
		}
		buffer = make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, v)
	}
	bytes := length / 8
	if bytes*8 < length {
		bytes += 1
	}
	return buffer[len(buffer)-bytes:], length, nil
}

func parseBinaryTag(fieldName, tagKey, tagValue, fullPath string, eType reflect.Type) (*binaryTag, error) {
	switch tagKey {
	case "bits":
		var length int
		literal := []byte{}
		mb := maxBits(eType.Kind())
		if isLiteral(tagValue) {
			var err error
			literal, length, err = parseLiteral(tagValue, mb, 1)
			if err != nil {
				return nil, err
			}
		} else if isLabel(tagValue) {
			return &binaryTag{
				Type:  Bits,
				EType: eType,
				Label: tagValue,
			}, nil
		} else {
			v, err := strconv.ParseUint(tagValue, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("pattern `%s` syntax error: invalid sub field: %w", tagValue, runtimescan.ErrParseTag)
			}
			length = int(v)
		}
		if length > mb && mb > 0 {
			return nil, fmt.Errorf("pattern `%s` length error: %s's max length is %d but %d: %w", tagValue, eType.Kind(), maxBits(eType.Kind()), length, runtimescan.ErrParseTag)
		}
		return &binaryTag{
			Type:    Bits,
			EType:   eType,
			Literal: literal,
			Bits:    length,
			Label:   "",
		}, nil
	case "bytes":
		var length int
		literal := []byte{}
		mb := maxBits(eType.Kind())
		if isLiteral(tagValue) {
			var err error
			literal, length, err = parseLiteral(tagValue, mb, 8)
			if err != nil {
				return nil, err
			}
		} else if isLabel(tagValue) {
			return &binaryTag{
				Type:  Bytes,
				EType: eType,
				Label: tagValue,
			}, nil
		} else {
			v, err := strconv.ParseUint(tagValue, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("pattern `%s` syntax error: invalid sub field: %w", tagValue, runtimescan.ErrParseTag)
			}
			length = int(v) * 8
		}
		if length > mb && mb > 0 {
			return nil, fmt.Errorf("pattern `%s` length error: %s's max length is %d but %d: %w", tagValue, eType.Kind(), maxBits(eType.Kind())/8, length/8, runtimescan.ErrParseTag)
		}
		return &binaryTag{
			Type:    Bytes,
			EType:   eType,
			Literal: literal,
			Bits:    length,
			Label:   "",
		}, nil
	case "length":
		label := tagValue
		i := strings.LastIndex(tagValue, "/")
		length := maxBits(eType.Kind())
		if i != -1 {
			label = tagValue[:i]
			l, err := strconv.ParseUint(tagValue[i+1:], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid option(parse length size option error) '%s': %w", tagValue, runtimescan.ErrParseTag)
			}
			length = int(l)
		}
		if !strings.HasPrefix(tagValue, "@") {
			return nil, fmt.Errorf("invalid length label (label should starts with '@') '%s': %w", tagValue, runtimescan.ErrParseTag)
		}
		return &binaryTag{
			Type:  Length,
			EType: eType,
			Bits:  length,
			Label: label,
		}, nil
	case "":
		return &binaryTag{
			Type:    Bits,
			EType:   eType,
			Literal: []byte{},
			Bits:    maxBits(eType.Kind()),
			Label:   "",
		}, nil
	default:
		return nil, fmt.Errorf("pattern `%s` syntax error: unknown format: %w", tagValue, runtimescan.ErrParseTag)
	}
	return nil, fmt.Errorf("pattern `%s` syntax error: unknown format: %w", tagValue, runtimescan.ErrParseTag)
}

func isLabel(tagValue string) bool {
	return strings.HasPrefix(tagValue, "@")
}
