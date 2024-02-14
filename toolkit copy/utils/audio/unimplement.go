package audio

import (
	"errors"
	"io"
)

type UnimplementedDecoder struct{}
type UnimplementedSeeker struct {
}

func NewUnImplDecoder() Decoder {
	return &UnimplementedDecoder{}
}

// WithBytes
func (decoder *UnimplementedDecoder) Decode() (Decoder, error) {
	return nil, errors.New("Decode not implemented")
}

// Duration
func (decoder *UnimplementedDecoder) Duration() (float64, error) {
	return 0, errors.New("Duration not implemented")
}

// Samplerate
func (decoder *UnimplementedDecoder) SampleRate() (float64, error) {
	return 0, errors.New("SampleRate not implemented")
}

// WithBytes
func (decoder *UnimplementedDecoder) WithBytes(io.ReadSeeker) Decoder {
	return decoder
}

// WithBytes
func (decoder *UnimplementedDecoder) Bytes() io.ReadSeeker {
	return nil
}

// Close
func (decoder *UnimplementedDecoder) Close() error {
	return nil
}
