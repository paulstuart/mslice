package mslice

import (
	"fmt"
	"os"
	"reflect"
	"unsafe"

	"github.com/edsrzf/mmap-go"
)

func mfile(filename string) (*mmap.MMap, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("can't open %q -- %w", filename, err)
	}
	mm, err := mmap.Map(f, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("mmap failed: %w", err)
	}
	return &mm, nil
}

func New[T any](filename string, length, cap int) ([]T, error) {
	if length > cap {
		cap = length
	}
	var one T
	size := unsafe.Sizeof(one)
	f, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("can't create file %q -- %w", filename, err)
	}
	if err = f.Close(); err != nil {
		return nil, fmt.Errorf("can't close file %q -- %w", filename, err)
	}
	err = os.Truncate(filename, int64(cap*int(size)))
	if err != nil {
		return nil, fmt.Errorf("could not create %q (%d bytes) -- %w", filename, size, err)
	}
	mm, err := mfile(filename)
	if err != nil {
		return nil, err
	}
	head := (*reflect.SliceHeader)(unsafe.Pointer(&mm))
	ptr := (*T)(unsafe.Pointer(head.Data))
	slice := unsafe.Slice(ptr, cap)
	head = (*reflect.SliceHeader)(unsafe.Pointer(&slice))
	head.Len = int(length)
	reply := ([]T)(slice)
	return reply, nil
}
