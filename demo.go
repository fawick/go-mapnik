package main

// This file contains various demo applications of the go-mapnik package

import (
	"fmt"
	"github.com/fawick/go-mapnik/mapnik"
	"github.com/fawick/go-mapnik/maptiles"
	"io/ioutil"
	"net/http"
	"os"
)

const NUM_THREADS = 4

// Render a map of euripe
func SimpleExample() {
	m := mapnik.NewMap(1600, 1200)
	defer m.Free()
	m.Load("sampledata/stylesheet.xml")
	fmt.Println(m.SRS())
	// perform projection, only neccessary because stylesheet.xml is using
	// EPSG:3857 rather than WGS84
	p := m.Projection()
	ll := p.Forward(mapnik.Coord{0, 35})  // 0 degrees longitude, 35 degrees north
	ur := p.Forward(mapnik.Coord{16, 70}) // 16 degrees east, 70 degrees north
	m.ZoomToMinMax(ll.X, ll.Y, ur.X, ur.Y)
	ioutil.WriteFile("mapnik.png", m.RenderToMemoryPng(), 0644)
}

// This function resembles the OSM python script 'generate_tiles.py'
// The original script is found here:
// http://svn.openstreetmap.org/applications/rendering/mapnik/generate_tiles.py
func GenerateOSMTiles() {
	g := maptiles.Generator{}
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

func TileserverWithCaching() {
	cache := "gomapnikcache.sqlite"
	os.Remove(cache)
	t := maptiles.NewTileServer(cache)
	t.AddMapnikLayer("", "sampledata/stylesheet.xml")
	http.ListenAndServe(":8080", t)
}

func main() {
	SimpleExample()
	//GenerateOSMTiles()
	TileserverWithCaching()
}
