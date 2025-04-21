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
	"golang.org/x/sys/unix"
	"syscall"
	"time"
)

func GetTickCount64() uint64 {
	return uint64(time.Now().UnixMilli())
}

func SetThreadAffinity(idx int) error {

	var (
		err error
	)

	target := unix.CPUSet{}
	target.Set(idx)

	if err = unix.SchedSetaffinity(syscall.Gettid(), &target); err != nil {
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
