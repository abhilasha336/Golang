package audio_test

import (
	"fmt"
	"os"
	"testing"

	"gitlab.com/tuneverse/toolkit/utils/audio"
)

func TestMp3Decoder_Duration(t *testing.T) {

	// t.Run("remot file", func(t *testing.T) {
	// 	config := audio.DecoderConfig{
	// 		FilePath:   "https://www.soundhelix.com/examples/mp3/SoundHelix-Song-1.mp3",
	// 		RemoteFile: true,
	// 	}

	// 	decoder := audio.NewMP3Decoder(config)

	// 	duration, err := decoder.Duration()
	// 	if err != nil {
	// 		t.Errorf("Error calculating duration: %v", err)
	// 	}
	// 	if duration <= 0 {
	// 		t.Errorf("Invalid duration: %v", duration)
	// 	}

	// 	fmt.Println(duration)
	// })

	t.Run("local file", func(t *testing.T) {
		config := audio.DecoderConfig{
			FilePath:     "samples/audio.mp3",
			DecimalLimit: "3",
		}

		decoder, err := audio.NewMP3Decoder(config).Decode()
		if err != nil {
			t.Errorf("Error on deciode: %v", err)
		}

		_ = decoder
		duration, err := decoder.Duration()
		if err != nil {
			t.Errorf("Error calculating duration: %v", err)
		}

		if duration <= 0 {
			t.Errorf("Invalid duration: %v", duration)
		}

		fmt.Println(duration)
	})

	t.Run("with custom bytes array", func(t *testing.T) {
		config := audio.DecoderConfig{
			DecimalLimit: "3",
		}

		decoder := audio.NewMP3Decoder(config)
		f, _ := os.Open("samples/audio.mp3")
		defer f.Close()

		decode, _ := decoder.WithBytes(f).Decode()
		duration, err := decode.Duration()
		if err != nil {
			t.Errorf("Error calculating duration: %v", err)
		}
		if duration <= 0 {
			t.Errorf("Invalid duration: %v", duration)
		}

		t.Log("custom duration::", duration)

	})
}
