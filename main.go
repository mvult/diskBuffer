package diskBuffer

import (
	_ "fmt"
	"io"
	"os"
	"sync"
)

type Buffer struct {
	f               *os.File
	offset          int64
	inboundComplete bool
	lock            sync.Mutex
}

func New(f *os.File) *Buffer {
	return &Buffer{f: f, inboundComplete: false}
}

func (fb *Buffer) SetInboundComplete(b bool) {
	fb.lock.Lock()
	defer fb.lock.Unlock()
	fb.inboundComplete = b
}

func (fb *Buffer) GetInboundComplete() bool {
	fb.lock.Lock()
	defer fb.lock.Unlock()
	return fb.inboundComplete
}

func (fb *Buffer) Remove() error {
	fb.lock.Lock()
	fb.lock.Unlock()
	return os.Remove(fb.f.Name())
}

func (fb *Buffer) Read(b []byte) (int, error) {
	fb.lock.Lock()
	defer fb.lock.Unlock()
	n, err := fb.f.ReadAt(b, fb.offset)
	fb.offset += int64(n)
	// fmt.Println(n, err)
	if err == io.EOF {
		if fb.inboundComplete {
			return n, io.EOF
		}
		return n, nil
	}
	return n, err
}

func (fb *Buffer) Write(p []byte) (int, error) {
	fb.lock.Lock()
	defer fb.lock.Unlock()
	return fb.f.Write(p)
}

func (fb *Buffer) Close() error {
	fb.lock.Lock()
	defer fb.lock.Unlock()
	return fb.f.Close()
}

func (fb *Buffer) CloseAndRemove() error {
	fb.lock.Lock()
	if err := fb.f.Close(); err != nil {
		return err
	}
	return os.Remove(fb.f.Name())
}
