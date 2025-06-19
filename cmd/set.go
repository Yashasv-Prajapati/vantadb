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

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set [filename] [key] [value]",
	Short: "Set a key to a value",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		
		filePath := args[0]
		key := args[1]
		value := args[2]

		disk, err := fs.Mount(filePath)
		if err != nil {
			fmt.Println("Mount failed:", err)
			return
		}
		kv.Init(disk)

		ok, err := kv.Set(key, value)
		if err != nil || !ok {
			fmt.Println("Set failed:", err)
		} else {
			fmt.Println("OK")
		}
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	
}
