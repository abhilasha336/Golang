package audio_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/tuneverse/toolkit/utils/audio"
)

func TestFlacDecoder_Duration(t *testing.T) {

	t.Run("local file", func(t *testing.T) {
		config := audio.DecoderConfig{
			FilePath:     "samples/audio.flac",
			DecimalLimit: "3",
		}

		decoder, err := audio.NewFlacDecoder(config).Decode()
		require.Nil(t, err)

		duration, err := decoder.Duration()
		require.Nil(t, err)

		require.NotNil(t, duration)

		fmt.Println("duration::", duration)

	})

	// t.Run("remote file", func(t *testing.T) {
	// 	config := audio.DecoderConfig{
	// 		FilePath:     "https://filesamples.com/samples/audio/flac/Symphony%20No.6%20(1st%20movement).flac",
	// 		DecimalLimit: "3",
	// 		RemoteFile:   true,
	// 	}

	// 	decoder := audio.NewFlacDecoder(config)

	// 	duration, err := decoder.Duration()
	// 	if err != nil {
	// 		t.Errorf("Error calculating duration: %v", err)
	// 	}
	// 	if duration <= 0 {
	// 		t.Errorf("Invalid duration: %v", duration)
	// 	}

	// 	fmt.Println("duration::", duration)
	// })

	t.Run("with custom bytes array", func(t *testing.T) {
		config := audio.DecoderConfig{
			DecimalLimit: "3",
		}

		decoder := audio.NewFlacDecoder(config)
		f, err := os.Open("samples/audio.flac")
		require.Nil(t, err)

		defer f.Close()

		decode, err := decoder.WithBytes(f).Decode()
		require.Nil(t, err)

		duration, err := decode.Duration()
		require.Nil(t, err)

		require.NotNil(t, duration)

		t.Log("custom duration::", duration)

	})
}
