package audio

import (
	"io"

	"github.com/go-audio/aiff"
	"gitlab.com/tuneverse/toolkit/utils"
)

type AiffDecoder struct {
	UnimplementedDecoder
	DecoderConfig
	duration, sampleRate float64
}

// NewAIFFDecoder
func NewAIFFDecoder(config DecoderConfig) Decoder {
	return &AiffDecoder{
		DecoderConfig: config,
	}
}

// WithBytes
func (decoder *AiffDecoder) WithBytes(bytes io.ReadSeeker) Decoder {
	decoder.bytes = bytes
	return decoder
}

// Bytes
func (decoder *AiffDecoder) Bytes() io.ReadSeeker {
	return decoder.bytes
}

// Close
func (decoder *AiffDecoder) Close() error {
	if decoder.closer == nil {
		return nil
	}
	return decoder.closer.Close()
}

// decoder
func (decoder *AiffDecoder) Decode() (Decoder, error) {
	var readSeeker io.ReadSeeker

	// Check the bytes was passed from the client side
	// if it is not passed from the client side, then fetch it
	if readSeeker = decoder.Bytes(); readSeeker == nil {
		// load the file
		r, closer, err := LoadFile(decoder.FilePath, decoder.RemoteFile)
		if err != nil {
			return nil, err
		}
		readSeeker = r
		defer closer.Close()
	}

	d := aiff.NewDecoder(readSeeker)

	timeDuration, err := d.Duration()
	if err != nil {
		return nil, err
	}

	decoder.duration = timeDuration.Seconds()
	decoder.sampleRate = float64(d.SampleRate)

	return decoder, nil
}

// Duration
func (decoder *AiffDecoder) Duration() (float64, error) {
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
func (decoder *AiffDecoder) SampleRate() (float64, error) {
	var err error
	// duration calculation with point
	_, sampleRate, err := utils.FormatWithDecimalLimit(decoder.sampleRate, decoder.DecimalLimit)
	if err != nil {
		return sampleRate, err
	}

	return sampleRate, nil
}
