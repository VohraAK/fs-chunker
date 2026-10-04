/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fs-chunker/internal/logger"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fs-chunker",
	Short: "A concurrent file indexer and chunker",
	SilenceUsage:  true,
	SilenceErrors:  true,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.Init()
	},
}


// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		// exec failure path
		logger.Init()
			
		slog.Error("CLI execution failed",
			slog.String("error", err.Error()),
		)
		
		os.Exit(1)
	}
}


// init function defines the flags and config settings
func init() {	
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
