package audio

import (
	"io"
	"log"

	"github.com/tcolgate/mp3"
	"gitlab.com/tuneverse/toolkit/utils"
)

type Mp3Decoder struct {
	UnimplementedDecoder
	DecoderConfig
	duration, sampleRate float64
}

// NewMP3Decoder
func NewMP3Decoder(config DecoderConfig) Decoder {
	return &Mp3Decoder{
		DecoderConfig: config,
	}
}

// WithBytes
func (decoder *Mp3Decoder) WithBytes(bytes io.ReadSeeker) Decoder {
	decoder.bytes = bytes
	return decoder
}

// Bytes
func (decoder *Mp3Decoder) Bytes() io.ReadSeeker {
	return decoder.bytes
}

// Close
func (decoder *Mp3Decoder) Close() error {
	if decoder.closer == nil {
		return nil
	}
	return decoder.closer.Close()
}

// decoder
func (decoder *Mp3Decoder) Decode() (Decoder, error) {
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
		decoder.closer = closer
	}

	// decode
	var skipped int
	var duration float64
	var totalSamples int

	d := mp3.NewDecoder(readSeeker)
	var frame mp3.Frame
	for {

		if err := d.Decode(&frame, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			log.Println("an error happened when decode the audio file")
		}

		duration = duration + frame.Duration().Seconds()
		// Accumulate the total samples
		totalSamples += int(frame.Samples())
	}

	// set duration
	decoder.duration = duration
	// set sample rate
	decoder.sampleRate = float64(totalSamples) / duration

	return decoder, nil
}

// Duration
func (decoder *Mp3Decoder) Duration() (float64, error) {
	var err error
	// duration calculation with point
	_, duration, err := utils.FormatWithDecimalLimit(decoder.duration, decoder.DecimalLimit)
	if err != nil {
		return duration, err
	}

	return duration, nil
}

// SampleRate
func (decoder *Mp3Decoder) SampleRate() (float64, error) {
	var err error
	// duration calculation with point
	_, sampleRate, err := utils.FormatWithDecimalLimit(decoder.sampleRate, decoder.DecimalLimit)
	if err != nil {
		return sampleRate, err
	}

	return sampleRate, nil
}
