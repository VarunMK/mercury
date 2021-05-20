package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func getHandler(path string, s *http.Server) func(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open(path)
	if err != nil {
		return func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("Error: %v\n", err)
			w.WriteHeader(http.StatusNotFound)
			s.Shutdown(context.Background())
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		fstat, err := f.Stat()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			s.Shutdown(context.Background())
		}
		var sz = fstat.Size()
		w.Header().Add("Content-Length", strconv.Itoa(int(sz)))
		bar := progressbar.DefaultBytes(
			sz,
			"uploading",
		)
		_, err = io.Copy(io.MultiWriter(w, bar), f)
		if err != nil {
			fmt.Printf("Unsuccessful request from %v, Error: %v\n", r.RemoteAddr, err)
			s.Shutdown(context.Background())
		}
		fmt.Printf("Successful request from %v\n", r.RemoteAddr)
	}
}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send file to other device",
	Long: `
	This command allows you to send a file to another device.
	To run: 
		mercury send [filename] 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Error: One arg should be provided, the path")
			return
		}
		var path string = args[0]
		m := http.NewServeMux()
		s := http.Server{Addr: ":3000", Handler: m}
		m.HandleFunc("/", getHandler(path, &s))
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
}
