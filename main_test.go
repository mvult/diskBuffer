package diskBuffer

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
)

func TestFillFile(t *testing.T) {
	source, err := os.Create("source.data")
	if err != nil {
		panic(err)
	}
	fmt.Println(source.Name())

	chunk := make([]byte, 1024*1024)
	var n int

	for i := 0; i < 200; i++ {
		n, err = rand.Read(chunk)
		if err != nil {
			t.Error(err)
			panic(err)
		}
		source.Write(chunk[:n])
	}

	off, err := source.Seek(int64(0), 0)
	if off != int64(0) || err != nil {
		t.Error(err)
	}

	db := writeToBuffer(source)
	destination := writeBufferToDestination(db)

	h1 := sha256.New()
	if _, err := io.Copy(h1, source); err != nil {
		panic(err)
	}

	h2 := sha256.New()
	if _, err := io.Copy(h2, destination); err != nil {
		panic(err)
	}

	if string(h1.Sum(nil)) != string(h2.Sum(nil)) {
		t.Error("Unexpected hash values")
	}
	source.Close()
	destination.Close()
	if err = os.Remove("source.data"); err != nil {
		panic(err)
	}

	if err = os.Remove("destination.data"); err != nil {
		panic(err)
	}

}

func writeToBuffer(source *os.File) *Buffer {
	// f, err := os.Create("buffer.data")
	f, err := ioutil.TempFile("", "buffers")
	if err != nil {
		panic(err)
	}
	fmt.Println(f.Name())
	db := New(f)

	chunk := make([]byte, 1024*1024*3)
	var n int
	go func() {

		for {
			n, err = source.Read(chunk)
			if err != nil {
				if err == io.EOF {
					db.SetInboundComplete(true)
					return
				}
				panic(err)
			}
			if _, err = db.Write(chunk[:n]); err != nil {
				panic(err)
			}
		}
	}()
	return db
}

func writeBufferToDestination(db *Buffer) *os.File {

	f, err := os.Create("destination.data")
	if err != nil {
		panic(err)
	}

	var n int

	chunk := make([]byte, 1024*1024*3)
	for {
		n, err = db.Read(chunk)
		if err != nil {
			if err == io.EOF {
				if _, err = f.Write(chunk[:n]); err != nil {
					panic(err)
				}
				db.CloseAndRemove()
				return f
			}
			panic(err)
		}
		if _, err = f.Write(chunk[:n]); err != nil {
			panic(err)
		}
	}
	return &os.File{}
}
