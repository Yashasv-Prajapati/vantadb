/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"vantadb/internal/fs"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the vanta db database",
	Long: `Initializes the vanta db database and file with your given name, if you don't give a name, it used .vdsk by default`,
	Run: func(cmd *cobra.Command, args []string) {
		err := fs.CreateVDSKStorageData()
		if err != nil {
			fmt.Printf("Failed to create disk: %v\n", err)
			return
		}
		fmt.Println("Disk created:", args[0])
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
