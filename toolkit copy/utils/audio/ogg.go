package audio

import (
	"io"

	"github.com/jfreymuth/oggvorbis"
	"gitlab.com/tuneverse/toolkit/utils"
)

type OGGDecoder struct {
	UnimplementedDecoder
	DecoderConfig
	duration, sampleRate float64
}

// NewFlacDecoder
func NewOGGDecoder(config DecoderConfig) Decoder {
	return &OGGDecoder{
		DecoderConfig: config,
	}
}

// WithBytes
func (decoder *OGGDecoder) WithBytes(bytes io.ReadSeeker) Decoder {
	decoder.bytes = bytes
	return decoder
}

// Bytes
func (decoder *OGGDecoder) Bytes() io.ReadSeeker {
	return decoder.bytes
}

// Close
func (decoder *OGGDecoder) Close() error {
	if decoder.closer == nil {
		return nil
	}
	return decoder.closer.Close()
}

// decoder
func (decoder *OGGDecoder) Decode() (Decoder, error) {
	var readSeeker io.ReadSeeker

	// Check if the bytes were passed from the client side
	if readSeeker = decoder.Bytes(); readSeeker == nil {
		// Load the file
		r, closer, err := LoadFile(decoder.FilePath, decoder.RemoteFile)
		if err != nil {
			return nil, err
		}
		readSeeker = r
		defer closer.Close()
	}

	// Parse OGG Vorbis data
	oggReader, err := oggvorbis.NewReader(readSeeker)
	if err != nil {
		return nil, err
	}

	// Calculate duration
	duration := float64(oggReader.Length()) / float64(oggReader.SampleRate())

	// Set the duration and sample rate in the decoder
	decoder.duration = duration
	decoder.sampleRate = float64(oggReader.SampleRate())

	return decoder, nil
}

// Duration
func (decoder *OGGDecoder) Duration() (float64, error) {
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
func (decoder *OGGDecoder) SampleRate() (float64, error) {
	var err error
	// duration calculation with point
	_, sampleRate, err := utils.FormatWithDecimalLimit(decoder.sampleRate, decoder.DecimalLimit)
	if err != nil {
		return sampleRate, err
	}

	return sampleRate, nil
}
