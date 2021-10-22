package syscall

func Mmap(fd, length int) ([]byte, error) {
	return mmap(fd, length, 0)
}

func MmapOffset(fd, length int, offset int64) ([]byte, error) {
	return mmap(fd, length, offset)
}
