//go:build aix || freebsd || linux || netbsd

package tools

import (
	"syscall"
	"time"
	"unsafe"
)

const (
	_SysSchedSetAffinity = 203
	_CPUSetSize          = 1024
)

type _CPUSet struct {
	bits [_CPUSetSize / 64]uint64
}

func GetTickCount64() uint64 {
	return uint64(time.Now().UnixMilli())
}

func SetThreadAffinity(idx int) error {

	var (
		cpuSet _CPUSet
		err    syscall.Errno
	)

	cpuSet.bits[idx/64] |= 1 << (uint(idx) % 64)

	if _, _, err = syscall.Syscall(
		_SysSchedSetAffinity,
		uintptr(syscall.Gettid()),
		uintptr(unsafe.Sizeof(cpuSet)),
		uintptr(unsafe.Pointer(&cpuSet)),
	); err != 0 {
		return err
	}

	return nil
}
