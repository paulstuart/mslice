package mslice

import (
	"fmt"
)

// "github.com/tidwall/mmap"

type SliceFile interface {
	Flusher
	Cap() int
	Len() int
	Append(...Encoder) error
	Set(int, Encoder) error
	Get(int, Decoder) error
}

type sfile struct {
	b    Byter
	size int
	len  int
	cap  int
}

func (sf *sfile) Close() error {
	return sf.b.Close()
}

func (sf *sfile) Flush() error {
	return sf.b.Flush()
}

func (sf *sfile) Cap() int {
	return sf.cap
}

func (sf *sfile) Len() int {
	return sf.len
}

// ErrAppend is when we go past the mmap limit
// TODO: make it extendable
var ErrAppend = fmt.Errorf("append exceeds size")

// Append objects to the slice
func (sf *sfile) Append(objs ...Encoder) error {
	add := len(objs)
	want := sf.len + add
	if want > sf.cap {
		return fmt.Errorf("size:%d cap:%d add:%d -- %w", want, sf.cap, add, ErrAppend)
	}
	for i, obj := range objs {
		if err := sf.b.Encode(sf.len*obj.Size(), obj); err != nil {
			return fmt.Errorf("error encoding (%d/%d): %w", i+1, add, err)
		}
		sf.len++
	}
	return nil
}

func (sf *sfile) Size() int {
	return sf.b.Size()
}

func (sf *sfile) Set(idx int, obj Encoder) error {
	if idx > sf.len {
		return fmt.Errorf("idx %d exceeds len %d", idx, sf.len)
	}
	return sf.b.Encode(sf.size*idx, obj)
}

func (sf *sfile) Get(idx int, obj Decoder) error {
	if idx > sf.len {
		return fmt.Errorf("idx %d exceeds len %d", idx, sf.len)
	}
	return sf.b.Decode(idx*obj.Size(), obj)
}

// NewSlice takes a pointer to an empty slice and converts it into a mmap'd slice of the requested capacity
// NOTE: the object must comprise only fixed size elements, i.e., numbers
func NewSlice(filename string, length, cap int, slicePtr interface{}) (SliceFile, error) {
	if length > cap {
		cap = length
	}

	size := elemSize(slicePtr)
	total := int64(cap * size)
	b, err := NewByteFile(filename, total)
	if err != nil {
		return nil, fmt.Errorf("no slicefile for you: %w", err)
	}

	return &sfile{b: b, len: length, cap: cap, size: size}, nil
}

// OpenAny takes a pointer to an empty slice and converts it into a mmap'd slice of the requested capacity
func OpenSlice(filename string, writable bool, slicePtr interface{}) (SliceFile, error) {
	b, err := OpenByteFile(filename, writable)
	if err != nil {
		return nil, fmt.Errorf("no slicefile for you: %w", err)
	}

	size := elemSize(slicePtr)
	length := b.Size() / size

	return &sfile{b: b, len: length, cap: length, size: size}, nil
}
