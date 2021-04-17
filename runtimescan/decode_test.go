package runtimescan

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mapDecoder struct {
	values map[string]interface{}
}

func (d mapDecoder) ParseTag(name, tagKey, tagStr, pathStr string, eType reflect.Type) (interface{}, error) {
	return tagStr, nil
}

func (d mapDecoder) ExtractValue(tag interface{}) (interface{}, error) {
	tagStr := tag.(string)
	v, ok := d.values[tagStr]
	if ok {
		return v, nil
	}
	return nil, Skip
}

type SampleStruct struct {
	Value string
}

func Test_decode(t *testing.T) {
	d := mapDecoder{
		values: map[string]interface{}{
			"int":    12345,
			"uint":   uint(12345),
			"string": "string",
			"float":  123.45,
			"bool":   true,
			"struct": &SampleStruct{
				Value: "struct sample",
			},
			"interface": errors.New("error interface"),
		},
	}
	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "assign values",
			check: func(t *testing.T) {
				type Target struct {
					Int    int    `map:"int"`
					String string `map:"string"`
				}
				target := Target{}
				v, err := newParser(&d, []string{"map"}, &target)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				err = decode(&target, v, &d)
				assert.NoError(t, err)
				assert.Equal(t, 12345, target.Int)
				assert.Equal(t, "string", target.String)
			},
		},
		{
			name: "skip assigning values",
			check: func(t *testing.T) {
				type Target struct {
					Int    int    `map:"not-found"` // if tag is not found, test decoder returns Skip
					String string `map:"not-found"` // if tag is not found, test decoder returns Skip
				}
				target := Target{}
				v, err := newParser(&d, []string{"map"}, &target)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				err = decode(&target, v, &d)
				assert.NoError(t, err)
				assert.Equal(t, 0, target.Int)
				assert.Equal(t, "", target.String)
			},
		},
		{
			name: "skip private values",
			check: func(t *testing.T) {
				type Target struct {
					int    int    `map:"int"`    // this is private
					string string `map:"string"` // this is private
				}
				target := Target{}
				v, err := newParser(&d, []string{"map"}, &target)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				err = decode(&target, v, &d)
				assert.NoError(t, err)
				assert.Equal(t, 0, target.int)
				assert.Equal(t, "", target.string)
			},
		},
		{
			name: "struct and interface",
			check: func(t *testing.T) {
				type Target struct {
					Sample *SampleStruct `map:"struct"`
					Error  error         `map:"interface"`
				}
				target := Target{}
				v, err := newParser(&d, []string{"map"}, &target)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				err = decode(&target, v, &d)
				assert.NoError(t, err)
				assert.NotNil(t, target.Sample)
				assert.Equal(t, "struct sample", target.Sample.Value)
				assert.NotNil(t, target.Error)
				assert.Equal(t, "error interface", target.Error.Error())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t)
		})
	}
}
