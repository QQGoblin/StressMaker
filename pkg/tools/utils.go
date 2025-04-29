package tools

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/sets"
	"runtime"
	"strconv"
	"strings"
)

func ParseCPUs(sel []string) ([]int, error) {

	var (
		err              error
		result           = sets.NewInt()
		cpuBegin, cpuEnd int
	)

	numCPUs := runtime.NumCPU()

	for _, s := range sel {

		sd := strings.Split(s, "-")
		if len(sd) == 0 || len(sd) > 2 {
			return nil, fmt.Errorf("error cpu select %v", s)
		}

		cpuBegin, err = strconv.Atoi(sd[0])
		if err != nil {
			return nil, err
		}
		cpuEnd = cpuBegin
		if len(sd) == 2 {
			cpuEnd, err = strconv.Atoi(sd[0])
			if err != nil {
				return nil, err
			}
		}

		for cpu := cpuBegin; cpu < cpuEnd+1; cpu += 1 {
			if numCPUs > cpu {
				result.Insert(cpu)
			}
		}
	}
	return result.List(), nil
}
