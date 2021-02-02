package cli

import (
	"github.com/bwarren2/mandelbrot"
	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Generate a mandelbrot in a file",
	Long: `Mandelbrot generates images of the mandelbrot set.

It outputs in PNG, and is configurable for image size,
range, domain, and iterations`,
	Run: func(cmd *cobra.Command, args []string) {
		colors := mandelbrot.NewPalette(maxIterations)
		img := mandelbrot.MandelbrotBuilder{}.Draw(width, height, maxIterations, xMin, xMax, yMin, yMax, colors)
		mandelbrot.WritePng(img, filename)
	},
}

var xMin, xMax, yMin, yMax float64

func init() {
	newCmd.Flags().StringVarP(&filename, "filename", "f", "", "File to write to")
	newCmd.Flags().Uint16Var(&width, "width", 900, "Width of output image, in pixels")
	newCmd.Flags().Uint16Var(&height, "height", 450, "Height of output image, in pixels")
	newCmd.Flags().Uint8Var(&maxIterations, "maxIterations", 100, "Number of iterations to run the mandelbrot loop")
	newCmd.Flags().Float64Var(&xMin, "xMin", -2.5, "The min x value for the mandelbrot space to plot")
	newCmd.Flags().Float64Var(&xMax, "xMax", 1, "The max x value for the mandelbrot space to plot")
	newCmd.Flags().Float64Var(&yMin, "yMin", -1, "The min y value for the mandelbrot space to plot")
	newCmd.Flags().Float64Var(&yMax, "yMax", 1, "The max y value for the mandelbrot space to plot")
	rootCmd.AddCommand(newCmd)
}
