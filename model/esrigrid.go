package model

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// The EsriGrid interface defines a data structure that holds Esri Grid data representing a 
// surface within some mapping system.  The format is described here:
// https://en.wikipedia.org/wiki/Esri_gridPoint.  Equipment such as Lidar mapping sensors produce
// data in this format and there is a host of software available to process and visualise it.
//
// The UK Environment Agency publish mapping data in this format which can be downloaded for
// free here: https://environment.data.gov.uk/ds/survey/#/survey.
//
// The interface includes a method that reads a data file in the ASCII (plain text) version of
// the Esri format and sets up an EsriGrid.
type EsriGrid interface {
	// NCols returns the number of columns.
	Ncols() int
	// Nrows returns the number of rows.
	Nrows() int
	// Xllcorner returns the xllcorner (the east component of the map reference of the bottom left cell).
	Xllcorner() float32
	// Yllcorner returns the yllcorner (the north component of the map reference of the bottom left cell).
	Yllcorner() float32
	// Cellsize returns the size of the cells.
	CellSize() float32
	// NoDataValue returns the height value set when there is no data.
	NoDataValue() float32
	// MaxHeight returns the biggest height value (the ceiling). 
	MaxHeight() float32
	// MinHeight returns the smallest height value (the floor).
	MinHeight() float32
	// Height returns the height at the intersection of a row and column.
	Height(row, col int) float32
	// SetNCols sets the number of columns. 
	SetNCols(ncols int)
	// SetNRows sets the number of rows.
	SetNRows(nrows int)
	// SetXllcorner sets the xllcorner value.
	SetXllcorner(xllcorner float32)
	// SetYllcorner sets the yllcorner value.
	SetYllcorner(yllcorner float32)
	// SetCellSize sets the cell size.
	SetCellSize(cellsize float32)
	// SetNoDataValue sets the no data value.
	SetNoDataValue(noDataValue float32)
	// SetHeight sets the height at the intersection of a row and column.
	SetHeight(row, col int, height float32)
	// ReadEsrigridFromFile reads the data from a plain text EsriGrid file and sets the fields.
	ReadEsriGridFromFile(filename string, verbose bool) error
}

type ConcreteEsriGrid struct {
	ncols        int
	nrows        int
	xllcorner    float32
	yllcorner    float32
	cellsize     float32
	noDataValue  float32
	maxHeight    float32
	minHeight    float32
	height       [][]float32	// The grid of height data
	minHeightSet bool			// False until minHeight is set
	maxHeightSet bool			// False until maxHeight is set.
	verbose      bool			// Verbose logging mode.
}

// MakeEsriGrid creates and returns an EsriGrid object.
func MakeEsriGrid() EsriGrid {
	return &ConcreteEsriGrid{}
}

// NCols returns the number of columns. 
func (ceg ConcreteEsriGrid) Ncols() int {
	return ceg.ncols
}
// Nrows returns the number of rows.
func (ceg ConcreteEsriGrid) Nrows() int {
	return ceg.nrows
}

// Xllcorner returns the xllcorner (the east component of the map reference of the bottom left cell).
func (ceg ConcreteEsriGrid) Xllcorner() float32 {
	return ceg.xllcorner
}

// Yllcorner returns the yllcorner (the north component of the map reference of the bottom left cell).
func (ceg ConcreteEsriGrid) Yllcorner() float32 {
	return ceg.yllcorner
}

// Cellsize returns the size of the cells.
func (ceg ConcreteEsriGrid) CellSize() float32 {
	return ceg.cellsize
}

// NoDataValue returns the height value set when there is no data.
func (ceg ConcreteEsriGrid) NoDataValue() float32 {
	return ceg.noDataValue
}

// MaxHeight returns the biggest height value (the ceiling).
func (ceg ConcreteEsriGrid) MaxHeight() float32 {
	return ceg.maxHeight
}

// MinHeight returns the smallest height value (the floor).
func (ceg ConcreteEsriGrid) MinHeight() float32 {
	return ceg.minHeight
}

// Height returns the height at the intersection of a row and column 
func (ceg ConcreteEsriGrid) Height(row, col int) float32 {
	return ceg.height[row][col]
}

// SetNCols sets the number of columns.
func (ceg *ConcreteEsriGrid) SetNCols(ncols int) {
	ceg.ncols = ncols
}

// SetNRows sets the number of rows.
func (ceg *ConcreteEsriGrid) SetNRows(nrows int) {
	ceg.nrows = nrows
}

// SetXllcorner sets the xllcorner value.
func (ceg *ConcreteEsriGrid) SetXllcorner(xllcorner float32) {
	ceg.xllcorner = xllcorner
}

// SetYllcorner sets the yllcorner value.
func (ceg *ConcreteEsriGrid) SetYllcorner(yllcorner float32) {
	ceg.yllcorner = yllcorner
}

// SetCellSize sets the cell size.
func (ceg *ConcreteEsriGrid) SetCellSize(cellsize float32) {
	ceg.cellsize = cellsize
}

// SetNoDataValue sets the no data value
func (ceg *ConcreteEsriGrid) SetNoDataValue(noDataValue float32) {
	ceg.noDataValue = noDataValue
}

// SetHeight sets the height at the intersection of a row and column.
func (ceg *ConcreteEsriGrid) SetHeight(row, col int, height float32) {

	if row >= ceg.nrows || col >= ceg.ncols {
		log.Printf("SetHeight(%d,%d) - row or column out of range", row, col)
		return
	}
	ceg.height[row][col] = height

	if height = noDataValue {
		return
	}
	if ceg.maxHeightSet {
		if height > ceg.maxHeight {
			ceg.maxHeight = height
		}
	} else {
		// This is the first value.
		ceg.maxHeight = height
		ceg.maxHeightSet = true
	}

	if ceg.minHeightSet {
		if height < ceg.minHeight {
			ceg.minHeight = height
		}
	} else {
		// This is the first value.
		ceg.minHeight = height
		ceg.minHeightSet = true
	}
}

// ReadEsrigridFromFile reads the data from a plain text EsriGrid file and sets the fields.
func (ceg *ConcreteEsriGrid) ReadEsriGridFromFile(filename string, verbose bool) error {
	m := "ReadEsriGridFromFile"
	if verbose {
		log.Printf("%s: %s", m, filename)
	}

	// This is a very simple example of the input file:
	//
	// ncols         4
	// nrows         6
	// xllcorner     0.0
	// yllcorner     0.0
	// cellsize      50.0
	// NODATA_value  -9999
	// -9999 -9999 5 2
	// -9999 20 100 36
	// 3 8 35 10
	// 32 42 50 6
	// 88 75 27 9
	// 13 5 1 -9999
	//
	// The file starts with six header lines defining the rest of the data.  ncols is the number
	// of columns, nrows the number of rows.  xllcorner gives the x map reference of the bottom
	// left corner of the grid, yllcorner the y map reference of the same point.  cellsize is the 
	// size of the grid cells.  The NODATA value is used for points on the grid where the sensor 
	// couldn't figure out the height.
	//
	// The header lines are followed by the rows and columns of height data.  The values can be
	// floating point numbers, here they are integers.  This example defines a four by four grid.
	// The first row defines the top (most northern) line of the grid and the last row describes
	// the bottom (most Southern) line, so the first number of the last line is the height at
	// (xllcorner, yllcorner).

	in, err := os.Open(filename)
	if err != nil {
		log.Printf(filename + err.Error())
		return err
	}

	r := bufio.NewReader(in)

	lineNum := 0
	fieldName := "ncols"
	ceg.ncols, err = readIntFromHeader(r, fieldName, verbose)
	if err != nil {
		return err
	}
	lineNum++
	if verbose {
		log.Printf("%s: %s %d", m, fieldName, ceg.ncols)
	}

	fieldName = "nrows"
	ceg.nrows, err = readIntFromHeader(r, fieldName, verbose)
	if err != nil {
		return err
	}
	lineNum++
	if verbose {
		log.Printf("%s: %s %d", m, fieldName, ceg.nrows)
	}

	ceg.height = make([][]float32, ceg.nrows)

	for i := 0; i < ceg.nrows; i++ {
		ceg.height[i] = make([]float32, ceg.ncols)
	}

	fieldName = "xllcorner"
	ceg.xllcorner, err = readFloat32FromHeader(r, fieldName, verbose)
	if err != nil {
		return err
	}
	lineNum++
	if verbose {
		log.Printf("%s: %s %f", m, fieldName, ceg.xllcorner)
	}

	fieldName = "yllcorner"
	ceg.yllcorner, err = readFloat32FromHeader(r, fieldName, verbose)
	if err != nil {
		return err
	}
	lineNum++
	if verbose {
		log.Printf("%s: %s %f", m, fieldName, ceg.yllcorner)
	}

	fieldName = "cellsize"
	ceg.cellsize, err = readFloat32FromHeader(r, fieldName, verbose)
	if err != nil {
		return err
	}
	lineNum++
	if verbose {
		log.Printf("%s: %s %f", m, fieldName, ceg.cellsize)
	}

	fieldName = "NODATA_value"
	ceg.noDataValue, err = readFloat32FromHeader(r, fieldName, verbose)
	if err != nil {
		return err
	}
	lineNum++

	if verbose {
		log.Printf("NODATA_value %f", ceg.noDataValue)
	}

	// Read nrows of lines each containing ncols floats, space separated.
	if verbose {
		log.Printf("%s: reading %d data lines", m, ceg.nrows)
	}

	linesExpected := ceg.nrows + 6

	for row := 0; ; row++ {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		lineNum++
		if lineNum > linesExpected {
			log.Printf("%s: warning: file %s has too many lines - expected %d\n", m, filename, linesExpected)
			break
		}
		line, err = stripSpaces(line)
		if err != nil {
			log.Printf("%s: stripSpaces failed - %s", m, err.Error())
			return err
		}
		if verbose {
			log.Println(line)
		}

		numbers := strings.Split(line, " ")
		if len(numbers) > ceg.ncols {
			log.Printf("warning: line %d has too many columns - got %d expected %d\n",
				lineNum, len(numbers), ceg.ncols)
			continue
		}
		if len(numbers) < ceg.ncols {
			log.Printf("warning: line %d has too few columns - got %d expected %d\n",
				lineNum, len(numbers), ceg.ncols)
			continue
		}
		for col := range numbers {
			var f float32
			_, err := fmt.Sscanf(numbers[col], "%f", &f)
			if err != nil {
				log.Printf("%d %d %s", row, col, err.Error())
				return err
			}

			// Set height, maxheight and minHeight
			ceg.SetHeight(row, col, f)

			if verbose {
				log.Printf("height[%d][%d] %f", row, col, ceg.height[row][col])
			}
		}
	}

	if lineNum < linesExpected {
		log.Printf("warning: file %s has too few lines - got %d expected %d\n",
			filename, lineNum, linesExpected)
	}

	fmt.Printf("floor %f ceiling %f", ceg.maxHeight, ceg.minHeight)

	return nil
}

func readIntFromHeader(r *bufio.Reader, fieldName string, verbose bool) (int, error) {
	m := "readIntHeader"
	line, err := r.ReadString('\n')
	if err != nil {
		return 0, err
	}
	if verbose {
		log.Printf("%s: line %s", m, line)
	}
	line, err = stripSpaces(line)
	field := strings.Split(line, " ")
	if field[0] != fieldName {
		log.Printf("%s: expected %s, got %s", m, fieldName, line)
	}
	var result int
	_, err = fmt.Sscanf(field[1], "%d", &result)
	if err != nil {
		return 0, err
	}
	if verbose {
		log.Printf("%s: %s %d", m, fieldName, result)
	}

	return result, nil
}

func readFloat32FromHeader(r *bufio.Reader, fieldName string, verbose bool) (float32, error) {
	m := "readFloat32FromHeader"
	line, err := r.ReadString('\n')
	if err != nil {
		return 0, err
	}
	if verbose {
		log.Printf("%s: line %s", m, line)
	}
	line, err = stripSpaces(line)
	field := strings.Split(line, " ")
	if field[0] != fieldName {
		log.Printf("%s: expected %s, got %s", m, fieldName, line)
	}
	var result float32
	_, err = fmt.Sscanf(field[1], "%f", &result)
	if err != nil {
		return 0, err
	}
	if verbose {
		log.Printf("%s: %s %f", m, fieldName, result)
	}

	return result, nil
}

// stripSpaces removes all white space from the start and end of a string and replaces all
// white space within the string with a single space.
func stripSpaces(s string) (string, error) {
	s = strings.TrimSpace(s)
	re, err := regexp.Compile("  +")
	if err != nil {
		return s, err
	}

	return re.ReplaceAllLiteralString(s, " "), nil
}
