package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"solana-program-scanner/config"
	"solana-program-scanner/factory"
	"solana-program-scanner/log"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "solana-program-scanner",
		Short: "solana raydium amm program scanner",
	}

	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "Configuration related commands",
	}

	var configExportCmd = &cobra.Command{
		Use:   "export",
		Short: "Export configuration",
		Run: func(cmd *cobra.Command, args []string) {
			if cmd.Flags().Changed("input") {
				config.LoadConfig(config.F)
			}
			config.SaveConfig(config.FT)
		},
	}

	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run the application with the specified configuration file",
		Run: func(cmd *cobra.Command, args []string) {
			watchSignal()
			config.LoadConfig(config.F)
			f := factory.Factory{}
			f.Assemble().Run()
			log.Logger.Sync()
		},
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of the application",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(GetVersion())
		},
	}

	configExportCmd.Flags().StringVarP(&config.F, "input", "i", "config.json", "input configuration file to load")
	configExportCmd.Flags().StringVarP(&config.FT, "output", "o", "config.json.template", "output path for the exported configuration")
	runCmd.Flags().StringVarP(&config.F, "config", "c", "config.json", "configuration file path")

	configCmd.AddCommand(configExportCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
