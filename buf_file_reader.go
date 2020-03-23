package buf_file

import (
	"io"
	"os"
)

type BufFileReader struct {
	bufFile    *BufFile
	offset     int64
	fileReader *os.File
}

// reader need close
func (r *BufFileReader) Close() error {
	return r.fileReader.Close()
}

func (r *BufFileReader) Read(p []byte) (n int, err error) {
	n, err = r.ReadAt(p, r.offset)
	r.offset += int64(n)

	return
}

func (r *BufFileReader) ReadAt(p []byte, offset int64) (n int, err error) {
	r.bufFile.fileSizeLock.Lock()
	defer r.bufFile.fileSizeLock.Unlock()
	if offset < r.bufFile.fileSize {
		n, err = r.fileReader.ReadAt(p, offset)
	}

	if err != nil {
		if err == io.EOF {
			err = nil
		} else {
			return n, err
		}
	}

	if n == 0 && r.bufFile.writerStopped() {
		return n, io.EOF
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

	return n, nil
}
