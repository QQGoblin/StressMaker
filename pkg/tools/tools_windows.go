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

func GetTickCount64() uint64 {
	ret, _, _ := procGetTickCount64.Call()
	return uint64(ret)
}

func SetThreadAffinity(idx int) error {

	var (
		threadHandle        = windows.CurrentThread()
		mask         uint64 = 1 << (idx - 1)
		err          error
		ret          uintptr
	)

	if ret, _, err = procSetThreadAffinityMask.Call(uintptr(threadHandle), uintptr(mask)); ret == 0 {
		return fmt.Errorf("SetThreadAffinityMask failed: %v", err)
	}

	return nil
}

func SetRealtimeScheduling() error {
	return nil
}
