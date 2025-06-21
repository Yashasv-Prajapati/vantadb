/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/Yashasv-Prajapati/vantadb/internal/fs"
	"github.com/Yashasv-Prajapati/vantadb/internal/kv"
	"github.com/Yashasv-Prajapati/vantadb/internal/wal"

	"github.com/spf13/cobra"
)

var recover bool

// walCmd represents the wal command
var walCmd = &cobra.Command{
	Use:   "wal",
	Short: "Get WAL logs",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		disk, err := fs.Mount("")
		if err != nil {
			fmt.Println("Mount failed:", err)
			return
		}
		kv.Init(disk)
		
		if recover {
			kv.RecoverFromLogs()
			return
		}

		fmt.Println(wal.GetAllWALRecords())
	},
}

func init() {
	rootCmd.AddCommand(walCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// walCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	walCmd.Flags().BoolVarP(&recover,"recover", "r", false, "Recover the DB using the WAL file")
	// serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")
}
