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

// 定义 cSet 结构体
type _CPUSet struct {
	bits [_CPUSetSize / 64]uint64
}

// GetTickCount64 获取系统启动后的毫秒数
func GetTickCount64() uint64 {
	return uint64(time.Now().UnixMilli())
}

// SetThreadAffinity 设置线程的 CPU 亲和性
func SetThreadAffinity(idx int) error {

	var cpuSet _CPUSet
	cpuSet.bits[idx/64] |= 1 << (uint(idx) % 64)
	_, _, err := syscall.Syscall(
		_SysSchedSetAffinity,
		uintptr(syscall.Gettid()),
		uintptr(unsafe.Sizeof(cpuSet)),
		uintptr(unsafe.Pointer(&cpuSet)),
	)

	if err != 0 {
		return err
	}
	return nil
}
