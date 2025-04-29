package freq

import (
	"github.com/QQGoblin/StressMaker/pkg/freq"
	"github.com/spf13/cobra"
)

var (
	selectCPUs []string
	offlineCPU bool
)

func init() {
	Command.PersistentFlags().StringSliceVarP(&selectCPUs, "cpu", "c", []string{"1"}, "")
	onlineCommand.PersistentFlags().BoolVarP(&offlineCPU, "offline", "s", false, "")
	Command.AddCommand(onlineCommand)
}

var Command = &cobra.Command{
	Use:   "freq",
	Short: "",
}

var onlineCommand = &cobra.Command{
	Use:   "cpu-online",
	Short: "",
	RunE: func(cmd *cobra.Command, args []string) error {
		return freq.OnlineCPUs(selectCPUs, !offlineCPU)
	},
}
