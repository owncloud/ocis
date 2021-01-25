package cmd

import (
	"github.com/owncloud/ocis/ocis/pkg/runtime/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "pman",
		Short: "RPC Process Manager",
	}
)

// RootCmd returns a configured root command.
func RootCmd(cfg *config.Config) *cobra.Command {
	rootCmd.PersistentFlags().StringVarP(&cfg.Hostname, "hostname", "n", "localhost", "host with a running OCIS runtime.")
	rootCmd.PersistentFlags().StringVarP(&cfg.Port, "port", "p", "10666", "port to send messages to the rpc OCIS runtime.")
	rootCmd.PersistentFlags().BoolVarP(&cfg.KeepAlive, "keep-alive", "k", false, "restart supervised processes that abruptly die.")

	viper.BindPFlag("hostname", rootCmd.PersistentFlags().Lookup("hostname"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))

	rootCmd.AddCommand(List(cfg))
	rootCmd.AddCommand(Run(cfg))
	rootCmd.AddCommand(Kill(cfg))

	return rootCmd
}
