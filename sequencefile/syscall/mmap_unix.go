// +build !windows,!plan9

package syscall

import "syscall"

func mmap(fd, length int, offset int64) ([]byte, error) {
	return syscall.Mmap(
		fd,
		offset,
		length,
		syscall.PROT_READ,
		syscall.MAP_SHARED,
	)
}
