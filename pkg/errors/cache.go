package errors

const (
	CacheNotFound ErrorType = 200 + iota
	CacheGetError
	CacheSetError
	CacheDeleteError
	CacheKeysError
)
