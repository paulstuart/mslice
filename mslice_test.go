package mslice

import (
	"log"
	"os"
	"path"
	"testing"
)

const (
	testFile = "testdata/testmmap.bin"
	testLen  = 10
)

type testType struct {
	High, Low int64
}

func prepDir(t *testing.T) {
	t.Helper()
	if err := os.MkdirAll(path.Dir(testFile), os.ModePerm); err != nil {
		t.Fatal(err)
	}
}

func TestMSliceNewInt(t *testing.T) {
	prepDir(t)
	m, h, err := New[int64](testFile, testLen, testLen)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(m); i++ {
		m[i] = int64(i+1) * 100
	}
	// mm := h.MMap()
	// fmt.Printf("MM PTR NOW: %p\n", mm)

	//	fmt.Printf("MM NOW: %+v\n", mm)

	t.Logf("size: %d", len(m))
	for i, v := range m {
		t.Logf("%02d: %v\n", i, v)
	}
	// fmt.Printf("MM IS NOW: %+v\n", mm)
	if err := h.Close(); err != nil {
		log.Fatalf("can't close: %v", err)
	}
}

func TestMSliceIntAppend(t *testing.T) {
	prepDir(t)
	m, _, err := New[int64](testFile, 0, testLen)
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

func TestOpen(t *testing.T) {
	prepDir(t)
	m, h, err := New[int64](testFile, 0, testLen)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < cap(m); i++ {
		m = append(m, int64(i+1)*100)
	}
	for i, v := range m {
		t.Logf("%02d: %v\n", i, v)
	}
	// mm := h.MMap()
	// fmt.Printf("MM NOW: %p\n", mm)
	newLen := len(m)
	t.Logf("new len: %d", newLen)
	h.Close()

	m, err = Open[int64](testFile, newLen, testLen)
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range m {
		t.Logf("%02d: %v\n", i, v)
	}
}

func TestMSliceInt32Append(t *testing.T) {
	prepDir(t)
	m, h, err := New[int32](testFile, 0, testLen)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < cap(m); i++ {
		m = append(m, int32(i+1)*100)
	}
	if err := h.Flush(); err != nil {
		t.Errorf("flush failed: %v", err)
	}
	t.Logf("size: %d", len(m))
	for i, v := range m {
		t.Logf("%02d: %v\n", i, v)
	}
}

func TestMSliceFloat(t *testing.T) {
	prepDir(t)
	m, _, err := New[float64](testFile, testLen, testLen)
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

func TestNewAny(t *testing.T) {
	prepDir(t)
	m := []testType{}
	h, err := NewAny(testFile, testLen, testLen, &m)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(m); i++ {
		plus := int64(i + 1)
		m[i] = testType{plus * 100, plus}
	}
	t.Logf("size: %d", len(m))
	for i, v := range m {
		t.Logf("%02d: %v\n", i, v)
	}
	if err := h.Close(); err != nil {
		t.Fatalf("bummer: %v", err)
	}

	t.Logf("iterate saved file:")

	h, err = OpenAny(testFile, testLen, testLen, &m)
	if err != nil {
		t.Fatalf("can't recycle: %v", err)
	}

	for i, v := range m {
		t.Logf("%02d: %v\n", i, v)
	}
}

func TestNewAnyAuto(t *testing.T) {
	prepDir(t)
	m := []testType{}
	h, err := NewAny(testFile, testLen, testLen, &m)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < len(m); i++ {
		plus := int64(i + 1)
		m[i] = testType{plus * 100, plus}
	}
	t.Logf("size: %d", len(m))
	for i, v := range m {
		t.Logf("%02d: %v\n", i, v)
	}
	if err := h.Close(); err != nil {
		t.Fatalf("bummer: %v", err)
	}

	t.Logf("iterate saved file:")

	h, err = OpenAny(testFile, 0, 0, &m)
	if err != nil {
		t.Fatalf("can't recycle: %v", err)
	}

	for i, v := range m {
		t.Logf("%02d: %v\n", i, v)
	}
}
