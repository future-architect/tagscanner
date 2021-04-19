package binarypatternmatch

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"

	"gitlab.com/osaki-lab/tagscanner/runtimescan"
)

var (
	ErrNotMatch = errors.New("not match")
)

type binaryDecoder struct {
	reader       io.Reader
	remainedByte byte
	offset       int
	error        error
	lengths      map[string]uint64
	buffer       []byte
}

func (b binaryDecoder) ParseTag(name, tagKey, tagStr, pathStr string, elemType reflect.Type) (tag interface{}, err error) {
	return parseBinaryTag(name, tagKey, tagStr, pathStr, elemType)
}

func (b *binaryDecoder) ExtractValue(tag interface{}) (value interface{}, err error) {
	t := tag.(*binaryTag)
	switch t.Type {
	case Bits:
		err := b.readBits(t.Bits)
		if err != nil {
			return nil, err
		}
	case Bytes:
		err := b.readBits(t.Bits)
		if err != nil {
			return nil, err
		}
	case Length:
	}

	// pattern match
	if len(t.Literal) > 0 && !bytes.Equal(t.Literal, b.buffer){
		return nil, ErrNotMatch
	}

	switch t.EType.Kind() {
	case reflect.Bool:
		if b.buffer[0] == 0 {
			return false, nil
		} else {
			return true, nil
		}
	case reflect.Uint8:
		return uint(maskBits(b.buffer[0], t.Bits)), nil
	case reflect.Uint16:
	case reflect.Int32:
		if len(b.buffer) <= 4 {
			buffer := make([]byte, 4)
			copy(buffer[4-len(b.buffer):], b.buffer)
			return int(binary.BigEndian.Uint32(buffer)), nil
		}
	case reflect.Uint32:
		if len(b.buffer) < 4 {
			buffer := make([]byte, 4)
			copy(buffer[4-len(b.buffer):], b.buffer)
			return uint(binary.BigEndian.Uint32(buffer)), nil
		}
		return uint(binary.BigEndian.Uint32(b.buffer)), nil
	case reflect.Uint64:
		if len(b.buffer) < 8 {
			buffer := make([]byte, 8)
			copy(buffer[8-len(b.buffer):], b.buffer)
			return uint(binary.BigEndian.Uint64(buffer)), nil
		}
		return uint(binary.BigEndian.Uint64(b.buffer)), nil
	case reflect.String:
		return string(b.buffer), nil
	case reflect.Slice:
		if t.EType.Elem().Kind() == reflect.Uint8 {
			return b.buffer, nil
		}
	}
	return nil, fmt.Errorf("error: unsupported type to decode to variable: %v", t.EType.Kind())
}

func (b *binaryDecoder) readBits(bits int) error {
	if b.offset == 8 {
		byl := bits / 8
		remain := false
		if bits%8 != 0 {
			byl += 1
			remain = true
		}

		if !remain {
			err := b.readBytes(byl)
			if err != nil {
				return err
			}
			b.remainedByte = 0
			return nil
		}
		err := b.readBytes(byl)
		if err != nil {
			return err
		}
		b.offset = (b.offset + bits) % 8
		b.remainedByte = b.buffer[len(b.buffer)-1]
		return nil
	} else {
		needToRead := bits - (8 - b.offset)
		if needToRead > 0 {
			ntrb := needToRead / 8
			if needToRead % 7 != 0 {
				ntrb += 1
			}
			if ntrb > 0 {
				b.readBytes(ntrb)
			} else {
				b.buffer = b.buffer[:0]
			}
		} else {
			b.buffer = b.buffer[:0]
		}
		b.buffer, b.remainedByte, b.offset = copyWithOffset(b.remainedByte, b.buffer, b.offset, bits)
		return nil
	}
	return errors.New("🦑")
}

func copyWithOffset(preload byte, src []byte, offset, bits int) (result []byte, remainByte byte, newOffset int) {
	byl := bits / 8
	if bits % 8 != 0 {
		byl += 1
	}
	result = make([]byte, 0, byl)
	remainByte = preload
	if len(src) == 0 {
		result = append(result, maskBits(preload >> offset, bits))
		newOffset = offset + bits
		if newOffset == 8 {
			remainByte = 0
		}
		return
	}
	for _, b := range src {
		if bits < 8 {
			result = append(result, maskBits(b >> offset, bits) << offset | remainByte)
			remainByte = b
			newOffset = offset + bits
			bits = 0
		} else {
			result = append(result, remainByte >> offset | maskBits(b, offset) << (8-offset))
			newOffset = offset
			bits -= 8
			if newOffset == 8 {
				remainByte = 0
			} else {
				remainByte = b
			}
		}
	}
	if bits > 0 {
		result = append(result, maskBits(remainByte >> offset, bits))
		newOffset = offset + bits
		if newOffset == 8 {
			remainByte = 0
		}
	}
	return
}

func maskBits(i byte, l int) byte {
	switch l {
	case 1:
		return i & 0b00000001
	case 2:
		return i & 0b00000011
	case 3:
		return i & 0b00000111
	case 4:
		return i & 0b00001111
	case 5:
		return i & 0b00011111
	case 6:
		return i & 0b00111111
	case 7:
		return i & 0b01111111
	case 8:
		return i
	default:
		panic(fmt.Sprintf("maskBits only accepts 0 < l && l < 8, but %d", l))
	}
}

func (b *binaryDecoder) readBytes(l int) error {
	if len(b.buffer) == l {
		// do nothing
	} else if len(b.buffer) > l || cap(b.buffer) >= l {
		b.buffer = b.buffer[:l]
	} else { // len(buffer) < l && cap(buffer) < l
		b.buffer = make([]byte, l)
	}
	n, err := io.ReadFull(b.reader, b.buffer)
	if err != nil {
		return err
	}
	if n != l {
		panic(fmt.Sprintf("readBits %d/%d", n, l))
	}
	return nil
}

func Decode(dest interface{}, reader io.Reader) error {
	decoder := &binaryDecoder{
		reader: reader,
		lengths: make(map[string]uint64),
		offset: 8,
	}
	return runtimescan.Decode(dest, []string{"bits", "bytes", "length"}, decoder)
}
