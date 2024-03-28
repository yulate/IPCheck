package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "IPCheck",
	Short: "IPCheck IP检查工具",
	Long: `
当前功能
	给出ip.txt 自动去重
	按照low midden high 三个级别进行端口扫描
	可自定义设置线程
`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
