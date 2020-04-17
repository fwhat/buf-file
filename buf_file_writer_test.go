package buf_file

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

const hash = "6273962f98e8f3591daa492de531a209"

// 64 char
var s = []byte("test buf file the readWrite performance when write 400 * 1048576")
var repeatStr []byte // 4096 char

func init() {
	for i := 0; i < 64; i++ {
		repeatStr = append(repeatStr, s...)
	}
}

// 4096 *1024 * 10 = 400M 419430400
func writeFile(writer io.Writer) error {
	// 1024 * 10
	for i := 0; i < 102400; i++ {
		_, err := writer.Write(repeatStr)
		if err != nil {
			return err
		}
	}

	return nil
}

// 测试 (w *BufFileWriter) Write 返回的 n, err
func TestWriteSize(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Error(err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	bufFile, err := NewBufFile(tempFile.Name(), 1024*1024*4)

	if err != nil {
		t.Error(err)
	}

	writer, err := bufFile.GetWriter()

	if err != nil {
		t.Error(err)
	}
	defer writer.Close()

	for i := 0; i < 1024; i++ {
		n, err := writer.Write(s)
		if err != nil {
			t.Error(err)
		}
		if n != 64 {
			t.Error("write size invalid")
		}
	}
}

func TestWriteBuffered(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Error(err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	bufFile, err := NewBufFile(tempFile.Name(), 1024)

	if err != nil {
		t.Error(err)
	}

	writer, err := bufFile.GetWriter()

	if err != nil {
		t.Error(err)
	}
	defer writer.Close()

	var count = 0

	for i := 0; i < 1024; i++ {
		writer.Write(s)

		count += 64

		if count > 1024 {
			count = count - 1024
		}
		if bufFile.Buffered() != count {
			t.Error("buffer size invalid")
		}
	}
}

func TestWriteContent(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Error(err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	bufFile, err := NewBufFile(tempFile.Name(), 1024*1024*4)

	if err != nil {
		t.Error(err)
	}

	writer, err := bufFile.GetWriter()

	if err != nil {
		t.Error(err)
	}
	defer writer.Close()

	err = writeFile(writer)
	if err != nil {
		t.Error(err)
	}
	err = writer.Close()
	if err != nil {
		t.Error(err)
	}

	file, err := os.OpenFile(tempFile.Name(), os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		t.Error(err)
	}
	stat, _ := file.Stat()
	if stat.Size() != 419430400 {
		t.Error(err)
	}

	all, err := ioutil.ReadAll(file)

	if err != nil {
		t.Error(err)
	}

	if hashx(all) != hash {
		t.Error("write content invalid")
	}
}

func hashx(TestString []byte) string {
	Md5Inst := md5.New()
	Md5Inst.Write(TestString)
	Result := Md5Inst.Sum([]byte(""))

	return fmt.Sprintf("%x", Result)
}

func BenchmarkBufFileWriter_Write(b *testing.B) {
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

func BenchmarkBufioWriter_Write(b *testing.B) {
	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		b.Error(err)
	}

	writer := bufio.NewWriterSize(tempFile, 1024*1024*4)

	for i := 0; i < b.N; i++ {
		writer.Write(s)
	}

	tempFile.Close()
	os.Remove(tempFile.Name())
}
