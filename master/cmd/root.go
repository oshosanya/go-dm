package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-dm",
	Short: "go-dm is a download manager written in golang",
	Long: `A fast and flexible download manager written in golang
	You just have to love it`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	// Do Stuff Here
	// },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
