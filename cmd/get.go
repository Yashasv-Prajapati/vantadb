/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"vantadb/internal/fs"
	"vantadb/internal/kv"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [filename] [key]",
	Short: "Gets the value to a key",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		key := args[1]

		disk, err := fs.Mount(filePath)
		if err != nil {
			fmt.Println("Mount failed:", err)
			return
		}
		kv.Init(disk)

		value, err := kv.Get(key)
		if err != nil {
			fmt.Println("Get failed:", err)
		} else {
			fmt.Println(value)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
