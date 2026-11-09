package opus

/*
#cgo pkg-config: opus
#include <opus.h>
*/
import "C"
import "fmt"

type Decoder struct {
	dec     *C.OpusDecoder
	channel int
}

func (d *Decoder) Decode(data []byte, buf []int16) ([]int16, error) {
	n := C.opus_decode(
		d.dec,
		(*C.uchar)(&data[0]),
		C.opus_int32(len(data)),
		(*C.opus_int16)(&buf[0]),
		C.int(len(buf)/d.channel),
		0,
	)
	if n < 0 {
		return nil, fmt.Errorf("opus_decode failed: %d", n)
	}
	return buf[:int(n)*d.channel], nil
}

func (d *Decoder) DecodeFloat32(data []byte, buf []float32) ([]float32, error) {
	n := C.opus_decode_float(
		d.dec,
		(*C.uchar)(&data[0]),
		C.opus_int32(len(data)),
		(*C.float)(&buf[0]),
		C.int(len(buf)/d.channel),
		0,
	)
	if n < 0 {
		return nil, fmt.Errorf("opus_decode failed: %d", n)
	}
	return buf[:int(n)*d.channel], nil
}

func (d *Decoder) Close() {
	if d.dec != nil {
		C.opus_decoder_destroy(d.dec)
		d.dec = nil
	}
}

func NewDecoder(sampleRate, channels int) (*Decoder, error) {
	var errCode C.int
	dec := C.opus_decoder_create(
		C.opus_int32(sampleRate),
		C.int(channels),
		&errCode,
	)
	if errCode != C.OPUS_OK {
		return nil, fmt.Errorf("opus_decoder_create failed: %d", errCode)
	}
	return &Decoder{dec, channels}, nil
}
