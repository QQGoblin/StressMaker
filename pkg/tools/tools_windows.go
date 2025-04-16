//go:build windows

package tools

import (
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
)

var (
	kernel32                  = syscall.NewLazyDLL("kernel32.dll")
	procSetThreadAffinityMask = kernel32.NewProc("SetThreadAffinityMask")
	procGetTickCount64        = kernel32.NewProc("GetTickCount64")
)

// GetTickCount64 获取系统启动后的毫秒数
func GetTickCount64() uint64 {
	ret, _, _ := procGetTickCount64.Call()
	return uint64(ret)
}

// SetThreadAffinity 设置线程的 CPU 亲和性
func SetThreadAffinity(idx int) error {

	threadHandle := windows.CurrentThread()

	var mask uint64 = 1 << (idx - 1)
	r1, _, err := procSetThreadAffinityMask.Call(
		uintptr(threadHandle),
		uintptr(mask),
	)
	if r1 == 0 {
		return fmt.Errorf("SetThreadAffinityMask failed: %v", err)
	}

	return nil
}
