package mslice

import (
	"encoding/binary"
	"os"
	"path"
	"testing"
)

const (
	testFile = "testdata/testmmap.bin"
	testLen  = 100
)

type testType struct {
	High, Low int64
}

func (t testType) Size() int {
	return 16
}

func (t testType) Encode(b []byte) error {
	w := &Buffer{b}
	return binary.Write(w, binary.LittleEndian, t)
}

func (t *testType) Decode(b []byte) error {
	r := &Buffer{b}
	return binary.Read(r, binary.LittleEndian, t)
}

func prepDir(t *testing.T) {
	t.Helper()
	if err := os.MkdirAll(path.Dir(testFile), os.ModePerm); err != nil {
		t.Fatal(err)
	}
}

type Buffer struct {
	b []byte
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	return copy(b.b, p), nil
}

func (b *Buffer) Read(p []byte) (n int, err error) {
	return copy(p, b.b), nil
}

// ensure our test struct actually works
func TestTranscode(t *testing.T) {
	b := make([]byte, 16)
	want := testType{123, 456}
	var have testType
	if err := want.Encode(b); err != nil {
		t.Fatal(err)
	}
	if err := (&have).Decode(b); err != nil {
		t.Fatal(err)
	}
	if want != have {
		t.Fatalf("have %+v -- want %+v", have, want)
	}
}
