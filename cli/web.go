package cli

import (
	"github.com/bwarren2/mandelbrot/web"
	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the mandelbrot webserver",
	Long:  "Start the mandelbrot webserver",
	Run: func(cmd *cobra.Command, args []string) {
		web.Serve(7777)
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
}
