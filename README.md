# esrigrid
A Go object that holds an Esri Grid dataset for Geographical Information Systems (GIS)

The EsriGrid interface defines an object that can store an Esri Grid,
a rectangular grid of height or slope values,
often used to represent a geographical map.
The package also offers a function that 
reads a file in Esri Grid ASCII format (ARC/INFO ASCII GRID format).
and stores the contents.
ASCII format simply means that the file is in plain text.
It's essentially a list of height or slope measurements in rows and columns.

An implementation of the interface is also provided, and a demonstration program
showing how it can be used.

The file format is described here:
https://en.wikipedia.org/wiki/Esri_gridPoint.
This is a very small Esri Grid file with six rows and four columns:

    ncols         4
    nrows         6
    xllcorner     0.0
    yllcorner     0.0
    cellsize      50.0
    NODATA_value  -9999
    -9999 -9999 5 2
    -9999 20 100 36
    3 8 35 10
    32 42 50 6
    88 75 27 9
    13 5 1 -9999
   
The first six lines define the rest of the data.  Ncols is the number
of columns, nrows the number of rows.  Xllcorner gives the x map reference of the bottom
left corner of the grid, yllcorner the y map reference.  Cellsize is the size of the cells
in the grid.  The NODATA value is used for points on the grid where the sensor couldn't
figure out the height.

The header lines are followed by the rows and columns of height or slope data.  The values can be
floating point numbers, here they are integers.
The first row defines the top (most northern) line of the grid and the last row describes the
bottom (most Southern) line, so the first number of the last line is the height at
(xllcorner, yllcorner).

Real Esri Grid files are much bigger.

Lidar mapping sensors can measure the height at the tops of the trees and of the ground below.
The first of those measurements is used to create a Digital Surface
Model (DSM)
and the second a
Digital Terrain Model (DTM).

Lidar data is very useful for spotting features in the landscape.
It's been used for many years for applications including forestry, flood management and archaeology, 
but lately the sensors have become smaller and cheaper.
It's now feasible to mount one on a drone and produce your own high-definition maps.
If you are interested in that,
you may find [this discussion group](https://groups.google.com/forum/#!forum/lidar-mapping) useful.


# Downloading the Example Data

The UK Environment Agency publishes Lidar data in Esri Grid ASCII format.
These data are available for free and cover most of the
UK at 1m resolution or better.
They will cover the whole of England by the
end of 2018
and eventually the whole of the UK.
The measurements in the files are heights above sea level in metres.
The corner values are a UK Ordnance Survey numeric map reference.

The data on the Environment Agency website is not presented in a very obvious way.
It's arranged as collections of,
DTMs and DSMs at various resolutions,
grouped into sets of four zip files, each covering
the NE, NW, SE and SW part of a 100 Km square.
Each file in the zip covers a 1 Km square.

For example, the 1 Km square TQ1652 covers an area to the North of Dorking in Surrey.
It's part of the 100 Km square TQ15.
To download the data for TQ1652,
start at this page(https://environment.data.gov.uk/ds/survey/#/survey),
which displays a map of the UK.
Move across the map to the South of London.
As you hover over each square,
its reference number appears.
Find square TQ15 and click on it.
That produces some download buttons
at the bottom of the page.
Each button downloads a zip file containing some Esri Grid data files.
Download the SE part of the 1m resolution DTM set and
extract the files.
You need the file tq1652_DTM_1m.asc.


Installation

To install the software
you first need to install the go compiler.
The official instructions for that are [here](https://golang.org/doc/install).
If you haven't done this sort of thing before,
the instructions [here](http://www.goblimey.com/scaffolder/1.1.installing.go.html) are a bit gentler.

Once you've installed Go, you can install the esrigrid software like so:

    go get github.com/goblimey/esrigrid


# The Demonstration Program

The demonstration program uses the esrigrid module to
read a file and render it as a picture in .png format.
The picture has one pixel per grid cell and paints each cell in one of 256 shades of grey,
white for the lowest values (the floor) and pure black for the highest (the ceiling).

By default the demonstration program figures out the floor and ceiling values
as it reads the file.
You can also set your own floor and ceiling,
wich is useful if you want to compare two grid files.

Run the demo like so:

    esrigrid -i <esrigrid_file> -o <picture>

To see the complete list of options, run it like so:

    esrigrid -h

Process file tq1652_DTM_1m.asc like so:

    esrigrid -i tq1652_DTM_1M.asc -o tq1652.png

You can clearly see the River Mole, the A24 and the railway line in the picture.
The high ground on the right (in darker shades) is part of Box Hill.

This program is a very simple example showing how Lidar data can be visualised.
There are already plenty of tools available off the shelf that process the data
in much more sophisticated ways to make all sorts of structures stand out.


'''
PACKAGE DOCUMENTATION

package model
    import "github.com/goblimey/esrigrid/model"


TYPES

    The EsriGrid interface defines a data structure that holds Esri Grid data 
    representing a surface within some mapping system.  The format is described
    here:  https://en.wikipedia.org/wiki/Esri_gridPoint. Equipment such as Lidar
    mapping sensors produce data in this format and there is a host of
    software available to process and visualise it.

    The UK Environment Agency publish mapping data in this format which can
    be downloaded for free here:
    https://environment.data.gov.uk/ds/survey/#/survey. 

    The interface includes a method that reads a data file in the ASCII
    (plain text) version of the Esri format and sets up an EsriGrid.

type EsriGrid interface {
    NCols returns the number of columns.
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
    // contains filtered or unexported fields
}

func (ceg ConcreteEsriGrid) CellSize() float32
    Cellsize returns the size of the cells.

func (ceg ConcreteEsriGrid) Height(row, col int) float32
    Height returns the height at the intersection of a row and column

func (ceg ConcreteEsriGrid) MaxHeight() float32
    MaxHeight returns the biggest height value (the ceiling).

func (ceg ConcreteEsriGrid) MinHeight() float32
    MinHeight returns the smallest height value (the floor).

func (ceg ConcreteEsriGrid) Ncols() int
    NCols returns the number of columns.

func (ceg ConcreteEsriGrid) NoDataValue() float32
    NoDataValue returns the height value set when there is no data.

func (ceg ConcreteEsriGrid) Nrows() int
    Nrows returns the number of rows.

func (ceg *ConcreteEsriGrid) ReadEsriGridFromFile(filename string, verbose bool) error
    ReadEsrigridFromFile reads the data from a plain text EsriGrid file and
    sets the fields.

func (ceg *ConcreteEsriGrid) SetCellSize(cellsize float32)
    SetCellSize sets the cell size.

func (ceg *ConcreteEsriGrid) SetHeight(row, col int, height float32)
    SetHeight sets the height at the intersection of a row and column.

func (ceg *ConcreteEsriGrid) SetNCols(ncols int)
    SetNCols sets the number of columns.

func (ceg *ConcreteEsriGrid) SetNRows(nrows int)
    SetNRows sets the number of rows.

func (ceg *ConcreteEsriGrid) SetNoDataValue(noDataValue float32)
    SetNoDataValue sets the no data value

func (ceg *ConcreteEsriGrid) SetXllcorner(xllcorner float32)
    SetXllcorner sets the xllcorner value.

func (ceg *ConcreteEsriGrid) SetYllcorner(yllcorner float32)
    SetYllcorner sets the yllcorner value.

func (ceg ConcreteEsriGrid) Xllcorner() float32
    Xllcorner returns the xllcorner (the east component of the map reference
    of the bottom left cell).

func (ceg ConcreteEsriGrid) Yllcorner() float32
    Yllcorner returns the yllcorner (the north component of the map
    reference of the bottom left cell).

func MakeEsriGrid() EsriGrid
    MakeEsriGrid creates and returns an EsriGrid object.
'''
