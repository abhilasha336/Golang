package audio

import (
	"io"
	"time"

	"github.com/go-flac/go-flac"
	"gitlab.com/tuneverse/toolkit/utils"
)

type FlacDecoder struct {
	UnimplementedDecoder
	DecoderConfig
	duration, sampleRate float64
}

// NewFlacDecoder
func NewFlacDecoder(config DecoderConfig) Decoder {
	return &FlacDecoder{
		DecoderConfig: config,
	}
}

// WithBytes
func (decoder *FlacDecoder) WithBytes(bytes io.ReadSeeker) Decoder {
	decoder.bytes = bytes
	return decoder
}

// Bytes
func (decoder *FlacDecoder) Bytes() io.ReadSeeker {
	return decoder.bytes
}

// Close
func (decoder *FlacDecoder) Close() error {
	if decoder.closer == nil {
		return nil
	}
	return decoder.closer.Close()
}

// decoder
func (decoder *FlacDecoder) Decode() (Decoder, error) {
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
	file, err := flac.ParseBytes(readSeeker)
	if err != nil {
		return nil, err
	}

	// Get the stream info
	stream, err := file.GetStreamInfo()
	if err != nil {
		return nil, err
	}

	decoder.duration = time.Duration((float64(stream.SampleCount) / float64(stream.SampleRate)) * float64(time.Second)).Seconds()
	decoder.sampleRate = float64(stream.SampleRate)

	return decoder, nil
}

// Duration
func (decoder *FlacDecoder) Duration() (float64, error) {
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
func (decoder *FlacDecoder) SampleRate() (float64, error) {
	var err error
	// duration calculation with point
	_, sampleRate, err := utils.FormatWithDecimalLimit(decoder.sampleRate, decoder.DecimalLimit)
	if err != nil {
		return sampleRate, err
	}

	return sampleRate, nil
}
