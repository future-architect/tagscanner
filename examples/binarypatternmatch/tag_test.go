package binarypatternmatch

import (
	"reflect"
	"testing"
)

func TestParseBinaryTag_Bit(t *testing.T) {
	type args struct {
		fieldName string
		tagKey    string
		tagValue  string
		eType     reflect.Type
	}
	tests := []struct {
		name    string
		args    args
		want    *binaryTag
		wantErr bool
	}{
		{
			name: "bits with length: ok",
			args: args{
				fieldName: "Bits",
				tagKey:    "bits",
				tagValue:  "8",
				eType:     reflect.TypeOf(byte(0)),
			},
			want: &binaryTag{
				Type:    Bits,
				EType:   reflect.TypeOf(byte(0)),
				Literal: []byte{},
				Bits:    8,
				Label:   "",
			},
			wantErr: false,
		},
		{
			name: "bits with length: ok(2)",
			args: args{
				fieldName: "Bits",
				tagKey:    "bits",
				tagValue:  "1",
				eType:     reflect.TypeOf(true),
			},
			want: &binaryTag{
				Type:    Bits,
				EType:   reflect.TypeOf(true),
				Literal: []byte{},
				Bits:    1,
				Label:   "",
			},
			wantErr: false,
		},
		{
			name: "bits with length: ok (3): omit tags",
			args: args{
				fieldName: "Bits",
				tagKey:    "",
				tagValue:  "",
				eType:     reflect.TypeOf(byte(1)),
			},
			want: &binaryTag{
				Type:    Bits,
				EType:   reflect.TypeOf(byte(1)),
				Literal: []byte{},
				Bits:    8,
				Label:   "",
			},
			wantErr: false,
		},
		{
			name: "bits with length: ok (3): hex literal",
			args: args{
				fieldName: "Bits",
				tagKey:    "bits",
				tagValue:  "<<0x10>>",
				eType:     reflect.TypeOf(byte(1)),
			},
			want: &binaryTag{
				Type:    Bits,
				EType:   reflect.TypeOf(byte(1)),
				Literal: []byte("\x10"),
				Bits:    8,
				Label:   "",
			},
			wantErr: false,
		},
		{
			name: "bits with length: ok (3): hex literal with length",
			args: args{
				fieldName: "Bits",
				tagKey:    "bits",
				tagValue:  "<<0x10/6>>",
				eType:     reflect.TypeOf(byte(1)),
			},
			want: &binaryTag{
				Type:    Bits,
				EType:   reflect.TypeOf(byte(1)),
				Literal: []byte("\x10"),
				Bits:    6,
				Label:   "",
			},
			wantErr: false,
		},
		{
			name: "bits with length: ok (4): binary literal with length",
			args: args{
				fieldName: "Bits",
				tagKey:    "bits",
				tagValue:  "<<0b110000/6>>",
				eType:     reflect.TypeOf(byte(1)),
			},
			want: &binaryTag{
				Type:    Bits,
				EType:   reflect.TypeOf(byte(1)),
				Literal: []byte("\x30"),
				Bits:    6,
				Label:   "",
			},
			wantErr: false,
		},
		{
			name: "bits with length: ok (5): byte literal",
			args: args{
				fieldName: "Bits",
				tagKey:    "bits",
				tagValue:  "<<head/32>>",
				eType:     reflect.TypeOf(int32(1)),
			},
			want: &binaryTag{
				Type:    Bits,
				EType:   reflect.TypeOf(int32(1)),
				Literal: []byte("head"),
				Bits:    32,
				Label:   "",
			},
			wantErr: false,
		},
		{
			name: "bits with length: ng(1)",
			args: args{
				fieldName: "Bits",
				tagKey:    "bits",
				tagValue:  "2",
				eType:     reflect.TypeOf(true),
			},
			wantErr: true,
		},
		{
			name: "bits with length: ng(2)",
			args: args{
				fieldName: "Bits",
				tagKey:    "bits",
				tagValue:  "9",
				eType:     reflect.TypeOf(byte(1)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseBinaryTag(tt.args.fieldName, tt.args.tagKey, tt.args.tagValue, tt.args.fieldName, tt.args.eType)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseBinaryTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseBinaryTag() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseBinaryTag_Bytes(t *testing.T) {
	type args struct {
		fieldName string
		tagKey    string
		tagValue  string
		eType     reflect.Type
	}
	tests := []struct {
		name    string
		args    args
		want    *binaryTag
		wantErr bool
	}{
		{
			name: "bytes with length: ok",
			args: args{
				fieldName: "Bytes",
				tagKey:    "bytes",
				tagValue:  "1",
				eType:     reflect.TypeOf(byte(0)),
			},
			want: &binaryTag{
				Type:    Bytes,
				EType:   reflect.TypeOf(byte(0)),
				Literal: []byte{},
				Bits:    8,
				Label:   "",
			},
			wantErr: false,
		},
		{
			name: "bytes with length: ok(2)",
			args: args{
				fieldName: "Bytes",
				tagKey:    "bytes",
				tagValue:  "3",
				eType:     reflect.TypeOf(int32(1)),
			},
			want: &binaryTag{
				Type:    Bytes,
				EType:   reflect.TypeOf(int32(1)),
				Literal: []byte{},
				Bits:    24,
				Label:   "",
			},
			wantErr: false,
		},
		{
			name: "bytes with length: ok (3): hex literal",
			args: args{
				fieldName: "Bytes",
				tagKey:    "bytes",
				tagValue:  "<<0x10>>",
				eType:     reflect.TypeOf(int64(1)),
			},
			want: &binaryTag{
				Type:    Bytes,
				EType:   reflect.TypeOf(int64(1)),
				Literal: []byte("\x00\x00\x00\x00\x00\x00\x00\x10"),
				Bits:    64,
				Label:   "",
			},
			wantErr: false,
		},
		{
			name: "bytes with length: ok (4): hex literal with length",
			args: args{
				fieldName: "Bytes",
				tagKey:    "bytes",
				tagValue:  "<<0x10/6>>",
				eType:     reflect.TypeOf(int64(1)),
			},
			want: &binaryTag{
				Type:    Bytes,
				EType:   reflect.TypeOf(int64(1)),
				Literal: []byte("\x00\x00\x00\x00\x00\x10"),
				Bits:    48,
				Label:   "",
			},
			wantErr: false,
		},
		{
			name: "bytes with length: ok (5): byte literal",
			args: args{
				fieldName: "Bytes",
				tagKey:    "bytes",
				tagValue:  "<<head/4>>",
				eType:     reflect.TypeOf(int32(1)),
			},
			want: &binaryTag{
				Type:    Bytes,
				EType:   reflect.TypeOf(int32(1)),
				Literal: []byte("head"),
				Bits:    32,
				Label:   "",
			},
			wantErr: false,
		},
		{
			name: "bytes with length: ng(1)",
			args: args{
				fieldName: "Bytes",
				tagKey:    "bytes",
				tagValue:  "3",
				eType:     reflect.TypeOf(uint16(1)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseBinaryTag(tt.args.fieldName, tt.args.tagKey, tt.args.tagValue, tt.args.fieldName, tt.args.eType)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseBinaryTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseBinaryTag() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseBinaryTag_Length(t *testing.T) {
	type args struct {
		fieldName string
		tagKey    string
		tagValue  string
		eType     reflect.Type
	}
	tests := []struct {
		name    string
		args    args
		want    *binaryTag
		wantErr bool
	}{
		{
			name: "bytes with length: ok",
			args: args{
				fieldName: "Length",
				tagKey:    "length",
				tagValue:  "@length",
				eType:     reflect.TypeOf(int32(0)),
			},
			want: &binaryTag{
				Type:  Length,
				EType: reflect.TypeOf(int32(0)),
				Bits:  32,
				Label: "@length",
			},
			wantErr: false,
		},
		{
			name: "bytes with length: ok(2)",
			args: args{
				fieldName: "Length",
				tagKey:    "length",
				tagValue:  "@length/24",
				eType:     reflect.TypeOf(int32(0)),
			},
			want: &binaryTag{
				Type:  Length,
				EType: reflect.TypeOf(int32(0)),
				Bits:  24,
				Label: "@length",
			},
			wantErr: false,
		},
		{
			name: "length: ng",
			args: args{
				fieldName: "Length",
				tagKey:    "length",
				tagValue:  "--",
				eType:     reflect.TypeOf(uint16(1)),
			},
			wantErr: true,
		},
		{
			name: "bits with external length: ok",
			args: args{
				fieldName: "Bits",
				tagKey:    "bits",
				tagValue:  "@length",
				eType:     reflect.TypeOf(int32(1)),
			},
			want: &binaryTag{
				Type:  Bits,
				EType: reflect.TypeOf(int32(1)),
				Label: "@length",
			},
			wantErr: false,
		},
		{
			name: "bytes with external length: ok/slice",
			args: args{
				fieldName: "Bytes",
				tagKey:    "bytes",
				tagValue:  "@length",
				eType:     reflect.TypeOf([]byte{}),
			},
			want: &binaryTag{
				Type:  Bytes,
				EType: reflect.TypeOf([]byte{}),
				Label: "@length",
			},
			wantErr: false,
		},
		{
			name: "bytes with external length: ok/string",
			args: args{
				fieldName: "Bytes",
				tagKey:    "bytes",
				tagValue:  "@length",
				eType:     reflect.TypeOf(""),
			},
			want: &binaryTag{
				Type:  Bytes,
				EType: reflect.TypeOf(""),
				Label: "@length",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseBinaryTag(tt.args.fieldName, tt.args.tagKey, tt.args.tagValue, tt.args.fieldName, tt.args.eType)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseBinaryTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseBinaryTag() got = %v, want %v", got, tt.want)
			}
		})
	}
}
