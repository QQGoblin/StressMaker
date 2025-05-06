//go:build windows

package tools

import (
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

var (
	kernel32                  = syscall.NewLazyDLL("kernel32.dll")
	procSetThreadAffinityMask = kernel32.NewProc("SetThreadAffinityMask")
	procGetTickCount64        = kernel32.NewProc("GetTickCount64")
	user32                    = syscall.NewLazyDLL("user32.dll")
	enumWindowsProc           = user32.NewProc("EnumWindows")
	getWindowTextW            = user32.NewProc("GetWindowTextW")
	getWindowTextLen          = user32.NewProc("GetWindowTextLengthW")
	isWindowVisible           = user32.NewProc("IsWindowVisible")
	getWindowThreadProcID     = user32.NewProc("GetWindowThreadProcessId")
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

type EnumWindowsProc func(hwnd syscall.Handle, lParam uintptr) uintptr

// enumWindows 是对 EnumWindows 的封装
func enumWindows(callback func(hwnd syscall.Handle, lParam uintptr) uintptr, lParam uintptr) bool {
	ret, _, _ := enumWindowsProc.Call(
		syscall.NewCallback(callback),
		lParam,
	)
	return ret != 0
}

func GetWindowsByProcess(targetPID uint32) ([]string, error) {
	var windowTitles []string

	// 调用 EnumWindows 函数
	success := enumWindows(func(hwnd syscall.Handle, lParam uintptr) uintptr {

		// 获取窗口所属的进程 ID
		var processID uint32
		getWindowThreadProcID.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&processID)))
		if processID != targetPID {
			return 1 // 跳过不属于目标进程的窗口
		}

		// 获取窗口标题长度
		textLen, _, _ := getWindowTextLen.Call(uintptr(hwnd))
		if textLen == 0 {
			return 1 // 跳过没有标题的窗口
		}

		// 获取窗口标题
		buffer := make([]uint16, textLen+1)
		getWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer)))

		// 转换为 Go 字符串并保存
		windowTitle := syscall.UTF16ToString(buffer)
		windowTitles = append(windowTitles, windowTitle)

		return 1 // 继续枚举
	}, 0)

	if !success {
		return nil, fmt.Errorf("failed to enumerate windows")
	}

	return windowTitles, nil
}
