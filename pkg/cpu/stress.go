package cpustress

import (
	"context"
	"github.com/QQGoblin/StressMaker/pkg/tools"
	log "github.com/sirupsen/logrus"
	"runtime"
	"sync"
	"time"
)

const Million = 1000000

func StressOneCore(ctx context.Context, idx int, load func()) error {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if err := tools.SetThreadAffinity(idx); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			load()
		}
	}
}

func StressAllCore(ctx context.Context, load func(), targetCPUs []int, rt bool) {

	if rt {
		if err := tools.SetRealtimeScheduling(); err != nil {
			log.WithError(err).Fatal("Set process real-time scheduling")
		}
	}

	numCPUs := runtime.NumCPU()
	var wg sync.WaitGroup

	for _, t := range targetCPUs {
		if t > numCPUs {
			log.Warnf("CPU<%d> is not found", t)
			continue
		}

		log.Infof("Stress on CPU<%d>", t)
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			StressOneCore(ctx, idx, load)
		}(t)
	}

	wg.Wait()
}

// StaticStress 产生恒定负载，使 CPU 使用率控制特定百分比
func StaticStress(ctx context.Context, load float64, targetCPUs []int, rt bool) {
	loadInMillSecond := uint64(load * 1000)

	if loadInMillSecond > 990 {
		loadInMillSecond = 990
	}

	log.Infof("CPU load %vms\n", loadInMillSecond)

	StressAllCore(ctx, func() {
		startTime := tools.GetTickCount64()
		for tools.GetTickCount64()-startTime < loadInMillSecond {
		}

		time.Sleep(time.Duration(1000-loadInMillSecond) * time.Millisecond)
	}, targetCPUs, rt)
}

// CalcStress 指定四则运算负载
func CalcStress(ctx context.Context, loop int64, targetCPUs []int, rt bool) {

	log.Infof("CPU loop %v\n", loop)

	StressAllCore(ctx, func() {
		startTime := tools.GetTickCount64()
		value := loop * Million
		for value > 0 {
			value = value - 1
		}
		cost := tools.GetTickCount64() - startTime
		if cost < 1000 {
			time.Sleep(time.Duration(1000-cost) * time.Millisecond)
		}
	}, targetCPUs, rt)
}
