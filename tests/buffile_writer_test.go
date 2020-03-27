package tests

import (
	"crypto/md5"
	"fmt"
	buffile "github.com/Dowte/buf-file"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

const hash = "27c8c5fc3875c1ab804885a2cccc8bcb"

// 40 char
var s = []byte("test buf file the readWrite performance.")

// 1024000 * 40 = 40M
func writeFile(writer io.Writer) error {
	for i := 0; i < 1024000; i++ {
		_, err := writer.Write(s)
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

	bufFile, err := buffile.NewBufFile(tempFile.Name(), 1024*1024*4)

	if err != nil {
		t.Error(err)
	}

	writer, err := bufFile.GetWriter()

	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 1024; i++ {
		n, err := writer.Write(s)
		if err != nil {
			t.Error(err)
		}
		if n != 40 {
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

	bufFile, err := buffile.NewBufFile(tempFile.Name(), 1024)

	if err != nil {
		t.Error(err)
	}

	writer, err := bufFile.GetWriter()

	if err != nil {
		t.Error(err)
	}

	var count = 0

	for i := 0; i < 1024; i++ {
		writer.Write(s)

		count += 40

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

	bufFile, err := buffile.NewBufFile(tempFile.Name(), 1024*1024*4)

	if err != nil {
		t.Error(err)
	}

	writer, err := bufFile.GetWriter()

	if err != nil {
		t.Error(err)
	}
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
	if stat.Size() != 40960000 {
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
