package runtimescan

import (
	"log"
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
			name: "string to bool",
			check: func(t *testing.T) {
				var v bool
				src := "true"
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.True(t, v)
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
				assert.Equal(t, 12346, v)
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
				assert.Equal(t, 12346, *v)
			},
		},
		{
			name: "*float to int",
			check: func(t *testing.T) {
				var v int
				src := 12345.6
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.Equal(t, 12346, v)
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
				assert.Equal(t, 12346, *v)
			},
		},

		{
			name: "float to uint",
			check: func(t *testing.T) {
				var v uint
				src := 12345.6
				err := FuzzyAssign(&v, src)
				assert.NoError(t, err)
				assert.Equal(t, uint(12346), v)
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
				assert.Equal(t, uint(12346), *v)
			},
		},
		{
			name: "*float to uint",
			check: func(t *testing.T) {
				var v uint
				src := 12345.6
				err := FuzzyAssign(&v, &src)
				assert.NoError(t, err)
				assert.Equal(t, uint(12346), v)
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
				assert.Equal(t, uint(12346), *v)
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
