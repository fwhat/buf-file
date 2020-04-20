package buf_file

import (
	"bufio"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
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
	e := &E{t: t}
	WriteAndManyReaderRandomString(e, 10, 1222166, 333343)
}

func TestBufFileReader_ReadAndWriteLessBuffSize(t *testing.T) {
	e := &E{t: t}
	// 77 char * 1222166 ~ 89M
	WriteAndManyReaderRandomString(e, 10, 1222166, 133)
}

func BenchmarkBufFileReader_ReadAndWriteWithManyReader(b *testing.B) {
	e := &E{b: b}

	WriteAndManyReaderRandomString(e, 10, b.N, 1024*1024*4)
}

// only write
func BenchmarkBufFileReader_ReadAndWrite(b *testing.B) {
	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		b.Error(err)
	}

	bufFile, err := NewBufFile(tempFile.Name(), 1024*1024*4)

	if err != nil {
		b.Error(err)
	}

	writer, err := bufFile.GetWriter()

	if err != nil {
		b.Error(err)
	}

	for i := 0; i < b.N; i++ {
		writer.Write(s)
	}
	writer.Close()

	tempFile.Close()
	os.Remove(tempFile.Name())
}

// 1000000	      1172 ns/op
func BenchmarkBufFileReader_ReadAndWriteSimpleWithManyReader(b *testing.B) {
	e := &E{b: b}
	WriteAndManyReaderSimple(e, 20, b.N, 1024*1024*4)
}

func BenchmarkBufioWrite(b *testing.B) {
	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		b.Error(err)
	}

	writer := bufio.NewWriterSize(tempFile, 1024*1024*4)

	for i := 0; i < b.N; i++ {
		writer.Write(s)
	}
	writer.Flush()

	tempFile.Close()
	os.Remove(tempFile.Name())
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

func WriteAndManyReaderRandomString(t *E, readerCount int, writeCount int, bufSize int) {
	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Error(err)
	}

	bufFile, err := NewBufFile(tempFile.Name(), bufSize)

	if err != nil {
		t.Error(err)
	}

	writer, err := bufFile.GetWriter()

	n := writeCount
	var writeContent []byte

	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		for i := 0; i < n; i++ {
			str := []byte(RandomString(77))
			writeContent = append(writeContent, str...)
			writer.Write(str)
		}
		writer.Close()
		wg.Done()
	}()

	for i := 0; i < readerCount-1; i++ {
		wg.Add(1)
		go func() {
			reader, err := bufFile.GetReader()
			if err != nil {
				t.Error(err)
				return
			}
			data := readFileWithReadAt(reader, 77*n)
			if hashx(data) != hashx(writeContent) {
				t.Error(err)
			}
			reader.Close()
			wg.Done()
		}()
	}

	if readerCount > 0 {
		reader, err := bufFile.GetReader()
		if err != nil {
			t.Error(err)
		}

		data := readFileWithReadAt(reader, 77*n)

		if hashx(data) != hashx(writeContent) {
			t.Error()
		}
		reader.Close()
	}

	wg.Wait()

	tempFile.Close()
	os.Remove(tempFile.Name())

	return
}

func WriteAndManyReaderSimple(t *E, readerCount int, writeCount int, bufSize int) {
	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Error(err)
	}

	bufFile, err := NewBufFile(tempFile.Name(), bufSize)

	if err != nil {
		t.Error(err)
	}

	writer, err := bufFile.GetWriter()

	n := writeCount

	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		for i := 0; i < n; i++ {
			writer.Write(s)
		}
		writer.Close()
		wg.Done()
	}()

	for i := 0; i < readerCount-1; i++ {
		wg.Add(1)
		go func() {
			reader, err := bufFile.GetReader()
			if err != nil {
				t.Error(err)
				return
			}
			readFileWithReadAt(reader, 64*n)

			reader.Close()
			wg.Done()
		}()
	}

	if readerCount > 0 {
		reader, err := bufFile.GetReader()
		if err != nil {
			t.Error(err)
		}

		readFileWithReadAt(reader, 64*n)
		reader.Close()
	}

	wg.Wait()

	tempFile.Close()
	os.Remove(tempFile.Name())

	return
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandomString(length int) string {
	return StringWithCharset(length, charset)
}

type E struct {
	b *testing.B
	t *testing.T
}

func (e E) Error(args ...interface{}) {
	if e.b != nil {
		e.b.Error(args...)
	} else {
		e.t.Error(args...)
	}
}
