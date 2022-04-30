//go:build windows
// +build windows

package fstat

import (
	"strconv"
	"syscall"
)

func FStat(filename string) string {
	var iNo uint64
	var fi syscall.ByHandleFileInformation

	h, e := syscall.Open(filename, syscall.FILE_SHARE_READ|syscall.O_RDONLY|syscall.OPEN_EXISTING, 0777)

	if e = syscall.GetFileInformationByHandle(h, &fi); e != nil {
		syscall.CloseHandle(h)
	}

	iNo = uint64(fi.FileIndexHigh)<<32 | uint64(fi.FileIndexLow)

	return strconv.FormatUint(iNo, 10)
}
