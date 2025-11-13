package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "builder",
	Short: "Universal builder for security platform agents and packages",
	Long: `Builder is a CLI tool for generating installable agents, SDKs, and packages
across multiple platforms: web frontends, backends, desktop, mobile, and network appliances.

It supports:
- Agent initialization and scaffolding
- Multi-language agent builds
- Cross-platform package generation (deb/rpm/msi/dmg/helm)
- Code signing and artifact verification
- Air-gapped bundle generation`,
	Version: "0.1.0",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.builder.yaml)")
	rootCmd.PersistentFlags().String("workspace", ".", "workspace directory")
	viper.BindPFlag("workspace", rootCmd.PersistentFlags().Lookup("workspace"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".builder")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
