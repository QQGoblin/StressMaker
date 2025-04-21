package cpustress

import (
	"context"
	cpustress "github.com/QQGoblin/StressMaker/pkg/cpu"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	staticLoad float64
	calcLoad   int64
	targetCPUs []int
	rt         bool
	all        bool
)

func init() {
	StaticCommand.PersistentFlags().Float64VarP(&staticLoad, "static-load", "l", 0.7, "")
	CalcCommand.PersistentFlags().Int64VarP(&calcLoad, "calc-load", "c", 1000, "")
	Command.PersistentFlags().IntSliceVar(&targetCPUs, "cpu", []int{1}, "")
	Command.PersistentFlags().BoolVarP(&rt, "real-time", "rt", false, "")
	Command.PersistentFlags().BoolVarP(&all, "all", "a", false, "")
	Command.AddCommand(StaticCommand)
	Command.AddCommand(CalcCommand)
}

var Command = &cobra.Command{
	Use:   "cpu",
	Short: "Create CPU stress",
}

func bindAllCPUs() []int {
	var (
		result = make([]int, 0)
	)
	for i := 1; i <= runtime.NumCPU(); i++ {
		result = append(result, i)
	}
	return result
}

var StaticCommand = &cobra.Command{
	Use:   "static",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx, cancel := context.WithCancel(context.Background())

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigChan
			log.Infof("Received interrupt signal. Shutting down...")
			cancel()
		}()

		if all {
			targetCPUs = bindAllCPUs()
		}

		cpustress.StaticStress(ctx, staticLoad, targetCPUs, rt)

		return nil
	},
}

var CalcCommand = &cobra.Command{
	Use:   "calc",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {

		ctx, cancel := context.WithCancel(context.Background())

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigChan
			log.Infof("Received interrupt signal. Shutting down...")
			cancel()
		}()

		if all {
			targetCPUs = bindAllCPUs()
		}

		cpustress.CalcStress(ctx, calcLoad, targetCPUs, rt)

		return nil
	},
}
