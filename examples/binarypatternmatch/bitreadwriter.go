package binarypatternmatch

import (
	"fmt"
	"io"
)

func maskBits(i byte, l int) byte {
	switch l {
	case 1:
		return i & 0b00000001
	case 2:
		return i & 0b00000011
	case 3:
		return i & 0b00000111
	case 4:
		return i & 0b00001111
	case 5:
		return i & 0b00011111
	case 6:
		return i & 0b00111111
	case 7:
		return i & 0b01111111
	case 8:
		return i
	default:
		panic(fmt.Sprintf("maskBits only accepts 0 < l && l < 8, but %d", l))
	}
}

func copyWithOffset(preload byte, src []byte, offset, bits int) (result []byte, remainByte byte, newOffset int) {
	byl := bits / 8
	if bits%8 != 0 {
		byl += 1
	}
	result = make([]byte, 0, byl)
	remainByte = preload
	if len(src) == 0 {
		result = append(result, maskBits(preload>>offset, bits))
		newOffset = offset + bits
		if newOffset == 8 {
			remainByte = 0
		}
		return
	}
	for _, b := range src {
		if bits < 8 {
			result = append(result, maskBits(b>>offset, bits)<<offset|remainByte)
			remainByte = b
			newOffset = offset + bits
			bits = 0
		} else {
			result = append(result, remainByte>>offset|maskBits(b, offset)<<(8-offset))
			newOffset = offset
			bits -= 8
			if newOffset == 8 {
				remainByte = 0
			} else {
				remainByte = b
			}
		}
	}
	if bits > 0 {
		result = append(result, maskBits(remainByte>>offset, bits))
		newOffset = offset + bits
		if newOffset == 8 {
			remainByte = 0
		}
	}
	return
}

type bitWriter struct {
	writer       io.Writer
	bitCap       int
	remainedByte byte
}

func newBitWriter(writer io.Writer) *bitWriter {
	return &bitWriter{
		writer: writer,
		bitCap: 8,
	}
}

func (b *bitWriter) Padding(bit int) {
	if bit == 0 {
		b.writer.Write([]byte{b.remainedByte})
		b.bitCap = 8
	} else {
		if b.bitCap <= bit {
			b.writer.Write([]byte{b.remainedByte})
			bit -= b.bitCap
			b.bitCap = 8
			for bit > 8 {
				b.writer.Write([]byte{0})
				bit -= 8
			}
		}
		b.bitCap -= bit
	}
}

func (b *bitWriter) WriteBool(data bool) {
	if data {
		if b.bitCap == 1 {
			b.writer.Write([]byte{(1 << 7) | b.remainedByte})
			b.bitCap = 8
		} else {
			b.remainedByte = b.remainedByte | 1<<(8-b.bitCap)
			b.bitCap -= 1
		}
	} else {
		if b.bitCap == 1 {
			b.writer.Write([]byte{b.remainedByte})
			b.bitCap = 8
		} else {
			b.bitCap -= 1
		}
	}
}

func (b *bitWriter) WriteUint(data uint64, bit int) {

}

func (b *bitWriter) WriteByte(data byte, bit int) {
	if b.bitCap == 8 && bit == 8 {
		b.writer.Write([]byte{data})
		return
	}
	if b.bitCap == 8 {
		b.remainedByte = maskBits(data, bit)
		b.bitCap -= bit
	} else if b.bitCap+bit <= 8 {
		b.remainedByte = b.remainedByte | maskBits(data, bit)<<(8-b.bitCap)
		b.bitCap -= bit
		if b.bitCap == 0 {
			b.bitCap = 8
			b.writer.Write([]byte{b.remainedByte})
			b.remainedByte = 0
		}
	} else {
		b.writer.Write([]byte{b.remainedByte | maskBits(data, b.bitCap)<<(8-b.bitCap)})
		b.remainedByte = data >> b.bitCap
		b.bitCap = b.bitCap + bit - 8
	}
}

func (b *bitWriter) WriteBytes(data []byte, bit int) {
	if b.bitCap == 8 && bit == len(data)*8 {
		b.writer.Write(data)
		return
	}
	if b.bitCap == 8 {
		bytes := bit / 8
		b.writer.Write(data[:bytes])
		b.WriteByte(data[len(data)-1], bit%8)
	}
}
