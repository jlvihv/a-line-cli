package cmd

import (
	_ "embed"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

//go:embed cicd.yml
var content string

var (
	yamlFile string
)

var rootCmd = &cobra.Command{
	Use: "github.com/hamster-shared/a-line-cli",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("start")

		Main(strings.NewReader(content))

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
