package objectstorage

import "io"

type ObjectStorage interface {
	Put(collection, filename string, reader io.Reader, size int64, contentType string) (int64, error)
	PutImage(collection, filename string, reader io.Reader, size int64) (int64, error)
}
