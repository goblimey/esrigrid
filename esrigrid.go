package main

import (
	"flag"
	"github.com/goblimey/esrigrid/model"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

var filename string   // The point cloud file to display.
var output string     // The .png results file.
var ceiling64 float64 // parameter - the maximum height expected.
var ceiling float32   // ceiling as a float32
var floor64 float64   // parameter - the minimum height expected.
var floor float32     // floor as a float32
var verbose bool      // verbose mode

var maxHeight float64 = 0
var maxHeightSupplied = false	// true if the maxHeight was supplied on the command line.
var minHeight float64 = 0
var minHeightSupplied = false	// true if the minHeight was supplied on the command line.
var NUMBER_OF_SHADES = 256;		// Number of shades of grey available.

func init() {
	flag.StringVar(&filename, "input", "", "point cloud data file")
	flag.StringVar(&filename, "i", "", "point cloud data file")
	flag.StringVar(&output, "output", "", ".png results file")
	flag.StringVar(&output, "o", "", ".png results file")
	flag.Float64Var(&ceiling64, "ceiling", 0.0, "maximum height expected")
	flag.Float64Var(&ceiling64, "c", 0.0, "maximum height expected")
	flag.Float64Var(&floor64, "floor", 0.0, "mimimum height expected")
	flag.Float64Var(&floor64, "f", 0.0, "minimum height expected")
	flag.BoolVar(&verbose, "verbose", false, "verbose mode")
	flag.BoolVar(&verbose, "v", false, "verbose mode")
}

func main() {

	// Get the command line arguments.
	flag.Parse()

	// flagset contains the names of the flags that were supplied on the command line.
	flagset := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { flagset[f.Name] = true })

	// Create the output .png file.
	out, err := os.Create(output)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	// Create an esrigrid object from the given data file.
	pc := model.MakeEsriGrid()
	err = pc.ReadEsriGridFromFile(filename, verbose)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	// If the floor is not already set from the command line, set it from the object.
	if !(flagset["floor"] || flagset["f"]) {
		floor = pc.MinHeight()
	}

	// If the ceiling is not already set from the command line, set it from the object.
	if !(flagset["ceiling"] || flagset["f"]) {
		ceiling = pc.MaxHeight()
	}

	// Create an RGBA image with one pixel for each grid cell.  The origin is at the top left,
	// same as the grid.
	img := image.NewRGBA(image.Rect(0, 0, pc.Nrows(), pc.Ncols()))
	for row := 0; row < pc.Nrows(); row++ {
		for col := 0; col < pc.Ncols(); col++ {
			s := shade(floor, ceiling, pc.Height(row, col))
			if verbose {
				log.Printf("shading cell[%d[%d] %d\n", row, col, s)
			}
			img.Set(col, row, s)
		}
	}

	// Write the RGBA to the PNG image file.
	err = png.Encode(out, img)
}

func shade(floor, ceiling, height float32) color.Color {
	// Get height and ceiling relative to the floor.
	height = height - floor
	ceiling = ceiling - floor
	shade := uint8(NUMBER_OF_SHADES-1) - uint8(height*float32(NUMBER_OF_SHADES)/ceiling)
	if verbose {
		log.Printf("shade %d", shade)
	}
	return color.Gray{shade}
}
