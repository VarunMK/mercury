package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var rcvCmd = &cobra.Command{
	Use:   "receive",
	Short: "Receive a file from another device",
	Long: `
	This command allows you to receive a file from another device.
	To run: 
		mercury receive [endpoint] [file name to save as] 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("Error: 2 args should be provided, the hostname and the downloaded file name")
			return
		}
		var hostname = args[0]
		var fname = args[1]
		resp, err := http.Get(hostname)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		f, err := os.Create(fname)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		bar := progressbar.DefaultBytes(
			resp.ContentLength,
			"downloading",
		)
		io.Copy(io.MultiWriter(f, bar), resp.Body)
	},
}

func init() {
	rootCmd.AddCommand(rcvCmd)
}
