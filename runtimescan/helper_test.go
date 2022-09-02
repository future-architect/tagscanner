package runtimescan

import (
	"reflect"
	"testing"
)

func TestStr2PrimitiveValue(t *testing.T) {
	type args struct {
		value    string
		elemType reflect.Type
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "int",
			args: args{
				value:    "12345",
				elemType: reflect.TypeOf(12345),
			},
			want:    int64(12345),
			wantErr: false,
		},
		{
			name: "uint",
			args: args{
				value:    "12345",
				elemType: reflect.TypeOf(uint(12345)),
			},
			want:    uint64(12345),
			wantErr: false,
		},
		{
			name: "float",
			args: args{
				value:    "12345.6",
				elemType: reflect.TypeOf(float64(12345)),
			},
			want:    12345.6,
			wantErr: false,
		},
		{
			name: "bool",
			args: args{
				value:    "true",
				elemType: reflect.TypeOf(true),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "string",
			args: args{
				value:    "hello world",
				elemType: reflect.TypeOf("hello"),
			},
			want:    "hello world",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Str2PrimitiveValue(tt.args.value, tt.args.elemType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Str2PrimitiveValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Str2PrimitiveValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPointerOfSliceOfStruct(t *testing.T) {
	type args struct {
		i any
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "true: pointer of slice of struct",
			args: args{
				i: &[]struct{}{},
			},
			want: true,
		},
		{
			name: "false: slice of struct",
			args: args{
				i: []struct{}{},
			},
			want: false,
		},
		{
			name: "false: pointer of slice of int",
			args: args{
				i: &[]int{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPointerOfSliceOfStruct(tt.args.i); got != tt.want {
				t.Errorf("IsSliceOfStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPointerOfStruct(t *testing.T) {
	type args struct {
		i any
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "true: pointer of struct",
			args: args{
				i: &struct{}{},
			},
			want: true,
		},
		{
			name: "false: struct",
			args: args{
				i: struct{}{},
			},
			want: false,
		},
		{
			name: "false: pointer of int",
			args: args{
				i: &[]int{1}[0],
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPointerOfStruct(tt.args.i); got != tt.want {
				t.Errorf("IsPointerOfStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPointerOfSliceOfPointerOfStruct(t *testing.T) {
	type args struct {
		i any
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "true: pointer of slice of pointer of struct",
			args: args{
				i: &[]*struct{}{},
			},
			want: true,
		},
		{
			name: "false: slice of pointer of struct",
			args: args{
				i: []*struct{}{},
			},
			want: false,
		},
		{
			name: "false: pointer of slice of pointer of int",
			args: args{
				i: &[]*int{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPointerOfSliceOfPointerOfStruct(tt.args.i); got != tt.want {
				t.Errorf("IsPointerOfSliceOfPointerOfStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStructInstance(t *testing.T) {
	type Struct struct {
		Name string
	}
	type args struct {
		i any
	}
	tests := []struct {
		name string
		args args
		test func(t *testing.T, r any, err error)
	}{
		{
			name: "success: instance of struct",
			args: args{
				i: &Struct{},
			},
			test: func(t *testing.T, r any, err error) {
				if _, ok := r.(*Struct); !ok {
					t.Errorf("result should be valid instance")
				}
				if err != nil {
					t.Errorf("error should not be nil")
				}
			},
		},
		{
			name: "success: instance of struct of slice",
			args: args{
				i: &[]Struct{},
			},
			test: func(t *testing.T, r any, err error) {
				if _, ok := r.(*Struct); !ok {
					t.Errorf("result should be valid instance")
				}
				if err != nil {
					t.Errorf("error should not be nil")
				}
			},
		},
		{
			name: "success: instance of *struct of slice",
			args: args{
				i: &[]*Struct{},
			},
			test: func(t *testing.T, r any, err error) {
				if _, ok := r.(*Struct); !ok {
					t.Errorf("result should be valid instance")
				}
				if err != nil {
					t.Errorf("error should not be nil")
				}
			},
		},
		{
			name: "error: other type",
			args: args{
				i: &[]int{1},
			},
			test: func(t *testing.T, r any, err error) {
				if err == nil {
					t.Errorf("error should be nil")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStructInstance(tt.args.i)
			tt.test(t, got, err)
		})
	}
}
