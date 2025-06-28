package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)

func CASPathTRansformFunc(key string) string {
	hash := sha1.Sum([]byte(key)) //[20]byte => []byte => [:]
	hashstr := hex.EncodeToString(hash[:])

	blocksize := 5
	sliceLen := len(hashstr) / blocksize

	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashstr[from:to]
	}

	return strings.Join(paths, "/")
}

type PathTransformFunc func(string) string

type storeOpts struct {
	PathTransformFunc PathTransformFunc
}

var DefaultpathTransformFunc = func(key string) string {
	return key
}

type store struct {
	storeOpts
}

func NewStore(opts storeOpts) *store {
	return &store{
		storeOpts: opts,
	}
}
func (s *store) writeStream(key string, r io.Reader) error {
	pathName := s.PathTransformFunc(key)

	if err := os.MkdirAll(pathName, os.ModePerm); err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	io.Copy(buf, r)

	fileNameBytes := md5.Sum(buf.Bytes())
	filename := hex.EncodeToString(fileNameBytes[:])
	pathAndFilename := pathName + "/" + filename

	f, err := os.Create(pathAndFilename)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, buf)
	if err != nil {
		return err
	}

	log.Printf("Written (%d) bytes to disk %s", n, pathAndFilename)

	return nil
}
