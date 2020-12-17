/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"image/gif"
	"os"

	mandelbrot "github.com/bwarren2/mandelbrot/pkg/mandelbrot"
	"github.com/spf13/cobra"
)

// gifCmd represents the gif command
var gifCmd = &cobra.Command{
	Use:   "gif",
	Short: "Generate a mandelbrot in a file",
	Long: `Mandelbrot generates images of the mandelbrot set.

It outputs in PNG, and is configurable for image size,
range, domain, and iterations`,
	Run: func(cmd *cobra.Command, args []string) {
		colors := mandelbrot.NewPalette(maxIterations)
		output := mandelbrot.Gif(width, height, frames, maxIterations, x, y, scaleIn, colors)
		f, _ := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0600)
		defer f.Close()
		gif.EncodeAll(f, output)
	},
}

// Should these all be in root.go?
var filename string
var width, height, frames uint16
var maxIterations uint8
var x, y, scaleIn float64

func init() {
	gifCmd.Flags().StringVarP(&filename, "filename", "f", "out.gif", "File to write to")
	gifCmd.Flags().Uint16Var(&width, "width", 900, "Width of output image, in pixels")
	gifCmd.Flags().Uint16Var(&height, "height", 450, "Height of output image, in pixels")
	gifCmd.Flags().Uint8Var(&maxIterations, "maxIterations", 100, "Number of iterations to run the mandelbrot loop")
	gifCmd.Flags().Float64Var(&x, "x", -2, "The x point to zoom in on")
	gifCmd.Flags().Float64Var(&y, "y", 0, "The y point to zoom in on")
	gifCmd.Flags().Float64Var(&scaleIn, "scaleIn", .99, "How much to scale the image per frame")
	gifCmd.Flags().Uint16Var(&frames, "frames", 10, "How many frames to draw")
	rootCmd.AddCommand(gifCmd)
}
