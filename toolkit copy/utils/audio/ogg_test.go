package audio_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/tuneverse/toolkit/utils/audio"
)

func TestOGGDecoder_Duration(t *testing.T) {

	t.Run("local file", func(t *testing.T) {
		config := audio.DecoderConfig{
			FilePath:     "samples/audio.ogg",
			DecimalLimit: "3",
		}

		decoder, err := audio.NewOGGDecoder(config).Decode()
		require.Nil(t, err)

		duration, err := decoder.Duration()
		require.Nil(t, err)

		require.NotNil(t, duration)
		fmt.Println("duration::", duration)

	})

	// t.Run("remote file", func(t *testing.T) {
	// 	config := audio.DecoderConfig{
	// 		FilePath:     "https://demo-bucket-sample.s3.amazonaws.com/audio.aiff",
	// 		DecimalLimit: "3",
	// 		RemoteFile:   true,
	// 	}

	// 	decoder := audio.NewAIFFDecoder(config)

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

		decoder := audio.NewOGGDecoder(config)
		f, err := os.Open("samples/audio.ogg")
		require.Nil(t, err)

		defer f.Close()

		decode, _ := decoder.WithBytes(f).Decode()
		duration, err := decode.Duration()
		require.Nil(t, err)
		require.NotNil(t, duration)

		t.Log("custom duration::", duration)

	})
}
