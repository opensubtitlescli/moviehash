package moviehash

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/opensubtitlescli/moviehash/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockFile struct {
	*os.File
	MockReadAt func (b []byte, off int64) (n int, err error)
	MockStat func () (os.FileInfo, error)
}

func (m *mockFile) ReadAt(b []byte, off int64) (n int, err error) {
	return m.MockReadAt(b, off)
}

func (f *mockFile) Stat() (os.FileInfo, error) {
	return f.MockStat()
}

type mockFileInfo struct {
	fs.FileInfo
	MockSize func () int64
}

func (m mockFileInfo) Size() int64 {
	return m.MockSize()
}

func TestSums(t *testing.T) {
	m, err := testdata.FilesMap()
	require.NoError(t, err)
	for _, v := range m {
		f, err := os.Open(v[0])
		if err == nil {
			h, err := Sum(f)
			require.NoError(t, err)
			assert.Equal(t, v[1], h)
		} else {
			if os.IsNotExist(err) {
				fmt.Printf("warn: %s not found\n", v[0])
			} else {
				require.NoError(t, err)
			}
		}
	}
}

func TestSumsWithTheLeadingPadding(t *testing.T) {
	f, teardown := setupMock(t)
	defer teardown()

	f.MockReadAt = func (b []byte, off int64) (n int, err error) {
		return chunkSize, nil
	}
	mockSize(f, fileMinSize)

	a, err := Sum(f)
	require.NoError(t, err)

	assert.Equal(t, "0000000000020000", a)
}

func TestReturnsAnErrorIfItFailsToReadTheStat(t *testing.T) {
	f, teardown := setupMock(t)
	defer teardown()

	f.MockStat = func () (os.FileInfo, error) {
		f.File = nil
		return f.File.Stat()
	}

	_, a := Sum(f)
	assert.EqualError(t, a, "invalid argument")
}

func TestReturnsAnErrorIfTheFileIsTooSmall(t *testing.T) {
	f, teardown := setupMock(t)
	defer teardown()
	mockSize(f, fileMinSize - 1)
	_, a := Sum(f)
	e := fmt.Sprintf("moviehash: the %q file is too small", f.Name())
	assert.EqualError(t, a, e)
}

func TestReturnsAnErrorIfTheFileIsTooLarge(t *testing.T) {
	f, teardown := setupMock(t)
	defer teardown()
	mockSize(f, fileMaxSize + 1)
	_, a := Sum(f)
	e := fmt.Sprintf("moviehash: the %q file is too large", f.Name())
	assert.EqualError(t, a, e)
}

func TestReturnsAnErrorIfItFailsToReadTheFirstChunk(t *testing.T) {
	f, teardown := setupMock(t)
	defer teardown()

	f.MockReadAt = func (b []byte, off int64) (n int, err error) {
		return f.File.ReadAt(nil, -1)
	}
	mockSize(f, fileMinSize)

	_, err := Sum(f)
	a := errors.Unwrap(err)
	assert.EqualError(t, a, "negative offset")
}

func TestReturnsAnErrorIfItReadsTheWrongNumberOfBytesFromTheFirstChunk(t *testing.T) {
	f, teardown := setupMock(t)
	defer teardown()

	f.MockReadAt = func (b []byte, off int64) (n int, err error) {
		return chunkSize - 1, nil
	}
	mockSize(f, fileMinSize)

	_, a := Sum(f)
	e := fmt.Sprintf("moviehash: reads the wrong number of bytes from the first chunk of the %q file", f.Name())
	assert.EqualError(t, a, e)
}

func TestReturnsAnErrorIfItFailsToReadTheLastChunk(t *testing.T) {
	f, teardown := setupMock(t)
	defer teardown()

	f.MockReadAt = func (b []byte, off int64) (n int, err error) {
		if off == 0 {
			return chunkSize, nil
		}
		return f.File.ReadAt(nil, -1)
	}
	mockSize(f, fileMinSize)

	_, err := Sum(f)
	a := errors.Unwrap(err)
	assert.EqualError(t, a, "negative offset")
}

func TestReturnsAnErrorIfItReadsTheWrongNumberOfBytesFromTheLastChunk(t *testing.T) {
	f, teardown := setupMock(t)
	defer teardown()

	f.MockReadAt = func (b []byte, off int64) (n int, err error) {
		if off == 0 {
			return chunkSize, nil
		}
		return chunkSize - 1, nil
	}
	mockSize(f, fileMinSize)

	_, a := Sum(f)
	e := fmt.Sprintf("moviehash: reads the wrong number of bytes from the last chunk of the %q file", f.Name())
	assert.EqualError(t, a, e)
}

func TestReturnsAnErrorIfItFailsToReadTheBothChunks(t *testing.T) {
	fmt.Println("skip: have no idea how to make binary.Read fail")
}

func setupMock(t *testing.T) (*mockFile, func ()) {
	d, err := os.MkdirTemp("", "osub")
	if err != nil {
		os.RemoveAll(d)
	}
	require.NoError(t, err)

	f, err := os.CreateTemp(d, "osub")
	if err != nil {
		f.Close()
	}
	require.NoError(t, err)

	m := &mockFile{
		File: f,
	}
	teardown := func () {
		os.RemoveAll(d)
		f.Close()
	}
	return m, teardown
}

func mockSize(f *mockFile, s int) {
	f.MockStat = func () (os.FileInfo, error) {
		i, err := f.File.Stat()
		if err != nil {
			return i, err
		}
		m := mockFileInfo{
			FileInfo: i,
			MockSize: func () int64 {
				return int64(s)
			},
		}
		return m, err
	}
}
