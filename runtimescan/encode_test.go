package runtimescan

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type mapEncoder struct {
	result map[string]interface{}
}

func (m mapEncoder) ParseTag(name, tagStr, pathStr string, elemType reflect.Type) (tag interface{}, err error) {
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
		name    string
		check   func(t *testing.T)
	}{
		{
			name: "simple tags",
			check: func(t *testing.T) {
				m := mapEncoder{
					result: make(map[string]interface{}),
				}
				type Source struct {
					Int    int    `map:"int"`
					String string `map:"string"`
				}
				source := Source{
					Int: 12345,
					String: "test string",
				}
				v, err := newParser(&m, "map", &source)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				err = encode(&m, v, &source)
				assert.NoError(t, err)
				assert.Equal(t, 12345, m.result["int"])
				assert.Equal(t, "test string", m.result["string"])
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
						Int: 12345,
						String: "test string",
					},
				}
				v, err := newParser(&m, "map", &source)
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