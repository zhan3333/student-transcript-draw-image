package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{}

func init() {
	//cobra.OnInitialize(initConfig) // viper是cobra集成的配置文件读取的库，以后我们会专门说
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(excelCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
