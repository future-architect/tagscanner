package binarypatternmatch

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	type args struct {
		dest   interface{}
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    interface{}
	}{
		{
			name: "single bit (1)",
			args: args{
				dest: &struct {
					B bool `bits:"1"`
				}{},
				reader: bytes.NewReader([]byte{0x1}),
			},
			want: &struct {
				B bool `bits:"1"`
			}{
				B: true,
			},
		},
		{
			name: "single bit (2)",
			args: args{
				dest: &struct {
					B bool // no tag -> default is as same as type
				}{},
				reader: bytes.NewReader([]byte{0x1}),
			},
			want: &struct {
				B bool
			}{
				B: true,
			},
		},
		{
			name: "single byte (1)",
			args: args{
				dest: &struct {
					B byte `bytes:"1"`
				}{},
				reader: bytes.NewReader([]byte{0x12}),
			},
			want: &struct {
				B byte `bytes:"1"`
			}{
				B: 0x12,
			},
		},
		{
			name: "single byte (2)",
			args: args{
				dest: &struct {
					B byte // no tag -> default is as same as type
				}{},
				reader: bytes.NewReader([]byte{0x34}),
			},
			want: &struct {
				B byte // no tag -> default is as same as type
			}{
				B: 0x34,
			},
		},
		{
			name: "multiple bits (1)",
			args: args{
				dest: &struct {
					B byte `bits:"6"`
				}{},
				reader: bytes.NewReader([]byte{0xa1}),
			},
			want: &struct {
				B byte `bits:"6"`
			}{
				B: 0xa1 & 0b00111111,
			},
		},
		{
			name: "multiple bits (2)",
			args: args{
				dest: &struct {
					B1 byte `bits:"4"`
				}{},
				reader: bytes.NewReader([]byte{0xab}),
			},
			want: &struct {
				B1 byte `bits:"4"`
			}{
				B1: 0x0b,
			},
		},
		{
			name: "multiple bits (3)",
			args: args{
				dest: &struct {
					B1 byte `bits:"4"`
					B2 byte `bits:"4"`
				}{},
				reader: bytes.NewReader([]byte{0xab}),
			},
			want: &struct {
				B1 byte `bits:"4"`
				B2 byte `bits:"4"`
			}{
				B1: 0x0b,
				B2: 0x0a,
			},
		},
		{
			name: "multiple bits (4)",
			args: args{
				dest: &struct {
					B1 byte `bits:"4"`
					B2 byte `bits:"2"`
					B3 bool `bits:"1"`
					B4 bool `bits:"1"`
				}{},
				reader: bytes.NewReader([]byte{0xfb}),
			},
			want: &struct {
				B1 byte `bits:"4"`
				B2 byte `bits:"2"`
				B3 bool `bits:"1"`
				B4 bool `bits:"1"`
			}{
				B1: 0x0b,
				B2: 0x03,
				B3: true,
				B4: true,
			},
		},
		{
			name: "multiple bytes",
			args: args{
				dest: &struct {
					B int32 `bytes:"2"`
				}{},
				reader: bytes.NewReader([]byte{0x11, 0x22}),
			},
			want: &struct {
				B int32 `bytes:"2"`
			}{
				B: 0x00001122,
			},
		},
		{
			name: "omit field",
			args: args{
				dest: &struct {
					B1 byte `bits:"4"`
					B2 byte `bits:"2"`
					_  bool `bits:"1"`
					B4 bool `bits:"1"`
				}{},
				reader: bytes.NewReader([]byte{0xfb}),
			},
			want: &struct {
				B1 byte `bits:"4"`
				B2 byte `bits:"2"`
				_  bool `bits:"1"`
				B4 bool `bits:"1"`
			}{
				B1: 0x0b,
				B2: 0x03,
				B4: true,
			},
		},
		{
			name: "pattern match field(1): OK",
			args: args{
				dest: &struct {
					B1 byte `bits:"4"`
					B2 byte `bits:"2"`
					B3 bool `bits:"<<0b1>>"`
					B4 bool `bits:"1"`
				}{},
				reader: bytes.NewReader([]byte{0xfb}),
			},
			want: &struct {
				B1 byte `bits:"4"`
				B2 byte `bits:"2"`
				B3 bool `bits:"<<0b1>>"`
				B4 bool `bits:"1"`
			}{
				B1: 0x0b,
				B2: 0x03,
				B3: true,
				B4: true,
			},
		},
		{
			name: "pattern match field(1): Error",
			args: args{
				dest: &struct {
					B1 byte `bits:"4"`
					B2 byte `bits:"2"`
					B3 bool `bits:"<<0b1>>"`
					B4 bool `bits:"1"`
				}{},
				reader: bytes.NewReader([]byte{0xbb}),
			},
			wantErr: true,
		},
		{
			name: "read long value",
			args: args{
				dest: &struct {
					F uint32 `bytes:"4"`
				}{},
				reader: bytes.NewReader([]byte{0x3F, 0xE0, 0x00, 0x00}),
			},
			want: &struct {
				F uint32 `bytes:"4"`
			}{
				F: 0x3FE00000,
			},
		},
		{
			name: "read string",
			args: args{
				dest: &struct {
					Greet string `bytes:"5"`
				}{},
				reader: bytes.NewReader([]byte{0x68, 0x65, 0x6c, 0x6c, 0x6f}),
			},
			want: &struct {
				Greet string `bytes:"5"`
			}{
				Greet: "hello",
			},
		},
		{
			name: "read bytes",
			args: args{
				dest: &struct {
					Greet []byte `bytes:"5"`
				}{},
				reader: bytes.NewReader([]byte{0x68, 0x65, 0x6c, 0x6c, 0x6f}),
			},
			want: &struct {
				Greet []byte `bytes:"5"`
			}{
				Greet: []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Decode(tt.args.dest, tt.args.reader); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr && !reflect.DeepEqual(tt.args.dest, tt.want) {
				t.Errorf("Decode() got = %v, want %v", tt.args.dest, tt.want)
			}
		})
	}
}

func Test_copyWithOffset(t *testing.T) {
	type args struct {
		preload byte
		src     []byte
		offset  int
		length  int
	}
	tests := []struct {
		name           string
		args           args
		wantResult     []byte
		wantRemainByte byte
		wantNewOffset  int
	}{
		{
			name: "copy within remained byte (1)",
			args: args{
				preload: 0x02,
				src:     nil,
				offset:  1,
				length:  1,
			},
			wantResult:     []byte{0x01},
			wantRemainByte: 0x02,
			wantNewOffset:  2,
		},
		{
			name: "copy within remained byte (2)",
			args: args{
				preload: 0x20,
				src:     nil,
				offset:  4,
				length:  4,
			},
			wantResult:     []byte{0x02},
			wantRemainByte: 0x0,
			wantNewOffset:  8,
		},
		{
			name: "copy with extra bytes (no remained byte)",
			args: args{
				preload: 0x20,
				src:     []byte{0x31},
				offset:  4,
				length:  12,
			},
			wantResult:     []byte{0x12, 0x03},
			wantRemainByte: 0x0,
			wantNewOffset:  8,
		},
		{
			name: "copy with extra bytes (with remained byte)",
			args: args{
				preload: 0x20,
				src:     []byte{0x31},
				offset:  4,
				length:  8,
			},
			wantResult:     []byte{0x12},
			wantRemainByte: 0x31,
			wantNewOffset:  4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotRemainByte, gotNewOffset := copyWithOffset(tt.args.preload, tt.args.src, tt.args.offset, tt.args.length)
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("copyWithOffset() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if gotRemainByte != tt.wantRemainByte {
				t.Errorf("copyWithOffset() gotRemainByte = %v, want %v", gotRemainByte, tt.wantRemainByte)
			}
			if gotNewOffset != tt.wantNewOffset {
				t.Errorf("copyWithOffset() gotNewOffset = %v, want %v", gotNewOffset, tt.wantNewOffset)
			}
		})
	}
}
