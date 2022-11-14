package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	yamlFile string
)

var rootCmd = &cobra.Command{
	Use: "github.com/hamster-shared/a-line-cli",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("start")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&yamlFile, "file", "f", "cicd.yaml", "yaml config file")
}
