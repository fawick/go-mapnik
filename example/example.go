package main

import (
	"fmt"
	"github.com/fawick/go-mapnik/mapnik"
	"io/ioutil"
)

func main() {
	m := mapnik.NewMap(1600, 1200)
	defer m.Free()
	m.Load("sample/stylesheet.xml")
	fmt.Println(m.SRS())
	m.ZoomToMinMax(0, 35, 16, 70)
	ioutil.WriteFile("mapnik.png", m.RenderToMemoryPng(), 0644)
}
