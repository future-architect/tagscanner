package runtimescan

import (
	"fmt"
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

func Test_Encode_parallel(t *testing.T) {
	parallelNum := 1000
	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "simple tags",
			check: func(t *testing.T) {
				type Source struct {
					Int       int     `map:"int"`
					String    string  `map:"string"`
					PtrInt    *int    `map:"ptr_int"`
					PtrString *string `map:"ptr_str"`
				}
				sources := []Source{
					{
						Int:       12345678900,
						String:    "test string 0",
						PtrInt:    &[]int{0}[0],
						PtrString: nil,
					},
					{
						Int:       12345678901,
						String:    "test string 1",
						PtrInt:    &[]int{1}[0],
						PtrString: nil,
					},
					{
						Int:       12345678902,
						String:    "test string 2",
						PtrInt:    &[]int{2}[0],
						PtrString: nil,
					},
					{
						Int:       12345678903,
						String:    "test string 3",
						PtrInt:    &[]int{3}[0],
						PtrString: nil,
					},
					{
						Int:       12345678904,
						String:    "test string 4",
						PtrInt:    &[]int{4}[0],
						PtrString: nil,
					},
					{
						Int:       12345678905,
						String:    "test string 5",
						PtrInt:    &[]int{5}[0],
						PtrString: nil,
					},
					{
						Int:       12345678906,
						String:    "test string 6",
						PtrInt:    &[]int{6}[0],
						PtrString: nil,
					},
					{
						Int:       12345678907,
						String:    "test string 7",
						PtrInt:    &[]int{7}[0],
						PtrString: nil,
					},
					{
						Int:       12345678908,
						String:    "test string 8",
						PtrInt:    &[]int{8}[0],
						PtrString: nil,
					},
					{
						Int:       12345678909,
						String:    "test string 9",
						PtrInt:    &[]int{9}[0],
						PtrString: nil,
					},
					{
						Int:       12345678910,
						String:    "test string 10",
						PtrInt:    &[]int{10}[0],
						PtrString: nil,
					},
				}
				for i := 0; i < parallelNum; i++ {
					num := i
					index := num % len(sources)
					source := sources[index]
					go func() {
						m := mapEncoder{
							result: make(map[string]any),
						}

						err := Encode(&source, []string{"map"}, &m)
						assert.NoError(t, err)
						assert.Equal(t, int(123456789*100+index), m.result["int"])
						assert.Equal(t, string(fmt.Sprintf("test string %v", index)), m.result["string"])
						assert.Equal(t, int(index), m.result["ptr_int"])
						assert.Equal(t, nil, m.result["ptr_string"])
					}()
				}
			},
		},
		{
			name: "nested struct",
			check: func(t *testing.T) {
				type Child struct {
					Int    int    `map:"int"`
					String string `map:"string"`
				}
				type Source struct {
					Child Child
				}
				sources := []Source{
					{
						Child: Child{
							Int:    12345678900,
							String: "test string 0",
						},
					},
					{
						Child: Child{
							Int:    12345678901,
							String: "test string 1",
						},
					},
					{
						Child: Child{
							Int:    12345678902,
							String: "test string 2",
						},
					},
					{
						Child: Child{
							Int:    12345678903,
							String: "test string 3",
						},
					},
					{
						Child: Child{
							Int:    12345678904,
							String: "test string 4",
						},
					},
					{
						Child: Child{
							Int:    12345678905,
							String: "test string 5",
						},
					},
					{
						Child: Child{
							Int:    12345678906,
							String: "test string 6",
						},
					},
					{
						Child: Child{
							Int:    12345678907,
							String: "test string 7",
						},
					},
					{
						Child: Child{
							Int:    12345678908,
							String: "test string 8",
						},
					},
					{
						Child: Child{
							Int:    12345678909,
							String: "test string 9",
						},
					},
					{
						Child: Child{
							Int:    12345678910,
							String: "test string 10",
						},
					},
				}
				for i := 0; i < parallelNum; i++ {
					num := i
					index := num % len(sources)
					source := sources[index]
					go func() {
						m := mapEncoder{
							result: make(map[string]any),
						}
						err := Encode(&source, []string{"map"}, &m)
						assert.NoError(t, err)
						assert.Equal(t, 123456789*100+index, m.result["int"])
						assert.Equal(t, fmt.Sprintf("test string %v", index), m.result["string"])
					}()
				}
			},
		},
		{
			name: "user defined type",
			check: func(t *testing.T) {
				type Int int
				type String string
				type Source struct {
					Int       Int     `map:"int"`
					String    String  `map:"string"`
					PtrInt    *Int    `map:"ptr_int"`
					PtrString *String `map:"ptr_str"`
				}

				sources := []Source{
					{
						Int:       12345678900,
						String:    "test string 0",
						PtrInt:    &[]Int{0}[0],
						PtrString: nil,
					},
					{
						Int:       12345678901,
						String:    "test string 1",
						PtrInt:    &[]Int{1}[0],
						PtrString: nil,
					},
					{
						Int:       12345678902,
						String:    "test string 2",
						PtrInt:    &[]Int{2}[0],
						PtrString: nil,
					},
					{
						Int:       12345678903,
						String:    "test string 3",
						PtrInt:    &[]Int{3}[0],
						PtrString: nil,
					},
					{
						Int:       12345678904,
						String:    "test string 4",
						PtrInt:    &[]Int{4}[0],
						PtrString: nil,
					},
					{
						Int:       12345678905,
						String:    "test string 5",
						PtrInt:    &[]Int{5}[0],
						PtrString: nil,
					},
					{
						Int:       12345678906,
						String:    "test string 6",
						PtrInt:    &[]Int{6}[0],
						PtrString: nil,
					},
					{
						Int:       12345678907,
						String:    "test string 7",
						PtrInt:    &[]Int{7}[0],
						PtrString: nil,
					},
					{
						Int:       12345678908,
						String:    "test string 8",
						PtrInt:    &[]Int{8}[0],
						PtrString: nil,
					},
					{
						Int:       12345678909,
						String:    "test string 9",
						PtrInt:    &[]Int{9}[0],
						PtrString: nil,
					},
					{
						Int:       12345678910,
						String:    "test string 10",
						PtrInt:    &[]Int{10}[0],
						PtrString: nil,
					},
				}

				for i := 0; i < parallelNum; i++ {
					num := i
					index := num % len(sources)
					source := sources[index]
					go func() {
						m := mapEncoder{
							result: make(map[string]any),
						}

						err := Encode(&source, []string{"map"}, &m)
						assert.NoError(t, err)
						assert.Equal(t, Int(123456789*100+index), m.result["int"])
						assert.Equal(t, String(fmt.Sprintf("test string %v", index)), m.result["string"])
						assert.Equal(t, Int(index), m.result["ptr_int"])
						assert.Equal(t, nil, m.result["ptr_string"])
					}()
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t)
		})
	}
}
