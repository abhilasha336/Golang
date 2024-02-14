## Audio Decoder
The `audio` package provides a set of interfaces and implementations for decoding `audio` files in various formats. The package includes support for `MP3`, `WAV`, `FLAC`, and `AIFF` file formats.


## AudioFileFormat
The `AudioFileFormat` type is an enumeration of supported audio file formats. The supported formats are MP3, WAV, FLAC, and AIFF.
```go
    const (
        MP3  AudioFileFormat = "mp3"
        Wav  AudioFileFormat = "wav"
        Flac AudioFileFormat = "flac"
        AIFF AudioFileFormat = "aiff"
    )
```

## Decoder
The `Decoder` interface defines methods for decoding audio files. The interface includes methods for getting the duration of an audio file and formatting the duration as a string.
```go
    type Decoder interface {
        Decode() (Decoder, error)
        WithBytes(io.ReadSeeker) Decoder
        Bytes() io.ReadSeeker
        Close() error

        Duration() (float64, error)
        SampleRate() (float64, error)
    }
```

- **Duration**
- **SampleRate* : get the duration
- **WithBytes** : the `bytes` field is set to the input `io.ReadSeeker` object .This allows the `Decoder` to read from the input `io.ReadSeeker` object when decoding a `audio` file.
- **Bytes** : Used to set the `bytes` to the decoder
## LoadFile
The `LoadFile` function is a Go function that loads a file from either a local `file system` or `a remote location` based on the `remoteFile` config.
**use cases** 
- Load a configuration file from a remote location during application startup.
- Load a data file from the local file system during application runtime.
- Load a file from a remote location and process its contents in memory without saving it to disk.

## DecoderConfig
The `DecoderConfig` struct holds configuration options for the audio decoder. The struct includes options for specifying the decimal limit and file path.

```go
    type DecoderConfig struct {
        DecimalLimit string
        FilePath     string
    }
```

## NewDecoder
The `NewDecoder` function creates a new audio decoder based on the specified file format and configuration options. The function returns a Decoder interface that can be used to decode audio files.
```go
func NewDecoder(format AudioFileFormat, config DecoderConfig) Decoder {
	switch format {
        case MP3:
            return NewMP3Decoder(config)
        case Wav:
            return NewWavDecoder(config)
        case Flac:
            return NewFlacDecoder(config)
        case AIFF:
            return NewAIFFDecoder(config)
        }
        return NewUnImplDecoder()
    }

```


# Decoders
## FlacDecoder
The `FlacDecoder` struct is an implementation of the Decoder interface for decoding FLAC audio files. The struct includes a `Duration()` method that calculates the duration of the audio file.
```go
    config := audio.DecoderConfig{
        FilePath:     "samples/audio.flac",
        DecimalLimit: "3",
    }

    decoder,err := audio.NewFlacDecoder(config).Decode()

    duration, err := decoder.Duration()
```

## Mp3Decoder
The `Mp3Decoder` struct is an implementation of the `Decoder` interface for decoding MP3 audio files. The struct includes a `Duration()` method that calculates the duration of the audio file.
```go
    config := audio.DecoderConfig{
        FilePath:     "samples/audio.flac",
        DecimalLimit: "3",
    }

    decoder, err := audio.NewMP3Decoder(config).Decode()

    duration, err := decoder.Duration()
```

## AiffDecoder
The `AiffDecoder` struct is an implementation of the `Decoder` interface for decoding `aiff` audio files. The struct includes a `Duration()` method that calculates the duration of the audio file.
```go
    config := audio.DecoderConfig{
        FilePath:     "samples/audio.flac",
        DecimalLimit: "3",
    }

    decoder,err := audio.NewAIFFDecoder(config).Decode()

    duration, err := decoder.Duration()
```

## WavDecoder
The `WavDecoder` struct is an implementation of the `Decoder` interface for decoding `wav` audio files. The struct includes a `Duration()` method that calculates the duration of the audio file.
```go
    config := audio.DecoderConfig{
        FilePath:     "samples/audio.flac",
        DecimalLimit: "3",
    }

    decoder,err := audio.NewWavDecoder(config).Decode()

    duration, err := decoder.Duration()
```



# Examples 

- **WithBytes**
```go
    config := audio.DecoderConfig{
        DecimalLimit: "3",
    }

    decoder := audio.NewFlacDecoder(config)
    f, _ := os.Open("samples/audio.flac")
    defer f.Close()
    de := decoder.WithBytes(f).decode()
    
    duration, err := de.Duration()
```