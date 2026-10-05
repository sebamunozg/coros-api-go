package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "coros",
	Short: "Coros CLI is a command-line interface for interacting with the Coros API.",
	Long:  `Coros CLI allows you to authenticate and download activities from the Coros API.`,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
