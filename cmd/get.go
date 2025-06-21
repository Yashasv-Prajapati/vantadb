/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/Yashasv-Prajapati/vantadb/internal/fs"
	"github.com/Yashasv-Prajapati/vantadb/internal/kv"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Gets the value to a key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[1]

		disk, err := fs.Mount("")
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
