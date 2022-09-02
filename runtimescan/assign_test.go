package runtimescan

import (
	"log"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func Test_FuzzyAssign_primitive(t *testing.T) {
	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "bool to float",
			check: func(t *testing.T) {
				var v float64
				src := true
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.InDelta(t, 1.0, v, 0.1)
			},
		},
		{
			name: "bool to *float",
			check: func(t *testing.T) {
				var v *float64
				src := true
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.InDelta(t, 1.0, *v, 0.1)
			},
		},
		{
			name: "*bool to float",
			check: func(t *testing.T) {
				var v float64
				src := true
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.InDelta(t, 1.0, v, 0.1)
			},
		},
		{
			name: "*bool to *float",
			check: func(t *testing.T) {
				var v *float64
				src := true
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.InDelta(t, 1.0, *v, 0.1)
			},
		},

		{
			name: "bool to int",
			check: func(t *testing.T) {
				var v int
				src := true
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.Equal(t, 1, v)
			},
		},
		{
			name: "bool to *int",
			check: func(t *testing.T) {
				var v *int
				src := true
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, 1, *v)
			},
		},
		{
			name: "*bool to int",
			check: func(t *testing.T) {
				var v int
				src := true
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.Equal(t, 1, v)
			},
		},
		{
			name: "*bool to *int",
			check: func(t *testing.T) {
				var v *int
				src := true
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, 1, *v)
			},
		},

		{
			name: "bool to uint",
			check: func(t *testing.T) {
				var v uint
				src := true
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.Equal(t, uint(1), v)
			},
		},
		{
			name: "bool to *uint",
			check: func(t *testing.T) {
				var v *uint
				src := true
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, uint(1), *v)
			},
		},
		{
			name: "*bool to uint",
			check: func(t *testing.T) {
				var v uint
				src := true
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.Equal(t, uint(1), v)
			},
		},
		{
			name: "*bool to *uint",
			check: func(t *testing.T) {
				var v *uint
				src := true
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, uint(1), *v)
			},
		},
		{
			name: "bool to string",
			check: func(t *testing.T) {
				var v string
				src := true
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.Equal(t, "true", v)
			},
		},
		{
			name: "bool to *string",
			check: func(t *testing.T) {
				var v *string
				src := true
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, "true", *v)
			},
		},
		{
			name: "*bool to string",
			check: func(t *testing.T) {
				var v string
				src := true
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.Equal(t, "true", v)
			},
		},
		{
			name: "*bool to *string",
			check: func(t *testing.T) {
				var v *string
				src := true
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, "true", *v)
			},
		},

		{
			name: "int to string",
			check: func(t *testing.T) {
				var v string
				src := 1234
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.Equal(t, "1234", v)
			},
		},
		{
			name: "int to *string",
			check: func(t *testing.T) {
				var v *string
				src := 1234
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, "1234", *v)
			},
		},
		{
			name: "*int to string",
			check: func(t *testing.T) {
				var v string
				src := 1234
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.Equal(t, "1234", v)
			},
		},
		{
			name: "*int to *string",
			check: func(t *testing.T) {
				var v *string
				src := 1234
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, "1234", *v)
			},
		},

		{
			name: "uint to string",
			check: func(t *testing.T) {
				var v string
				src := uint(1234)
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.Equal(t, "1234", v)
			},
		},
		{
			name: "uint to *string",
			check: func(t *testing.T) {
				var v *string
				src := uint(1234)
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, "1234", *v)
			},
		},
		{
			name: "*uint to string",
			check: func(t *testing.T) {
				var v string
				src := uint(1234)
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.Equal(t, "1234", v)
			},
		},
		{
			name: "*uint to *string",
			check: func(t *testing.T) {
				var v *string
				src := uint(1234)
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, "1234", *v)
			},
		},

		{
			name: "float to bool",
			check: func(t *testing.T) {
				var v bool
				src := 1.2345
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.True(t, v)
			},
		},
		{
			name: "float to *bool",
			check: func(t *testing.T) {
				var v *bool
				src := 1.2345
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.True(t, *v)
			},
		},
		{
			name: "*float to bool",
			check: func(t *testing.T) {
				var v bool
				src := 1.2345
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.True(t, v)
			},
		},
		{
			name: "*float to *bool",
			check: func(t *testing.T) {
				var v *bool
				src := 1.2345
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.True(t, *v)
			},
		},

		{
			name: "string to string",
			check: func(t *testing.T) {
				var v string
				src := "test"
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.Equal(t, src, v)
			},
		},
		{
			name: "string to *string",
			check: func(t *testing.T) {
				var v *string
				src := "test"
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, src, *(v))
			},
		},
		{
			name: "*string to string",
			check: func(t *testing.T) {
				var v string
				src := "test"
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.Equal(t, src, v)
			},
		},
		{
			name: "*string to *string",
			check: func(t *testing.T) {
				var v *string
				src := "test"
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, src, *(v))
			},
		},

		{
			name: "int to bool",
			check: func(t *testing.T) {
				var v bool
				src := 1
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.True(t, v)
			},
		},
		{
			name: "int to *bool",
			check: func(t *testing.T) {
				var v *bool
				src := 1
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.True(t, *v)
			},
		},
		{
			name: "*int to bool",
			check: func(t *testing.T) {
				var v bool
				src := 1
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.True(t, v)
			},
		},
		{
			name: "*int to *bool",
			check: func(t *testing.T) {
				var v *bool
				src := 1
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.True(t, *v)
			},
		},

		{
			name: "uint to bool",
			check: func(t *testing.T) {
				var v bool
				src := uint(1)
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.True(t, v)
			},
		},
		{
			name: "uint to *bool",
			check: func(t *testing.T) {
				var v *bool
				src := uint(1)
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.True(t, *v)
			},
		},
		{
			name: "*uint to bool",
			check: func(t *testing.T) {
				var v bool
				src := uint(1)
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.True(t, v)
			},
		},
		{
			name: "*uint to *bool",
			check: func(t *testing.T) {
				var v *bool
				src := uint(1)
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.True(t, *v)
			},
		},

		{
			name: "string to int",
			check: func(t *testing.T) {
				var v int
				src := "12345"
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.Equal(t, 12345, v)
			},
		},
		{
			name: "string to *int",
			check: func(t *testing.T) {
				var v *int
				src := "12345"
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, 12345, *v)
			},
		},
		{
			name: "*string to int",
			check: func(t *testing.T) {
				var v int
				src := "12345"
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.Equal(t, 12345, v)
			},
		},
		{
			name: "*string to *int",
			check: func(t *testing.T) {
				var v *int
				src := "12345"
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, 12345, *v)
			},
		},

		{
			name: "string to uint",
			check: func(t *testing.T) {
				var v uint
				src := "12345"
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.Equal(t, uint(12345), v)
			},
		},
		{
			name: "string to *uint",
			check: func(t *testing.T) {
				var v *uint
				src := "12345"
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, uint(12345), *v)
			},
		},
		{
			name: "*string to uint",
			check: func(t *testing.T) {
				var v uint
				src := "12345"
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.Equal(t, uint(12345), v)
			},
		},
		{
			name: "*string to *uint",
			check: func(t *testing.T) {
				var v *uint
				src := "12345"
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, uint(12345), *v)
			},
		},

		{
			name: "int to float",
			check: func(t *testing.T) {
				var v float64
				src := 12345
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.InDelta(t, 12345.0, v, 0.1)
			},
		},
		{
			name: "int to *float",
			check: func(t *testing.T) {
				var v *float64
				src := 12345
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.InDelta(t, 12345.0, *v, 0.1)
			},
		},
		{
			name: "*int to float",
			check: func(t *testing.T) {
				var v float64
				src := 12345
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.InDelta(t, 12345.0, v, 0.1)
			},
		},
		{
			name: "*int to *float",
			check: func(t *testing.T) {
				var v *float64
				src := 12345
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.InDelta(t, 12345.0, *v, 0.1)
			},
		},

		{
			name: "uint to float",
			check: func(t *testing.T) {
				var v float64
				src := uint(12345)
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.InDelta(t, 12345.0, v, 0.1)
			},
		},
		{
			name: "uint to *float",
			check: func(t *testing.T) {
				var v *float64
				src := uint(12345)
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.InDelta(t, 12345.0, *v, 0.1)
			},
		},
		{
			name: "*uint to float",
			check: func(t *testing.T) {
				var v float64
				src := uint(12345)
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.InDelta(t, 12345.0, v, 0.1)
			},
		},
		{
			name: "*uint to *float",
			check: func(t *testing.T) {
				var v *float64
				src := uint(12345)
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.InDelta(t, 12345.0, *v, 0.1)
			},
		},
		{
			name: "string to bool (true)",
			check: func(t *testing.T) {
				var v bool
				src := "true"
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.True(t, v)
			},
		},
		{
			name: "string to bool (false)",
			check: func(t *testing.T) {
				var v bool
				src := "false"
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.False(t, v)
			},
		},
		{
			name: "string to bool (no)",
			check: func(t *testing.T) {
				var v bool
				src := "no"
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.False(t, v)
			},
		},
		{
			name: "string to *bool",
			check: func(t *testing.T) {
				var v *bool
				src := "true"
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.True(t, *v)
			},
		},
		{
			name: "*string to bool",
			check: func(t *testing.T) {
				var v bool
				src := "true"
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.True(t, v)
			},
		},
		{
			name: "*string to *bool",
			check: func(t *testing.T) {
				var v *bool
				src := "true"
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.True(t, *v)
			},
		},

		{
			name: "float to string",
			check: func(t *testing.T) {
				var v string
				src := 1234.5
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.Equal(t, "1234.5", v)
			},
		},
		{
			name: "float to *string",
			check: func(t *testing.T) {
				var v *string
				src := 1234.5
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, "1234.5", *v)
			},
		},
		{
			name: "*float to string",
			check: func(t *testing.T) {
				var v string
				src := 1234.5
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.Equal(t, "1234.5", v)
			},
		},
		{
			name: "*float to *string",
			check: func(t *testing.T) {
				var v *string
				src := 1234.5
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, "1234.5", *v)
			},
		},
		{
			name: "float to float",
			check: func(t *testing.T) {
				var v float64
				src := 1234.5
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.InDelta(t, src, v, 0.1)
			},
		},
		{
			name: "float to *float",
			check: func(t *testing.T) {
				var v *float64
				src := 1234.5
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.InDelta(t, src, *v, 0.1)
			},
		},
		{
			name: "*float to float",
			check: func(t *testing.T) {
				var v float64
				src := 1234.5
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.InDelta(t, src, v, 0.1)
			},
		},
		{
			name: "*float to *float",
			check: func(t *testing.T) {
				var v *float64
				src := 1234.5
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.InDelta(t, src, *v, 0.1)
			},
		},
		{
			name: "float32 to float64",
			check: func(t *testing.T) {
				var v float64
				src := float32(1234.5)
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.InDelta(t, float64(src), v, 0.1)
			},
		},
		{
			name: "float64 to float32",
			check: func(t *testing.T) {
				var v float32
				src := 1234.5
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.InDelta(t, float32(src), v, 0.1)
			},
		},

		{
			name: "string to float",
			check: func(t *testing.T) {
				var v float64
				src := "1234.5"
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.InDelta(t, 1234.5, v, 0.1)
			},
		},
		{
			name: "string to *float",
			check: func(t *testing.T) {
				var v *float64
				src := "1234.5"
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.InDelta(t, 1234.5, *v, 0.1)
			},
		},
		{
			name: "*string to float",
			check: func(t *testing.T) {
				var v float64
				src := "1234.5"
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.InDelta(t, 1234.5, v, 0.1)
			},
		},
		{
			name: "*string to *float",
			check: func(t *testing.T) {
				var v *float64
				src := "1234.5"
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.InDelta(t, 1234.5, *v, 0.1)
			},
		},

		{
			name: "int to int",
			check: func(t *testing.T) {
				var v int
				src := 12345
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.Equal(t, src, v)
			},
		},
		{
			name: "int to *int",
			check: func(t *testing.T) {
				var v *int
				src := 12345
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, src, *v)
			},
		},
		{
			name: "*int to int",
			check: func(t *testing.T) {
				var v int
				src := 12345
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.Equal(t, src, v)
			},
		},
		{
			name: "*int to *int",
			check: func(t *testing.T) {
				var v *int
				src := 12345
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, src, *v)
			},
		},

		{
			name: "int to #int underlying type",
			check: func(t *testing.T) {
				type Int int
				var v Int
				src := 12345
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.Equal(t, Int(src), v)
			},
		},
		{
			name: "int to *int underlying type",
			check: func(t *testing.T) {
				type Int int
				var v *Int
				src := 12345
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, Int(src), *v)
			},
		},
		{
			name: "*int to int underlying type",
			check: func(t *testing.T) {
				type Int int
				var v Int
				src := 12345
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.Equal(t, Int(src), v)
			},
		},
		{
			name: "*int to *int underlying type",
			check: func(t *testing.T) {
				type Int int
				var v *Int
				src := 12345
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, Int(src), *v)
			},
		},

		{
			name: "int to #int",
			check: func(t *testing.T) {
				type Int int
				var v Int
				src := 12345
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, Int(src), v)
			},
		},
		{
			name: "#int to int",
			check: func(t *testing.T) {
				type Int int
				var v int
				src := Int(12345)
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.Equal(t, int(src), v)
			},
		},
		{
			name: "#int to #int",
			check: func(t *testing.T) {
				type Int int
				var v *Int
				src := Int(12345)
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, src, *v)
			},
		},

		{
			name: "uint to uint",
			check: func(t *testing.T) {
				var v uint
				src := uint(12345)
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.Equal(t, src, v)
			},
		},
		{
			name: "int to int16",
			check: func(t *testing.T) {
				var v int16
				src := 12345
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.Equal(t, int16(src), v)
			},
		},
		{
			name: "int16 to int",
			check: func(t *testing.T) {
				var v int
				src := int16(12345)
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.Equal(t, int(src), v)
			},
		},
		{
			name: "int to uint",
			check: func(t *testing.T) {
				var v uint
				src := 12345
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.Equal(t, uint(src), v)
			},
		},
		{
			name: "uint to int",
			check: func(t *testing.T) {
				var v int
				src := uint(12345)
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.Equal(t, int(src), v)
			},
		},
		{
			name: "uint to *uint",
			check: func(t *testing.T) {
				var v *uint
				src := uint(12345)
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, src, *v)
			},
		},
		{
			name: "*uint to uint",
			check: func(t *testing.T) {
				var v uint
				src := uint(12345)
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.Equal(t, src, v)
			},
		},
		{
			name: "*uint to *uint",
			check: func(t *testing.T) {
				var v *uint
				src := uint(12345)
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, src, *v)
			},
		},

		{
			name: "bool to bool",
			check: func(t *testing.T) {
				var v bool
				src := true
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.True(t, v)
			},
		},
		{
			name: "bool to *bool",
			check: func(t *testing.T) {
				var v *bool
				src := true
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.True(t, *v)
			},
		},
		{
			name: "*bool to bool",
			check: func(t *testing.T) {
				var v bool
				src := true
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.True(t, v)
			},
		},
		{
			name: "*bool to *bool",
			check: func(t *testing.T) {
				var v *bool
				src := true
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.True(t, *v)
			},
		},

		{
			name: "float to int",
			check: func(t *testing.T) {
				var v int
				src := 12345.6
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.Equal(t, 12345, v)
			},
		},
		{
			name: "float to *int",
			check: func(t *testing.T) {
				var v *int
				src := 12345.6
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, 12345, *v)
			},
		},
		{
			name: "*float to int",
			check: func(t *testing.T) {
				var v int
				src := 12345.6
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.Equal(t, 12345, v)
			},
		},
		{
			name: "*float to *int",
			check: func(t *testing.T) {
				var v *int
				src := 12345.6
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, 12345, *v)
			},
		},

		{
			name: "float to uint",
			check: func(t *testing.T) {
				var v uint
				src := 12345.6
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.Equal(t, uint(12345), v)
			},
		},
		{
			name: "float to *uint",
			check: func(t *testing.T) {
				var v *uint
				src := 12345.6
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, uint(12345), *v)
			},
		},
		{
			name: "*float to uint",
			check: func(t *testing.T) {
				var v uint
				src := 12345.6
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.Equal(t, uint(12345), v)
			},
		},
		{
			name: "*float to *uint",
			check: func(t *testing.T) {
				var v *uint
				src := 12345.6
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, uint(12345), *v)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t)
		})
	}
}

type E struct {
	s string
}

func (e E) Error() string {
	return e.s
}

func Test_FuzzyAssign_struct(t *testing.T) {
	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "struct to *struct",
			check: func(t *testing.T) {
				var v *E
				src := E{
					s: "test",
				}
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, "test", v.s)
			},
		},
		{
			name: "*struct to *struct",
			check: func(t *testing.T) {
				var v *E
				src := E{
					s: "test",
				}
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, "test", v.s)

			},
		},
		{
			name: "struct to interface",
			check: func(t *testing.T) {
				var v error
				src := E{
					s: "test",
				}
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, "test", v.Error())
			},
		},
		{
			name: "*struct to interface",
			check: func(t *testing.T) {
				var v error
				src := E{
					s: "test",
				}
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, "test", v.Error())
			},
		},
		{
			name: "struct to interface",
			check: func(t *testing.T) {
				var v error
				src := E{
					s: "test",
				}
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.NotNil(t, v)
				assert.Equal(t, "test", v.Error())
			},
		},
		{
			name: "user defined type",
			check: func(t *testing.T) {
				type Int int
				var v Int
				src := 12345
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.Equal(t, Int(src), v)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t)
		})
	}
}

func Test_directAssign(t *testing.T) {
	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "int to int",
			check: func(t *testing.T) {
				var i int = 10
				var o int
				ok := directAssign(reflect.ValueOf(&o), reflect.ValueOf(i))
				assert.True(t, ok)
				assert.Equal(t, i, o)
			},
		},
		{
			name: "int to int (user defined to base)",
			check: func(t *testing.T) {
				type Int int
				var i Int = 10
				var o int
				ok := directAssign(reflect.ValueOf(&o), reflect.ValueOf(i))
				assert.True(t, ok)
				assert.Equal(t, i, Int(o))
			},
		},
		{
			name: "int to int (base to user defined)",
			check: func(t *testing.T) {
				type Int int
				var i int = 10
				var o Int
				ok := directAssign(reflect.ValueOf(&o), reflect.ValueOf(i))
				assert.True(t, ok)
				assert.Equal(t, Int(i), o)
			},
		},
		{
			name: "int to int (user defined to user defined)",
			check: func(t *testing.T) {
				type Int int
				var i Int = 10
				var o Int
				ok := directAssign(reflect.ValueOf(&o), reflect.ValueOf(i))
				assert.True(t, ok)
				assert.Equal(t, i, o)
			},
		},
		{
			name: "int16 to int",
			check: func(t *testing.T) {
				var i int16 = 10
				var o int
				ok := directAssign(reflect.ValueOf(&o), reflect.ValueOf(i))
				assert.True(t, ok)
				assert.Equal(t, int(i), o)
			},
		},
		{
			name: "int to int16",
			check: func(t *testing.T) {
				var i int = 10
				var o int16
				ok := directAssign(reflect.ValueOf(&o), reflect.ValueOf(i))
				assert.True(t, ok)
				assert.Equal(t, int16(i), o)
			},
		},
		{
			name: "string to int",
			check: func(t *testing.T) {
				var i string = "hello"
				var o int
				ok := directAssign(reflect.ValueOf(&o), reflect.ValueOf(i))
				assert.False(t, ok)
				assert.NotEqual(t, i, o)
			},
		},
		{
			name: "string to string",
			check: func(t *testing.T) {
				var i string = "hello"
				var o string
				ok := directAssign(reflect.ValueOf(&o), reflect.ValueOf(i))
				assert.True(t, ok)
				assert.Equal(t, i, o)
			},
		},
		{
			name: "struct to struct",
			check: func(t *testing.T) {
				type S struct {
					Name string
				}
				var i S = S{Name: "hello"}
				var o S
				ok := directAssign(reflect.ValueOf(&o), reflect.ValueOf(i))
				assert.True(t, ok)
				assert.Equal(t, i, o)
			},
		},
		{
			name: "int to string should not be ok",
			check: func(t *testing.T) {
				var v string
				src := 1234
				ok := directAssign(reflect.ValueOf(&v), reflect.ValueOf(src))
				assert.False(t, ok)
				assert.Equal(t, v, "")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t)
		})
	}
}

func Test_unwrap(t *testing.T) {
	tests := []struct {
		name  string
		check func(t *testing.T)
	}{
		{
			name: "*v to *v (not convert)",
			check: func(t *testing.T) {
				var i int
				v, ok := unwrap(reflect.ValueOf(&i), true)
				assert.True(t, ok)
				assert.Equal(t, 1, countPointerDepth(v))
			},
		},
		{
			name: "*v to v",
			check: func(t *testing.T) {
				var i int
				v, ok := unwrap(reflect.ValueOf(&i), false)
				assert.True(t, ok)
				assert.Equal(t, 0, countPointerDepth(v))
			},
		},
		{
			name: "**v to v",
			check: func(t *testing.T) {
				var i int
				ip := &i
				v, ok := unwrap(reflect.ValueOf(&ip), false)
				assert.True(t, ok)
				assert.Equal(t, 0, countPointerDepth(v))
			},
		},
		{
			name: "**v to *v",
			check: func(t *testing.T) {
				var i int
				ip := &i
				v, ok := unwrap(reflect.ValueOf(&ip), true)
				assert.True(t, ok)
				assert.Equal(t, 1, countPointerDepth(v))
			},
		},
		{
			name: "v to *v (error)",
			check: func(t *testing.T) {
				var i int
				v, ok := unwrap(reflect.ValueOf(i), true)
				assert.False(t, ok)
				assert.False(t, v.IsValid())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t)
		})
	}
}
