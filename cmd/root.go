package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)
var rootCmd = &cobra.Command {
	Use:   "go-mitm",
	Short: "A basic man-in-the-middle proxy for debugging outgoing requests",
	Run: func(cmd *cobra.Command, args []string) { 
		mitmMain()
	},
}

var (
	port string
)

func init() {
	rootCmd.Flags().StringVarP(&port, "port","p", "8080", "Port to listen on")
}

// Framework boilerplate
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
