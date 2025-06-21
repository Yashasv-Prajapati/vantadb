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

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a key to a value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		
		key := args[1]
		value := args[2]

		disk, err := fs.Mount("")
		if err != nil {
			fmt.Println("Mount failed:", err)
			return
		}
		kv.Init(disk)

		fmt.Println(kv.Set(key, value))
		// if err != nil || !ok {
		// 	fmt.Println("Set failed:", err)
		// } else {
		// 	fmt.Println("OK")
		// }
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
