package audio

import (
	"io"

	"github.com/gopxl/beep/wav"
	"gitlab.com/tuneverse/toolkit/utils"
)

type WavDecoder struct {
	UnimplementedDecoder
	DecoderConfig
	duration, sampleRate float64
}

// NewWavDecoder
func NewWavDecoder(config DecoderConfig) Decoder {
	return &WavDecoder{
		DecoderConfig: config,
	}
}

// WithBytes
func (decoder *WavDecoder) WithBytes(bytes io.ReadSeeker) Decoder {
	decoder.bytes = bytes
	return decoder
}

// Bytes
func (decoder *WavDecoder) Bytes() io.ReadSeeker {
	return decoder.bytes
}

// Close
func (decoder *WavDecoder) Close() error {
	if decoder.closer == nil {
		return nil
	}
	return decoder.closer.Close()
}

func (decoder *WavDecoder) Decode() (Decoder, error) {
	var readSeeker io.ReadSeeker

	// Check if bytes were passed from the client side; if not, fetch them
	if readSeeker = decoder.Bytes(); readSeeker == nil {
		r, closer, err := LoadFile(decoder.FilePath, decoder.RemoteFile)
		if err != nil {
			return nil, err
		}
		readSeeker = r
		defer closer.Close()
	}

	decodedseeker, format, err := wav.Decode(readSeeker)
	if err != nil {
		return nil, err
	}
	duration := decodedseeker.Len() / int(format.SampleRate)

	decoder.sampleRate = float64(format.SampleRate)
	decoder.duration = float64(duration)

	return decoder, nil
}

// Duration
func (decoder *WavDecoder) Duration() (float64, error) {
	var duration float64
	duration = decoder.duration

	// duration calculation with point
	_, duration, err := utils.FormatWithDecimalLimit(duration, decoder.DecimalLimit)
	if err != nil {
		return duration, err
	}

	return duration, nil
}

// SampleRate
func (decoder *WavDecoder) SampleRate() (float64, error) {
	var err error
	// duration calculation with point
	_, sampleRate, err := utils.FormatWithDecimalLimit(decoder.sampleRate, decoder.DecimalLimit)
	if err != nil {
		return sampleRate, err
	}

	return sampleRate, nil
}
