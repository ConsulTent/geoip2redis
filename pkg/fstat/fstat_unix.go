//go:build !windows
// +build !windows

package fstat

import (
	"strconv"
	"syscall"
)

func FStat(fi string) string {
	var iNo uint64 = 0
	var stat syscall.Stat_t

	if err := syscall.Stat(fi, &stat); err == nil {
		iNo = stat.Ino
	}

	return strconv.FormatUint(iNo, 10)
}
