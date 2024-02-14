package audio

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

type AudioFileFormat string

const (
	MP3  AudioFileFormat = "mp3"
	Wav  AudioFileFormat = "wav"
	Flac AudioFileFormat = "flac"
	AIFF AudioFileFormat = "aiff"
	OGG  AudioFileFormat = "ogg"
)

type Decoder interface {
	Decode() (Decoder, error)
	WithBytes(io.ReadSeeker) Decoder
	Bytes() io.ReadSeeker
	Close() error

	Duration() (float64, error)
	SampleRate() (float64, error)
}

type DecoderConfig struct {
	DecimalLimit string
	FilePath     string
	RemoteFile   bool

	bytes  io.ReadSeeker
	closer io.Closer
}

// NewDecoder
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
	case OGG:
		return NewOGGDecoder(config)
	}
	return NewUnImplDecoder()
}

// LoadFile
// The LoadFile function is used to load a file from either a local file path or a remote URL.
// The function takes two arguments: filepath and remoteFile. If remoteFile is true, the function downloads the file from the specified URL and
// returns an io.Reader and io.Closer for the downloaded file. If remoteFile is false, the function opens the file at the specified file path and
// returns an io.Reader and io.Closer for the opened file.

// The function first checks if remoteFile is true. If it is, the function downloads the file from the specified URL using the http.Get() function and
// returns an io.Reader and io.Closer for the downloaded file.
// If remoteFile is false, the function opens the file at the specified file path using the os.Open() function and returns an io.Reader and io.Closer
// for the opened file.
// If there is an error while downloading or opening the file, the function returns nil for both the io.Reader and io.Closer, and the error.

func LoadFile(filepath string, remoteFile bool) (io.ReadSeeker, io.Closer, error) {
	var err error
	if remoteFile {
		// were the file is getting from remote
		req, err := http.NewRequest(http.MethodGet, filepath, nil)
		if err != nil {
			return nil, nil, err
		}
		client := &http.Client{
			// Timeout: time.Second * 10,
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, nil, fmt.Errorf("failed to retrieve the file. Status code: %d", resp.StatusCode)
		}

		// create a buffer and copy the contents of the reader to the buffer
		buffer := new(bytes.Buffer)
		_, err = io.Copy(buffer, resp.Body)
		if err != nil {
			return nil, nil, err
		}

		return bytes.NewReader(buffer.Bytes()), resp.Body, nil
	}

	// open the file
	file, err := os.Open(filepath)
	if err != nil {
		return nil, nil, err
	}

	return file, file, nil
}
