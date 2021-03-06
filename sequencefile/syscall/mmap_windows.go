// Copyright 2017 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package syscall

import (
	"os"
	"syscall"
	"unsafe"
)

// func mmap(fd, size int) ([]byte, error) {
// 	low, high := uint32(size), uint32(size>>32)
// 	h, errno := syscall.CreateFileMapping(syscall.Handle(fd), nil, syscall.PAGE_READONLY, high, low, nil)
// 	if h == 0 {
// 		return nil, os.NewSyscallError("CreateFileMapping", errno)
// 	}

// 	addr, errno := syscall.MapViewOfFile(h, syscall.FILE_MAP_READ, 0, 0, uintptr(size))
// 	if addr == 0 {
// 		return nil, os.NewSyscallError("MapViewOfFile", errno)
// 	}

// 	if err := syscall.CloseHandle(syscall.Handle(h)); err != nil {
// 		return nil, os.NewSyscallError("CloseHandle", err)
// 	}

// 	return (*[maxMapSize]byte)(unsafe.Pointer(addr))[:size], nil
// }

func mmap(fd, length int, offset int64) ([]byte, error) {
	low := uint32((offset + int64(length)) & 0xFFFFFFFF)
	high := uint32((offset + int64(length)) >> 32)
	h, errno := syscall.CreateFileMapping(syscall.Handle(fd), nil, syscall.PAGE_READONLY, high, low, nil)
	if h == 0 {
		return nil, os.NewSyscallError("CreateFileMapping", errno)
	}

	offsetLow := uint32(offset & 0xFFFFFFFF)
	offsetHigh := uint32(offset >> 32)
	addr, errno := syscall.MapViewOfFile(h, syscall.FILE_MAP_READ, offsetHigh, offsetLow, uintptr(length))
	if addr == 0 {
		return nil, os.NewSyscallError("MapViewOfFile", errno)
	}

	if err := syscall.CloseHandle(syscall.Handle(h)); err != nil {
		return nil, os.NewSyscallError("CloseHandle", err)
	}

	return (*[maxMapSize]byte)(unsafe.Pointer(addr))[offset : offset+int64(length)], nil
}