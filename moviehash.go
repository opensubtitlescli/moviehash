// Provides a hash function to match subtitle files against movie files.
package moviehash

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	Version = "0.0.1"

	chunkSize     =      65536 // byte
	fileMaxSize   = 9000000000 // byte
	fileMinSize   =     131072 // byte
	hashMinLength =         16 // char
)

type file interface {
	Name() string
	Stat() (os.FileInfo, error)
	ReadAt(b []byte, off int64) (n int, err error)
}

// Calculates the hash value of a file using the moviehash algorithm. It returns
// an error if the file cannot be read, or if the file is either too small or
// too large.
//
// [OpenSubtitles Reference]
//
// [OpenSubtitles Reference]: https://opensubtitles.stoplight.io/docs/opensubtitles-api/e3750fd63a100-getting-started#calculating-moviehash-of-video-file
func Sum(file file) (string, error) {
	fi, err := file.Stat()
	if err != nil {
		return "", err
	}

	s := fi.Size()
	if s < fileMinSize {
		return "", fmt.Errorf("moviehash: the %q file is too small", file.Name())
	}
	if s > fileMaxSize {
		return "", fmt.Errorf("moviehash: the %q file is too large", file.Name())
	}

	c := uint64(s)

	b := make([]byte, chunkSize * 2)

	n, err := file.ReadAt(b[:chunkSize], 0)
	if err != nil {
		return "", err
	}
	if n != chunkSize {
		return "", fmt.Errorf("moviehash: reads the wrong number of bytes from the first chunk of the %q file", file.Name())
	}

	n, err = file.ReadAt(b[chunkSize:], s - chunkSize)
	if err != nil {
		return "", err
	}
	if n != chunkSize {
		return "", fmt.Errorf("moviehash: reads the wrong number of bytes from the last chunk of the %q file", file.Name())
	}

	var d [chunkSize * 2 / 8]uint64
	r := bytes.NewReader(b)
	err = binary.Read(r, binary.LittleEndian, &d)
	if err != nil {
		return "", err
	}

	for _, p := range d {
		c += p
	}

	h := strconv.FormatUint(c, 16)
	l := len(h)
	if l < hashMinLength {
		o := strings.Repeat("0", hashMinLength - l)
		h = o + h
	}

	return h, nil
}
