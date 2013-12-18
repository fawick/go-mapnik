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

// Render a simple map of europe to a PNG file
func SimpleExample() {
	m := mapnik.NewMap(1600, 1200)
	defer m.Free()
	m.Load("sampledata/stylesheet.xml")
	fmt.Println(m.SRS())
	// Perform a projection that is only neccessary because stylesheet.xml
	// is using EPSG:3857 rather than WGS84
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

	// Modify this number according to your machine!
	g.Threads = 4

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
	g.Run(mapnik.Coord{0, 35.0}, mapnik.Coord{16, 70}, 1, 11, "Europe")
}

// Serve a single stylesheet via HTTP. Open view_tileserver.html in your browser
// to see the results.
// The created tiles are cached in an sqlite database (MBTiles 1.2 conform) so
// successive access a tile is much faster.
func TileserverWithCaching() {
	cache := "gomapnikcache.sqlite"
	os.Remove(cache)
	t := maptiles.NewTileServer(cache)
	t.AddMapnikLayer("", "sampledata/stylesheet.xml")
	http.ListenAndServe(":8080", t)
}

// Before uncommenting the GenerateOSMTiles call make sure you have
// the neccessary OSM sources. Consult OSM wiki for details.
func main() {
	SimpleExample()
	//GenerateOSMTiles()
	TileserverWithCaching()
}
