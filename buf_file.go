package buf_file

import (
	"io"
	"os"
	"sync"
	"sync/atomic"
)

type BufFile struct {
	file         *os.File
	buf          []byte
	buffSize     int64
	err          error
	fileSize     int64
	fileSizeLock *sync.Mutex
}

func (b *BufFile) Buffered() int { return int(atomic.LoadInt64(&b.buffSize)) }

func (b *BufFile) incBufferedSize(size int64) int {
	return int(atomic.AddInt64(&b.buffSize, size))
}

func (b *BufFile) setBufferedSize(size int64) {
	atomic.StoreInt64(&b.buffSize, size)
}

func (b *BufFile) Available() int { return len(b.buf) - b.Buffered() }

func (b *BufFile) Flush() error {
	b.fileSizeLock.Lock()
	defer b.fileSizeLock.Unlock()
	if b.err != nil {
		return b.err
	}
	if b.Buffered() == 0 {
		return nil
	}
	n, err := b.file.Write(b.buf[0:b.Buffered()])
	if n < b.Buffered() && err == nil {
		err = io.ErrShortWrite
	}
	b.fileSize += int64(n)
	if err != nil {
		if n > 0 && n < b.Buffered() {
			copy(b.buf[0:b.Buffered()-n], b.buf[n:b.Buffered()])
		}
		b.incBufferedSize(int64(-n))
		b.err = err
		return err
	}
	b.setBufferedSize(0)
	return nil
}

func (b *BufFile) Write(p []byte) (nn int, err error) {
	for len(p) > b.Available() && b.err == nil {
		var n int
		if b.Buffered() == 0 {
			b.fileSizeLock.Lock()
			// Large write, empty buffer.
			// Write directly from p to avoid copy.
			n, b.err = b.file.Write(p)

			b.fileSize += int64(n)
			b.fileSizeLock.Unlock()
		} else {
			n = copy(b.buf[b.Buffered():], p)
			b.incBufferedSize(int64(n))
			b.Flush()
		}
		nn += n
		p = p[n:]
	}
	if b.err != nil {
		return nn, b.err
	}
	n := copy(b.buf[b.Buffered():], p)
	b.incBufferedSize(int64(n))
	nn += n
	return nn, nil
}

func (b *BufFile) ReadAt(p []byte, offset int64) (n int, err error) {
	b.fileSizeLock.Lock()
	defer b.fileSizeLock.Unlock()
	if offset < b.fileSize {
		n, err = b.file.ReadAt(p, offset)
	}

	// 剩余空间由buff中读取
	if len(p)-n > 0 {
		cn := 0
		if n > 0 {
			cn = copy(p[n:], b.buf[0:b.Buffered()])
		} else {
			cn = copy(p[n:], b.buf[(offset-b.fileSize):b.Buffered()])
		}
		n = cn + n
	}

	return n, nil
}

func NewBufFile(file *os.File, writeBuffSize int) *BufFile {
	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	return &BufFile{
		buf:          make([]byte, writeBuffSize),
		file:         file,
		fileSize:     stat.Size(),
		fileSizeLock: &sync.Mutex{},
	}
}
