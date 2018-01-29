package cmd

import (
	"github.com/spf13/cobra"
	"monkey/bootstrap"
)

var serveCmd = &cobra.Command{
	Use:"serve",
	Short:"Starts monkey server",
	Long:"",
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap.Server()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}