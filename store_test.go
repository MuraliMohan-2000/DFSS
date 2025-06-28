package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPathTRansforFunc(t *testing.T) {
	key := "momsbestpicture"
	pathname := CASPathTRansformFunc(key)
	fmt.Println(pathname)
	expectedPathName := "68044/29f74/181a6/3c50c/3d81d/733a1/2f14a/353ff"
	if pathname != expectedPathName {
		t.Errorf("have %s want %s", pathname, expectedPathName)
	}
}

func TestStore(t *testing.T) {
	opts := storeOpts{
		PathTransformFunc: CASPathTRansformFunc,
	}

	s := NewStore(opts)
	data := bytes.NewReader([]byte("some jpg bytes"))
	if err := s.writeStream("myspecialpicture", data); err != nil {
		t.Error(err)
	}
}
