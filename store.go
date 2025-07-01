package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultRootFolderName = "mnetwork"

func CASPathTRansformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key)) //[20]byte => []byte => [:]
	hashstr := hex.EncodeToString(hash[:])

	blocksize := 5
	sliceLen := len(hashstr) / blocksize

	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashstr[from:to]
	}

	return PathKey{
		PathName: strings.Join(paths, "/"),
		FileName: hashstr,
	}
}

type PathTransformFunc func(string) PathKey

type PathKey struct {
	PathName string
	FileName string
}

func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0 {
		return ""
	}

	return paths[0]
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.FileName)
}

type storeOpts struct {
	//Root is the folder name of the root, containing all the folders/files
	//of the system
	Root string

	PathTransformFunc PathTransformFunc
}

var DefaultpathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		FileName: key,
	}
}

type store struct {
	storeOpts
}

func NewStore(opts storeOpts) *store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultpathTransformFunc
	}

	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
	}

	return &store{
		storeOpts: opts,
	}
}

func (s *store) Has(id string, key string) bool {
	pathkey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathkey.FullPath())

	_, err := os.Stat(fullPathWithRoot)

	return !errors.Is(err, os.ErrNotExist)
}

func (s *store) clear() error {
	return os.RemoveAll(s.Root)
}

func (s *store) Delete(id string, key string) error {
	pathkey := s.PathTransformFunc(key)

	defer func() {
		log.Printf("deleted [%s] from disk", pathkey.FileName)
	}()

	firstPathNameWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathkey.FirstPathName())

	return os.RemoveAll(firstPathNameWithRoot)

}

func (s *store) Write(id string, key string, r io.Reader) (int64, error) {
	return s.writeStream(id, key, r)
}

func (s *store) WriteDecrypt(encKey []byte, id string, key string, r io.Reader) (int64, error) {
	f, err := s.openFileForWritting(id, key)
	if err != nil {
		return 0, err
	}

	n, err := copyDecrypt(encKey, r, f)

	return int64(n), err
}

func (s *store) openFileForWritting(id string, key string) (*os.File, error) {
	pathkey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathkey.PathName)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return nil, err
	}

	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathkey.FullPath())
	return os.Create(fullPathWithRoot)

}

func (s *store) writeStream(id string, key string, r io.Reader) (int64, error) {
	f, err := s.openFileForWritting(id, key)
	if err != nil {
		return 0, err
	}
	return io.Copy(f, r)
}

func (s *store) Read(id string, key string) (int64, io.Reader, error) {
	return s.readstream(id, key)
}

func (s *store) readstream(id string, key string) (int64, io.ReadCloser, error) {
	pathkey := s.PathTransformFunc(key)
	fullPathWIthRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathkey.FullPath())

	file, err := os.Open(fullPathWIthRoot)
	if err != nil {
		return 0, nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return 0, nil, err
	}

	return fi.Size(), file, nil
}
