package mslice

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/edsrzf/mmap-go"
)

type Flusher interface {
	Flush() error
	Close() error
}

type handle struct {
	m mmap.MMap
	f *os.File
}

// Flush the mmap to disk
func (h handle) Flush() error {
	return h.m.Flush()
}

// Close will flush the mmap and close the backing file
func (h handle) Close() error {
	if err := h.m.Flush(); err != nil {
		log.Printf("error flushing mmap: %v", err)
	}

	if err := h.m.Unmap(); err != nil {
		log.Printf("error umapping mmap: %v", err)
	}

	err := h.f.Close()
	if err != nil {
		log.Printf("error closing %q -- %v", h.f.Name(), err)
	}
	return err
}

func mfile(filename string, write bool) (mmap.MMap, *os.File, error) {
	var fflags int
	if write {
		fflags |= os.O_RDWR | os.O_CREATE
	}
	f, err := os.OpenFile(filename, fflags, 0755)
	if err != nil {
		return nil, nil, fmt.Errorf("can't open %q -- %w", filename, err)
	}
	var flags int
	if write {
		flags |= mmap.RDWR
	}
	mm, err := mmap.Map(f, flags, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("mmap failed for %q: %w", filename, err)
	}
	return mm, f, nil
}

func emptyFile(filename string, size int64) error {
	// unlike the truncate command, we need to ensure the file exists
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("can't create file %q -- %w", filename, err)
	}
	if err = f.Close(); err != nil {
		return fmt.Errorf("can't close file %q -- %w", filename, err)
	}
	if err = os.Truncate(filename, size); err != nil {
		return fmt.Errorf("could not create %q (%d bytes) -- %w", filename, size, err)
	}
	return nil
}

func elemSize(slice interface{}) int {
	val := reflect.ValueOf(slice)
	elem := val.Elem()
	ty := elem.Type()
	return int(ty.Size())
}

type Encoder interface {
	Encode([]byte) error
	Size() int
}

type Decoder interface {
	Decode([]byte) error
	Size() int
}

type Transcoder interface {
	Encoder
	Decoder
}

type Byter interface {
	Size() int
	Flush() error
	Close() error
	Decode(int, Decoder) error
	Encode(int, Encoder) error
}

type bfile struct {
	f    *os.File
	m    mmap.MMap
	size int64
}

// Size returns the size of the slice in bytes
func (bf *bfile) Size() int {
	return int(bf.size)
}

// Close will flush and close the mmap'd backing file
func (bf *bfile) Close() error {
	if err := bf.m.Flush(); err != nil {
		log.Printf("flush failed: %v", err)
	}
	return bf.f.Close()
}

// Flush writes changes to disk without closing the file
func (bf *bfile) Flush() error {
	return bf.m.Flush()
}

// Encode takes a byte index
func (bf *bfile) Encode(idx int, obj Encoder) error {
	offset := idx //* obj.Size()
	ending := offset + obj.Size()
	bb := bf.m[offset:ending]
	return obj.Encode(bb)
}

// Decode takes a byte index
func (bf *bfile) Decode(idx int, obj Decoder) error {
	offset := idx
	ending := offset + obj.Size()
	bb := bf.m[offset:ending]
	return obj.Decode(bb)
}

func NewByteFile(filename string, length int64) (Byter, error) {
	if err := emptyFile(filename, length); err != nil {
		return nil, err
	}
	if err := os.Truncate(filename, length); err != nil {
		return nil, fmt.Errorf("truncate failure: %w", err)
	}
	mm, mf, err := mfile(filename, true)
	if err != nil {
		return nil, err
	}

	return &bfile{f: mf, m: mm, size: length}, nil
}

func OpenByteFile(filename string, writable bool) (Byter, error) {
	mm, mf, err := mfile(filename, writable)
	if err != nil {
		return nil, err
	}
	stat, err := mf.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat fail: %w", err)
	}
	return &bfile{f: mf, m: mm, size: stat.Size()}, nil
}

func fileSize(filename string) int {
	fs, err := os.Stat(filename)
	if err != nil {
		log.Printf("cannot stat %q -- %v", filename, err)
		return -1
	}
	return int(fs.Size())
}
