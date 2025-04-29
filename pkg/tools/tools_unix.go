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
	"k8s.io/apimachinery/pkg/util/sets"
	"os"
	"strconv"
	"strings"
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

func CPUNum() (int, error) {

	var (
		err    error
		files  []os.DirEntry
		result int
	)

	if files, err = os.ReadDir("/sys/devices/system/cpu/"); err != nil {
		return 0, err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "cpu") {
			if _, err = strconv.Atoi(file.Name()[3:]); err == nil {
				result++
			}
		}
	}
	return result, nil
}

func ParseCPUs(sel []string) ([]int, error) {

	var (
		result   = sets.NewInt()
		totalCPU int
		err      error
	)

	if totalCPU, err = CPUNum(); err != nil {
		return nil, err
	}

	for _, s := range sel {
		var (
			cpuBegin, cpuEnd int
			sd               = strings.Split(s, "-")
		)

		if len(sd) == 0 || len(sd) > 2 {
			return nil, fmt.Errorf("error cpu select %v", s)
		}

		if cpuBegin, err = strconv.Atoi(sd[0]); err != nil {
			return nil, err
		}

		cpuEnd = cpuBegin

		if len(sd) == 2 {
			if cpuEnd, err = strconv.Atoi(sd[1]); err != nil {
				return nil, err
			}
		}

		for cpu := cpuBegin; cpu < cpuEnd+1; cpu += 1 {
			if totalCPU > cpu {
				result.Insert(cpu)
			}
		}
	}
	return result.List(), nil
}
