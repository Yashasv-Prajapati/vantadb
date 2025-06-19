/*
Copyright © 2025 NAME HERE yamuprajapati05@gmail.com
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"vantadb/internal/fs"
	"vantadb/internal/kv"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
)

var port int
var filePath string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the vantadb server",
	Run: func(cmd *cobra.Command, args []string) {
		disk, err := fs.Mount(filePath)
		if err != nil {
			log.Fatalf("Failed to mount disk: %v", err)
		}

		kv.Init(disk)

		http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			if key == "" {
				http.Error(w, "Missing key", http.StatusBadRequest)
				return
			}
			val, err := kv.Get(key)
			if err != nil {
				http.Error(w, "Key not found", http.StatusNotFound)
				return
			}
			fmt.Fprint(w, val)
		})

		http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
			var payload struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
			ok, err := kv.Set(payload.Key, payload.Value)
			if err != nil || !ok {
				http.Error(w, "Failed to set value", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		})

		go func() {
			addr := fmt.Sprintf(":%d", port)
			fmt.Println("Serving on http://localhost" + addr)
			log.Fatal(http.ListenAndServe(addr, nil))
		}()

		rl, err := readline.New("> ")
		if err != nil {
			log.Fatalf("failed to start repl: %v", err)
		}
		defer rl.Close()

		for {
			line, err := rl.Readline()
			if err != nil {
				break
			}

			input := strings.TrimSpace(line)

			if input == "" {
				continue
			}

			parts := strings.SplitN(input, " ", 3)
			cmd := parts[0]

			switch cmd {
			case "exit":
				fmt.Printf("shutting down")
				return

			case "get":
				if len(parts) != 2 {
					fmt.Printf("invalid number of arguments")
					return
				}
				key := parts[1]
				value, err := kv.Get(key)
				if err != nil {
					fmt.Printf("could not get key value: %v", err)
					continue
				}
				fmt.Println(value)

			case "set":
				if len(parts) != 3 {
					fmt.Printf("invalid number of arguments")
					return
				}
				key := parts[1]
				value := parts[2]
				check, err := kv.Set(key, value)
				if err != nil {
					fmt.Printf("could not set key value: %v\n", err)
					continue
				}
				if !check {
					fmt.Printf("failed to set key value")
					continue
				}
				fmt.Println("OK")

			default:
				fmt.Println("unknown command")
			}

		}

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")
	serveCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the .vdsk file")
	serveCmd.MarkFlagRequired("file")
}
