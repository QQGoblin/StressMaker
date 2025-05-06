package excel

import (
	"github.com/QQGoblin/StressMaker/pkg/excel"
	"github.com/spf13/cobra"
)

var (
	bootFile string
	binFile  string
)

func init() {
	Command.PersistentFlags().StringVarP(&bootFile, "file", "f", "input.xlsx", "")
	Command.PersistentFlags().StringVarP(&binFile, "bin", "b", "", "")
}

var Command = &cobra.Command{
	Use:   "excel",
	Short: "Get excel boot cost",
	RunE: func(cmd *cobra.Command, args []string) error {
		return excel.BootCost(bootFile, binFile)
	},
}
