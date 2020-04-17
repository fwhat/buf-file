package buf_file

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestBufFileReader_Close(t *testing.T) {
	os.Remove("/tmp/buf_file")
	bufFile, err := NewBufFile("/tmp/buf_file", 1024*1024*4)
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

func TestBufFileReader_ReadAndWrite(t *testing.T) {
	WriteAndManyReader(t, 10, 1222166)
}

func readFileWithReadAt(read io.ReaderAt, total int) []byte {
	var readContent []byte
	var offset int64
	for true {
		buf := make([]byte, 1024*10)
		at, err := read.ReadAt(buf, offset)
		offset += int64(at)
		if at > 0 {
			readContent = append(readContent, buf[:at]...)
		}
		if len(readContent) == total {
			break
		}
		if err != nil {
			panic(err)
		}
	}

	return readContent
}

func WriteAndManyReader(t *testing.T, readerCount int, writeCount int) {
	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Error(err)
	}

	bufFile, err := NewBufFile(tempFile.Name(), 333113)

	if err != nil {
		t.Error(err)
	}

	writer, err := bufFile.GetWriter()

	n := writeCount

	go func() {
		for i := 0; i < n; i++ {
			writer.Write(s)
		}
		writer.Close()
	}()

	for i := 0; i < readerCount-1; i++ {
		go func() {
			reader, err := bufFile.GetReader()
			if err != nil {
				t.Error(err)
			}
			data := readFileWithReadAt(reader, 64*n)
			all, _ := ioutil.ReadFile(tempFile.Name())
			if hashx(data) != hashx(all) {
				t.Error(err)
			}
			reader.Close()
		}()
	}

	reader, err := bufFile.GetReader()
	if err != nil {
		t.Error(err)
	}

	data := readFileWithReadAt(reader, 64*n)

	all, _ := ioutil.ReadFile(tempFile.Name())
	if hashx(data) != hashx(all) {
		t.Error()
	}
	reader.Close()

	tempFile.Close()
	os.Remove(tempFile.Name())

	return
}
