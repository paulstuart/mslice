package mslice

import (
	"os"
	"path"
	"testing"
)

const (
	testFile = "testdata/testmmap.bin"
	testLen  = 10
)

func prepDir(t *testing.T) {
	t.Helper()
	if err := os.MkdirAll(path.Dir(testFile), os.ModePerm); err != nil {
		t.Fatal(err)
	}
}

func TestMSliceInt(t *testing.T) {
	prepDir(t)
	m, err := New[int64](testFile, testLen, testLen)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(m); i++ {
		m[i] = int64(i+1) * 100
	}
	t.Logf("size: %d", len(m))
	for i, v := range m {
		t.Logf("%02d: %v\n", i, v)
	}
}

func TestMSliceIntAppend(t *testing.T) {
	prepDir(t)
	m, err := New[int64](testFile, 0, testLen)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < cap(m); i++ {
		m = append(m, int64(i+1)*100)
	}
	t.Logf("size: %d", len(m))
	for i, v := range m {
		t.Logf("%02d: %v\n", i, v)
	}
}

func TestMSliceInt32Append(t *testing.T) {
	prepDir(t)
	m, err := New[int32](testFile, 0, testLen)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < cap(m); i++ {
		m = append(m, int32(i+1)*100)
	}
	t.Logf("size: %d", len(m))
	for i, v := range m {
		t.Logf("%02d: %v\n", i, v)
	}
}
func TestMSliceFloat(t *testing.T) {
	prepDir(t)
	m, err := New[float64](testFile, testLen, testLen)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(m); i++ {
		m[i] = (float64(i+1) * 100) + 0.12345
	}
	t.Logf("size: %d", len(m))
	for i, v := range m {
		t.Logf("%02d: %v\n", i, v)
	}
}
