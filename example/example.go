package main

import (
	"github.com/fawick/go-mapnik/mapnik"
)

func main() {
	mapnik.RegisterDatasources("/usr/lib/mapnik/input")
	m := mapnik.NewMap(1600, 1200)
	defer m.Free()
	m.Load("sample/stylesheet.xml")
	m.ZoomAll()
	m.RenderToFile("mapnik.png")
}
