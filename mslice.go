package mslice

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"unsafe"

	"github.com/edsrzf/mmap-go"
)

type Flusher interface {
	Flush() error
	Close() error
	// MMap() *mmap.MMap
}

type handle struct {
	m *mmap.MMap
	f *os.File
}

// func (h handle) MMap() *mmap.MMap {
// 	return h.m
// }

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

func mfile(filename string, write bool) (*mmap.MMap, *os.File, error) {
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
		return nil, nil, fmt.Errorf("mmap failed: %w", err)
	}
	// fmt.Printf("MM IS: %+v\n", &mm)
	// fmt.Printf("MM PTR: %p\n", &mm)
	return &mm, f, nil
}

func mslice[T any](mm *mmap.MMap, length, cap int) ([]T, error) {
	if length > cap {
		cap = length
	}
	head := (*reflect.SliceHeader)(unsafe.Pointer(mm))
	ptr := (*T)(unsafe.Pointer(head.Data))
	slice := unsafe.Slice(ptr, cap)
	head = (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	head.Len = int(length)
	reply := ([]T)(slice)
	return reply, nil
}

func aslice(mm *mmap.MMap, length, cap int, slicePtr interface{}) error {
	if length > cap {
		cap = length
	}
	mhead := (*reflect.SliceHeader)(unsafe.Pointer(mm))
	val := reflect.ValueOf(slicePtr)
	ptr := val.Pointer()
	uptr := unsafe.Pointer(ptr)
	shead := (*reflect.SliceHeader)(uptr)
	shead.Len = length
	shead.Cap = cap
	shead.Data = mhead.Data
	return nil
}

func New[T any](filename string, length, cap int) ([]T, Flusher, error) {
	if length > cap {
		cap = length
	}
	var one T
	size := unsafe.Sizeof(one)
	if err := emptyFile(filename, int64(cap*int(size))); err != nil {
		return nil, nil, err
	}

	mm, mf, err := mfile(filename, true)
	if err != nil {
		return nil, nil, err
	}
	slice, err := mslice[T](mm, length, cap)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create mmap slice: %w", err)
	}

	return slice, &handle{mm, mf}, nil
}

func emptyFile(filename string, size int64) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("can't create file %q -- %w", filename, err)
	}
	if err = f.Close(); err != nil {
		return fmt.Errorf("can't close file %q -- %w", filename, err)
	}
	err = os.Truncate(filename, size)
	if err != nil {
		return fmt.Errorf("could not create %q (%d bytes) -- %w", filename, size, err)
	}
	return nil
}

// NewAny takes a pointer to an empty slice and converts it into a mmap'd slice of the requested capacity
func NewAny(filename string, length, cap int, slicePtr interface{}) (Flusher, error) {
	if length > cap {
		cap = length
	}
	size := elemSize(slicePtr)
	if err := emptyFile(filename, int64(cap*int(size))); err != nil {
		return nil, err
	}

	mm, mf, err := mfile(filename, true)
	if err != nil {
		return nil, err
	}

	err = aslice(mm, length, cap, slicePtr)
	if err != nil {
		return nil, fmt.Errorf("could not create mmap slice: %w", err)
	}

	return handle{mm, mf}, nil
}

// OpenAny takes a pointer to an empty slice and converts it into a mmap'd slice of the requested capacity
func OpenAny(filename string, length, cap int, slicePtr interface{}) (Flusher, error) {
	fs, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot stat %q -- %w", filename, err)
	}
	fileSize := int(fs.Size())

	if cap == 0 {
		cap = fileSize / elemSize(slicePtr)
		if length == 0 {
			length = cap
		}
	}

	// TODO: return error?
	if length > cap {
		cap = length
	}

	mm, mf, err := mfile(filename, false)
	if err != nil {
		return nil, err
	}

	err = aslice(mm, length, cap, slicePtr)
	if err != nil {
		return nil, fmt.Errorf("could not create mmap slice: %w", err)
	}
	return handle{mm, mf}, nil
}

// Open returns a mmap-backed slice of type T
func Open[T any](filename string, length, cap int) ([]T, error) {
	var one T
	size := int64(unsafe.Sizeof(one))

	fs, err := os.Stat(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot stat %q -- %w", filename, err)
	}
	fileSize := fs.Size()

	if cap == 0 {
		cap = int(fileSize / size)
	}
	if length > cap {
		log.Printf("adjusting cap %d to match len %d", cap, length)
		cap = length
	}

	need := int64(cap * int(size))
	have := fileSize / size
	// log.Printf("file size (%d) is %d -- need %d", fs.Size(), have, need)
	if need > fs.Size() {
		return nil, fmt.Errorf("file holds %d elements but cap is less at %d", have, size)
	}
	mm, _, err := mfile(filename, false)
	if err != nil {
		return nil, err
	}
	return mslice[T](mm, length, cap)
}

func elemSize(slice interface{}) int {
	val := reflect.ValueOf(slice)
	elem := val.Elem()
	ty := elem.Type()
	return int(ty.Size())
}
