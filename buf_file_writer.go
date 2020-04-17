package buf_file

import (
	"errors"
	"io"
	"os"
)

var BufFileWriterStopped = errors.New("buf file stop write")

type BufFileWriter struct {
	bufFile    *BufFile
	fileWriter *os.File
	closed     bool
	err        error
}

func (w *BufFileWriter) Flush() error {
	if w.err != nil {
		return w.err
	}
	if w.bufFile.Buffered() == 0 {
		return nil
	}
	n, err := w.fileWriter.Write(w.bufFile.buf[0:w.bufFile.Buffered()])
	if n < w.bufFile.Buffered() && err == nil {
		err = io.ErrShortWrite
	}
	w.bufFile.incFileSize(int64(n))
	if err != nil {
		if n > 0 && n < w.bufFile.Buffered() {
			copy(w.bufFile.buf[0:w.bufFile.Buffered()-n], w.bufFile.buf[n:w.bufFile.Buffered()])
		}
		w.bufFile.incBufferedSize(-n)
		w.err = err
		return err
	}
	w.bufFile.setBufferedSize(0)
	return nil
}

func (w *BufFileWriter) Write(p []byte) (nn int, err error) {
	if w.closed {
		return 0, BufFileWriterStopped
	}
	for len(p) > w.bufFile.Available() && w.err == nil {
		var n int
		if w.bufFile.Buffered() == 0 {
			// Large write, empty buffer.
			// Write directly from p to avoid copy.
			n, w.err = w.fileWriter.Write(p)

			w.bufFile.incFileSize(int64(n))
		} else {
			n = copy(w.bufFile.buf[w.bufFile.Buffered():], p)
			w.bufFile.incBufferedSize(n)
			w.Flush()
		}
		nn += n
		p = p[n:]
	}
	if w.err != nil {
		return nn, w.err
	}
	n := copy(w.bufFile.buf[w.bufFile.Buffered():], p)
	w.bufFile.incBufferedSize(n)
	nn += n
	return nn, nil
}

func (w *BufFileWriter) Close() (err error) {
	err = w.Flush()
	w.closed = true
	err = w.fileWriter.Close()

	return
}
