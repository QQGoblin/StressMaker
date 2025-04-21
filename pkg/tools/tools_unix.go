//go:build aix || freebsd || linux || netbsd

package tools

/*
#include <sched.h>
#include <sys/types.h>
#include <unistd.h>

int set_realtime_scheduling(pid_t pid) {
    struct sched_param param;
    param.sched_priority = 99;
    return sched_setscheduler(pid, SCHED_FIFO, &param);
}
*/
import "C"

import (
	"fmt"
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

func SetRealtimeScheduling() error {
	if ret := C.set_realtime_scheduling(C.getpid()); ret != 0 {
		return fmt.Errorf("failed to set realtime scheduling, ret %v", ret)
	}
	return nil
}
