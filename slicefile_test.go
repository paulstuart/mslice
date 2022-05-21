package mslice

import (
	"testing"
)

func TestSliceFileNew(t *testing.T) {
	prepDir(t)
	var m testType
	h, err := NewSlice(testFile, 0, testLen, &m)
	if err != nil {
		t.Fatal(err)
	}
	want := make([]testType, testLen)
	// t.Logf("file size: %d", h.Size())
	for i := 0; i < testLen; i++ {
		plus := int64(i + 1)
		obj := testType{plus * 123456, 789 * plus}
		t.Logf("APPEND: %d", i)
		if err := h.Append(obj); err != nil {
			t.Fatalf("append failed: %v", err)
		}
		want[i] = obj
	}
	// t.Logf("size: %d", len(m))
	// for i, v := range m {
	// 	t.Logf("%02d: %v\n", i, v)
	// }
	if err := h.Close(); err != nil {
		t.Fatalf("bummer: %v", err)
	}
	t.Logf("iterate saved file:")

	sf, err := OpenSlice(testFile, false, &m)
	// h = SliceFile(sf)
	// h = sf.(SliceFile)
	if err != nil {
		t.Fatalf("can't recycle: %v", err)
	}

	t.Logf("LEN:%d", sf.Len())
	for i := 0; i < sf.Len(); i++ {
		var obj testType
		if err := sf.Get(i, &obj); err != nil {
			t.Fatalf("decode %d fail: %v", i, err)
		}
		t.Logf("======> GOT (%02d): %v\n", i, obj)
		if obj != want[i] {
			t.Fatalf("(%2d) want %v -- have %v", i, want[i], obj)
		}
	}
}

func TestSliceFileSet(t *testing.T) {
	prepDir(t)
	var m testType
	h, err := NewSlice(testFile, testLen, testLen, &m)
	if err != nil {
		t.Fatal(err)
	}
	want := make([]testType, testLen)
	for i := 0; i < testLen; i++ {
		plus := int64(i + 1)
		obj := testType{plus * 123456, 789 * plus}
		if err := h.Set(i, obj); err != nil {
			t.Fatalf("append failed: %v", err)
		}
		want[i] = obj
	}
	if err := h.Close(); err != nil {
		t.Fatalf("bummer: %v", err)
	}
	t.Logf("iterate saved file:")

	sf, err := OpenSlice(testFile, false, &m)
	if err != nil {
		t.Fatalf("can't recycle: %v", err)
	}

	for i := 0; i < sf.Len(); i++ {
		var obj testType
		if err := sf.Get(i, &obj); err != nil {
			t.Fatalf("decode %d fail: %v", i, err)
		}
		if obj != want[i] {
			t.Fatalf("(%2d) want %v -- have %v", i, want[i], obj)
		}
	}
}
