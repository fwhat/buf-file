package buf_file

import (
	"os"
	"sync"
	"sync/atomic"
)

type BufFile struct {
	Filepath       string
	buf            []byte
	bufferedSize   *atomic.Value
	fileSize       int64
	fileSizeLock   *sync.Mutex
	buffFileWriter *BufFileWriter
}

func (b *BufFile) writerStopped() bool {
	if b.buffFileWriter != nil {
		return b.buffFileWriter.closed
	}

	return true
}

func (b *BufFile) Buffered() int { return b.bufferedSize.Load().(int) }

func (b *BufFile) incBufferedSize(size int) {
	b.bufferedSize.Store(b.bufferedSize.Load().(int) + size)
}

func (b *BufFile) setBufferedSize(size int) {
	b.bufferedSize.Store(size)
}

func (b *BufFile) Available() int { return len(b.buf) - b.Buffered() }

func (b *BufFile) GetWriter() (*BufFileWriter, error) {
	if b.buffFileWriter == nil {
		file, err := os.OpenFile(b.Filepath, os.O_WRONLY, 0)
		if err != nil {
			return nil, err
		}

		b.buffFileWriter = &BufFileWriter{
			bufFile:    b,
			fileWriter: file,
		}
	}

	return b.buffFileWriter, nil
}

func (b *BufFile) GetReader() (*BufFileReader, error) {
	file, err := os.Open(b.Filepath)

	if err != nil {
		return nil, err
	}

	return &BufFileReader{
		bufFile:    b,
		fileReader: file,
	}, err
}

func NewBufFile(filepath string, writeBuffSize int) (*BufFile, error) {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	buffSize := &atomic.Value{}
	buffSize.Store(0)

	return &BufFile{
		buf:          make([]byte, writeBuffSize),
		bufferedSize: buffSize,
		Filepath:     filepath,
		fileSize:     stat.Size(),
		fileSizeLock: &sync.Mutex{},
	}, nil
}
