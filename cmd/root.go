package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{}

func init() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(teacherCommentRenderCmd)
	rootCmd.AddCommand(excelCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
