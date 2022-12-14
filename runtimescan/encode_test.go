package runtimescan

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mapEncoder struct {
	result map[string]any
}

func (m mapEncoder) ParseTag(name, tagKey, tagStr, pathStr string, elemType reflect.Type) (tag any, err error) {
	return tagStr, nil
}

func (m *mapEncoder) VisitField(tag, value any) (err error) {
	tagStr := tag.(string)
	m.result[tagStr] = value
	return nil
}

func (m mapEncoder) EnterChild(tag any) (err error) {
	panic("implement me")
}

func (m mapEncoder) LeaveChild(tag any) (err error) {
	panic("implement me")
}

func Test_encode(t *testing.T) {
	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "simple tags",
			check: func(t *testing.T) {
				m := mapEncoder{
					result: make(map[string]any),
				}
				type Source struct {
					Int       int     `map:"int"`
					String    string  `map:"string"`
					PtrInt    *int    `map:"ptr_int"`
					PtrString *string `map:"ptr_str"`
				}
				source := Source{
					Int:       12345,
					String:    "test string",
					PtrInt:    &[]int{10}[0],
					PtrString: nil,
				}
				v, err := newParser(&m, []string{"map"}, &source)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				err = encode(&m, v, &source)
				assert.NoError(t, err)
				assert.Equal(t, 12345, m.result["int"])
				assert.Equal(t, "test string", m.result["string"])
				assert.Equal(t, 10, m.result["ptr_int"])
				assert.Equal(t, nil, m.result["ptr_string"])
			},
		},
		{
			name: "nested struct",
			check: func(t *testing.T) {
				m := mapEncoder{
					result: make(map[string]any),
				}
				type Child struct {
					Int    int    `map:"int"`
					String string `map:"string"`
				}
				type Source struct {
					Child Child
				}
				source := Source{
					Child: Child{
						Int:    12345,
						String: "test string",
					},
				}
				v, err := newParser(&m, []string{"map"}, &source)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				err = encode(&m, v, &source)
				assert.NoError(t, err)
				assert.Equal(t, 12345, m.result["int"])
				assert.Equal(t, "test string", m.result["string"])
			},
		},
		{
			name: "user defined type",
			check: func(t *testing.T) {
				m := mapEncoder{
					result: make(map[string]any),
				}
				type Int int
				type String string
				type Source struct {
					Int       Int     `map:"int"`
					String    String  `map:"string"`
					PtrInt    *Int    `map:"ptr_int"`
					PtrString *String `map:"ptr_str"`
				}
				source := Source{
					Int:       12345,
					String:    "test string",
					PtrInt:    &[]Int{10}[0],
					PtrString: nil,
				}
				v, err := newParser(&m, []string{"map"}, &source)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				err = encode(&m, v, &source)
				assert.NoError(t, err)
				assert.Equal(t, Int(12345), m.result["int"])
				assert.Equal(t, String("test string"), m.result["string"])
				assert.Equal(t, Int(10), m.result["ptr_int"])
				assert.Equal(t, nil, m.result["ptr_string"])
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t)
		})
	}
}
