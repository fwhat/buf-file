package buf_file

import (
	"os"
	"sync"
	"sync/atomic"
)

type BufFile struct {
	Filepath       string
	buf            []byte
	bufferedSize   int64
	fileSize       int64
	fileSizeLock   *sync.RWMutex
	buffFileWriter *BufFileWriter
}

func (b *BufFile) writerStopped() bool {
	if b.buffFileWriter != nil {
		return b.buffFileWriter.closed
	}

	return true
}

func (b *BufFile) GetSize() int {
	b.fileSizeLock.RLock()
	defer b.fileSizeLock.RUnlock()

	return int(b.fileSize) + b.Buffered()
}

func (b *BufFile) Buffered() int { return (int)(atomic.LoadInt64(&b.bufferedSize)) }

func (b *BufFile) incBufferedSize(size int) {
	atomic.StoreInt64(&b.bufferedSize, atomic.LoadInt64(&b.bufferedSize)+int64(size))
}

func (b *BufFile) setBufferedSize(size int) {
	atomic.StoreInt64(&b.bufferedSize, int64(size))
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

	return &BufFile{
		buf:          make([]byte, writeBuffSize),
		Filepath:     filepath,
		fileSize:     stat.Size(),
		fileSizeLock: &sync.RWMutex{},
	}, nil
}
