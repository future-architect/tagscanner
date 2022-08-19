package binarypatternmatch

import (
	"bytes"
	"reflect"
	"testing"
)

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

func TestBitWriter(t *testing.T) {
	tests := []struct {
		name  string
		setup func(w *bitWriter)
		want  []byte
	}{
		{
			name: "single byte",
			setup: func(w *bitWriter) {
				w.WriteByte(0x12, 8)
			},
			want: []byte{0x12},
		},
		{
			name: "half byte",
			setup: func(w *bitWriter) {
				w.WriteByte(0x02, 4)
				w.WriteByte(0x01, 4)
			},
			want: []byte{0x12},
		},
		{
			name: "over byte align",
			setup: func(w *bitWriter) {
				w.WriteByte(0x01, 4)
				w.WriteByte(0x32, 8)
				w.WriteByte(0x04, 4)
			},
			want: []byte{0x21, 0x43},
		},
		{
			name: "write bytes",
			setup: func(w *bitWriter) {
				w.WriteBytes([]byte{0x12, 0x34}, 16)
			},
			want: []byte{0x12, 0x34},
		},
		{
			name: "write bytes (2)",
			setup: func(w *bitWriter) {
				w.WriteBytes([]byte{0x12, 0x34}, 12)
				w.WriteByte(0x01, 4)
			},
			want: []byte{0x12, 0x14},
		},
		{
			name: "write bool (1)",
			setup: func(w *bitWriter) {
				w.WriteBool(true)
				w.Padding(0)
			},
			want: []byte{0x01},
		},
		{
			name: "write bool (2)",
			setup: func(w *bitWriter) {
				w.WriteBool(false)
				w.Padding(0)
			},
			want: []byte{0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			bw := newBitWriter(&buf)
			tt.setup(bw)
			if !reflect.DeepEqual(tt.want, buf.Bytes()) {
				t.Errorf("%v != %v", tt.want, buf.Bytes())
			}
		})
	}
}
