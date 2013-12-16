package main

// This file contains various demo applications of the go-mapnik package

import (
	"fmt"
	"github.com/fawick/go-mapnik/mapnik"
	"io/ioutil"
	"os"
)

const NUM_THREADS = 4

func SimpleExample() {
	m := mapnik.NewMap(1600, 1200)
	defer m.Free()
	m.Load("sampledata/stylesheet.xml")
	fmt.Println(m.SRS())
	m.ZoomToMinMax(0, 35, 16, 70)
	ioutil.WriteFile("mapnik.png", m.RenderToMemoryPng(), 0644)
}

// This function resembles the OSM python script 'generate_tiles.py'
// The original script is found here:
// http://svn.openstreetmap.org/applications/rendering/mapnik/generate_tiles.py
func GenerateOSMTiles() {
	g := mapnik.Generator{}
	g.Threads = NUM_THREADS

	home := os.Getenv("HOME")
	g.MapFile = os.Getenv("MAPNIK_MAP_FILE")
	if g.MapFile == "" {
		g.MapFile = home + "/svn.openstreetmap.org/applications/rendering/mapnik/osm-local.xml"
	}
	g.TileDir = os.Getenv("MAPNIK_TILE_DIR")
	if g.TileDir == "" {
		g.TileDir = home + "/osm/tiles"
	}

	g.Run(mapnik.Coord{-180, -90}, mapnik.Coord{180, 90}, 0, 6, "World")
	g.Run(mapnik.Coord{11.4, 48.07}, mapnik.Coord{11.7, 48.22}, 1, 12, "Muenchen")
	g.Run(mapnik.Coord{11.3, 48.01}, mapnik.Coord{12.15, 48.44}, 7, 12, "Muenchen+")
	g.Run(mapnik.Coord{1.0, 10.0}, mapnik.Coord{20.6, 50.0}, 1, 11, "Europe+")
}

func main() {
	//SimpleExample()
	GenerateOSMTiles()
}
