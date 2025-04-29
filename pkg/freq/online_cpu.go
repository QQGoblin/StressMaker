package freq

import (
	"fmt"
	"github.com/QQGoblin/StressMaker/pkg/tools"
	log "github.com/sirupsen/logrus"
	"os"
)

func onlineCPU(cpu int, online bool) error {

	var (
		err            error
		ctrl           *os.File
		onlineFilePath = fmt.Sprintf("/sys/devices/system/cpu/cpu%d/online", cpu)
	)

	if ctrl, err = os.OpenFile(onlineFilePath, os.O_WRONLY, 0644); err != nil {
		return err
	}
	defer ctrl.Close()

	onlineValue := "1"
	if !online {
		onlineValue = "0"
	}

	_, err = ctrl.WriteString(onlineValue)
	return err
}

func OnlineCPUs(selectCPUs []string, online bool) error {

	var (
		err    error
		target []int
	)

	if target, err = tools.ParseCPUs(selectCPUs); err != nil {
		return err
	}

	for _, cpu := range target {

		if (cpu == 0 || cpu == 1) && !online {
			log.Warn("Never offline CPU0 or CPU1")
			continue
		}

		if err = onlineCPU(cpu, online); err != nil {
			return err
		}
	}

	return nil
}
