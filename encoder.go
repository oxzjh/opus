package opus

/*
#cgo pkg-config: opus
#include <opus.h>
*/
import "C"
import "fmt"

type Encoder struct {
	enc      *C.OpusEncoder
	channels int
}

func (e *Encoder) Encode(pcm []int16, buf []byte) ([]byte, error) {
	n := C.opus_encode(
		e.enc,
		(*C.opus_int16)(&pcm[0]),
		C.int(len(pcm)/e.channels),
		(*C.uchar)(&buf[0]),
		C.opus_int32(len(buf)),
	)
	if n < 0 {
		return nil, fmt.Errorf("opus_encode failed: %d", n)
	}
	return buf[:n], nil
}

func (e *Encoder) EncodeFloat32(pcm []float32, buf []byte) ([]byte, error) {
	n := C.opus_encode_float(
		e.enc,
		(*C.float)(&pcm[0]),
		C.int(len(pcm)/e.channels),
		(*C.uchar)(&buf[0]),
		C.opus_int32(len(buf)),
	)
	if n < 0 {
		return nil, fmt.Errorf("opus_encode failed: %d", n)
	}
	return buf[:n], nil
}

func (e *Encoder) Close() {
	if e.enc != nil {
		C.opus_encoder_destroy(e.enc)
		e.enc = nil
	}
}

func NewEncoder(sampleRate, channels int) (*Encoder, error) {
	var errCode C.int
	enc := C.opus_encoder_create(
		C.opus_int32(sampleRate),
		C.int(channels),
		C.OPUS_APPLICATION_VOIP,
		&errCode,
	)
	if errCode != C.OPUS_OK {
		return nil, fmt.Errorf("opus_encoder_create failed: %d", errCode)
	}
	return &Encoder{enc, channels}, nil
}
