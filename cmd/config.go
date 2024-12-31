/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Isaac-Fate/myst/internal/config"
	"github.com/Isaac-Fate/myst/internal/utils"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Set the password store and the passphrase",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")

		var c config.Config

		passwordStorePath, err := cmd.Flags().GetString("store")
		cobra.CheckErr(err)

		passwordStorePath, err = utils.ResolvePath(passwordStorePath)
		cobra.CheckErr(err)

		// Set the password store path
		c.PasswordStorePath = passwordStorePath

		err = c.Save("./myst.yaml")
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	configCmd.Flags().StringP("store", "s", "", "Filepath of the password store database")
}
