package runtimescan

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mapEncoder struct {
	result map[string]interface{}
}

func (m mapEncoder) ParseTag(name, tagKey, tagStr, pathStr string, elemType reflect.Type) (tag interface{}, err error) {
	return tagStr, nil
}

func (m *mapEncoder) VisitField(tag, value interface{}) (err error) {
	tagStr := tag.(string)
	m.result[tagStr] = value
	return nil
}

func (m mapEncoder) EnterChild(tag interface{}) (err error) {
	panic("implement me")
}

func (m mapEncoder) LeaveChild(tag interface{}) (err error) {
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
					result: make(map[string]interface{}),
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
					result: make(map[string]interface{}),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t)
		})
	}
}
