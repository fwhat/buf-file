package buf_file

import (
	"io"
	"os"
	"time"
)

type BufFileReader struct {
	bufFile    *BufFile
	offset     int64
	fileReader *os.File
	closed     bool
}

// reader need close
func (r *BufFileReader) Close() error {
	r.closed = true
	return r.fileReader.Close()
}

func (r *BufFileReader) Read(p []byte) (n int, err error) {
	n, err = r.ReadAt(p, r.offset)
	r.offset += int64(n)

	return
}

func (r *BufFileReader) Seek(offset int64, whence int) (int64, error) {
	r.offset = int64(whence) + offset

	return r.offset, nil
}

func (r *BufFileReader) ReadAt(p []byte, offset int64) (n int, err error) {
	r.bufFile.fileSizeLock.RLock()
	if offset < r.bufFile.fileSize {
		n, err = r.fileReader.ReadAt(p, offset)
	}

	if err != nil {
		if err == io.EOF {
			err = nil
		} else {
			goto end
		}
	}

	if n == 0 && r.bufFile.writerStopped() {
		err = io.EOF
		goto end
	}

	if n == 0 && r.closed {
		err = os.ErrClosed
		goto end
	}

	// 剩余空间由buff中读取
	if len(p)-n > 0 {
		cn := 0
		if n > 0 {
			cn = copy(p[n:], r.bufFile.buf[0:r.bufFile.Buffered()])
		} else {
			cn = copy(p[n:], r.bufFile.buf[(offset-r.bufFile.fileSize):r.bufFile.Buffered()])
		}
		n = cn + n
	}

end:
	r.bufFile.fileSizeLock.RUnlock()
	if n == 0 {
		time.Sleep(50 * time.Millisecond)
	}

	return n, err
}
