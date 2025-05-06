package excel

import (
	"fmt"
	"github.com/QQGoblin/StressMaker/pkg/tools"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
	"time"
)

func excelEXE() (string, error) {

	if _, err := os.Stat("C:\\Program Files (x86)\\Microsoft Office\\root\\Office16\\EXCEL.EXE"); err == nil {
		return "C:\\Program Files (x86)\\Microsoft Office\\root\\Office16\\EXCEL.EXE", nil
	}

	if _, err := os.Stat("C:\\Program Files\\Microsoft Office\\Office16\\EXCEL.EXE"); err == nil {
		return "C:\\Program Files\\Microsoft Office\\Office16\\EXCEL.EXE", nil
	}

	return "", fmt.Errorf("EXCEL.EXE is not found")
}

func openExcelFile(filePath, bin string) (*os.Process, error) {

	var (
		err error
	)

	if bin == "" {
		if bin, err = excelEXE(); err != nil {
			return nil, err
		}
	}

	cmd := exec.Command(bin, filePath)
	if err = cmd.Start(); err != nil {
		return nil, err
	}

	// 获取进程 ID
	if cmd.Process == nil {
		return nil, fmt.Errorf("pid is not found")
	}

	return cmd.Process, nil
}

func BootCost(filePath, binPath string) error {

	var (
		err         error
		excelProc   *os.Process
		openWindows []string
		startTime   = time.Now()
	)

	if excelProc, err = openExcelFile(filePath, binPath); err != nil {
		return err
	}

	defer func() {
		if excelProc != nil {
			excelProc.Release()
		}
	}()

	for {

		isOpenMainWindow := false
		isOpenLoadWindow := false

		if openWindows, err = tools.GetWindowsByProcess(uint32(excelProc.Pid)); err != nil {
			return err
		}

		if len(openWindows) > 0 {
			for _, winStr := range openWindows {

				if strings.HasPrefix(winStr, filePath) {
					isOpenMainWindow = true
				}
				if strings.HasPrefix(winStr, "正在打开") {
					isOpenLoadWindow = true
				}
			}

			if isOpenMainWindow && !isOpenLoadWindow {
				break
			}

		}
		time.Sleep(time.Millisecond * 100)
	}

	bootCost := time.Now().Sub(startTime)
	log.Infof("Boot excel file cost: %fs", bootCost.Seconds())

	excelProc.Wait()

	return nil
}
