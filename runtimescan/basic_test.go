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
		want    interface{}
		wantErr bool
	}{
		{
			name: "int",
			args: args{
				value: "12345",
				elemType: reflect.TypeOf(12345),
			},
			want: int64(12345),
			wantErr: false,
		},
		{
			name: "uint",
			args: args{
				value: "12345",
				elemType: reflect.TypeOf(uint(12345)),
			},
			want: uint64(12345),
			wantErr: false,
		},
		{
			name: "float",
			args: args{
				value: "12345.6",
				elemType: reflect.TypeOf(float64(12345)),
			},
			want: 12345.6,
			wantErr: false,
		},
		{
			name: "bool",
			args: args{
				value: "true",
				elemType: reflect.TypeOf(true),
			},
			want: true,
			wantErr: false,
		},
		{
			name: "string",
			args: args{
				value: "hello world",
				elemType: reflect.TypeOf("hello"),
			},
			want: "hello world",
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
