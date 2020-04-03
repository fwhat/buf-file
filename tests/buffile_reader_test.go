package tests

import (
	buf_file "github.com/Dowte/buf-file"
	"io"
	"os"
	"testing"
	"time"
)

func TestReaderClosed(t *testing.T) {
	bufFile, err := buf_file.NewBufFile("/tmp/buf_file", 1024*1024*4)
	if err != nil {
		panic(err)
	}
	writer, err := bufFile.GetWriter()

	reader, err := bufFile.GetReader()

	if err != nil {
		panic(err)
	}
	go func() {
		reader.Close()
		time.Sleep(time.Second * 1)
		writer.Write(s)
	}()

	written, err := io.Copy(writer, &io.LimitedReader{R: reader, N: 10})
	if written != 0 {
		t.Error("close fail")
	}

	if err != os.ErrClosed {
		t.Error("err invalid")
	}
}
