package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTRansforFunc(t *testing.T) {
	key := "momsbestpicture"
	pathkey := CASPathTRansformFunc(key)
	expectedFileName := "6804429f74181a63c50c3d81d733a12f14a353ff"
	expectedPathName := "68044/29f74/181a6/3c50c/3d81d/733a1/2f14a/353ff"
	if pathkey.PathName != expectedPathName {
		t.Errorf("have %s want %s", pathkey.PathName, expectedPathName)
	}

	if pathkey.FileName != expectedFileName {
		t.Errorf("have %s want %s", pathkey.FileName, expectedFileName)
	}
}

func TestStore(t *testing.T) {
	s := newStore()
	id := generateID()
	defer tearDown(t, s)

	for i := 0; i < 50; i++ {

		key := fmt.Sprintf("foo_%d", i)
		data := []byte("some jpg bytes")

		if _, err := s.writeStream(id, key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		if ok := s.Has(id, key); !ok {
			t.Errorf("expected to have key %s", key)
		}

		_, r, err := s.Read(id, key)
		if err != nil {
			t.Error(err)
		}

		b, _ := io.ReadAll(r)

		fmt.Println(string(b))

		if string(b) != string(data) {
			t.Errorf("want %s have %s", data, b)
		}

		err = s.Delete(id, key)
		if err != nil {
			t.Errorf("Deletion Failed for key %s", key)
		}

		if ok := s.Has(id, key); ok {
			t.Errorf("expected to not have the key %s", key)
		}
	}

}

func newStore() *store {

	opts := storeOpts{
		PathTransformFunc: CASPathTRansformFunc,
	}

	return NewStore(opts)
}

func tearDown(t *testing.T, s *store) {
	if err := s.clear(); err != nil {
		t.Error(err)

	}
}
