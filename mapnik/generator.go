package mapnik

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type GeneratorJob struct {
	LowLeft, UpRight Coord
	MinZoom, MaxZoom uint64
	Name             string
}

type Generator struct {
	MapFile string
	TileDir string
	Threads int
}

func ensureDirExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}
}

// Generates tile files as a <zoom>/<x>/<y>.png file hierarchy in the current
// work directory.
func (g *Generator) Run(lowLeft, upRight Coord, minZ, maxZ uint64, name string) {
	c := make(chan TileCoord)
	q := make(chan bool)

	log.Println("starting job", name)

	ensureDirExists(g.TileDir)

	for i := 0; i < g.Threads; i++ {
		go func(id int, ctc <-chan TileCoord, q chan bool) {
			requests := NewTileRendererChan(g.MapFile)
			results := make(chan TileFetchResult)
			for t := range ctc {
				requests <- TileFetchRequest{t, results}
				r := <-results
				ioutil.WriteFile(r.Coord.OSMFilename(), r.BlobPNG, 0644)
			}
			q <- true
		}(i, c, q)
	}

	ll0 := [2]float64{lowLeft.X, upRight.Y}
	ll1 := [2]float64{upRight.X, lowLeft.Y}

	for z := minZ; z <= maxZ; z++ {
		px0 := fromLLtoPixel(ll0, z)
		px1 := fromLLtoPixel(ll1, z)

		ensureDirExists(fmt.Sprintf("%d", z))
		for x := uint64(px0[0] / 256.0); x <= uint64(px1[0]/256.0); x++ {
			ensureDirExists(fmt.Sprintf("%d/%d", z, x))
			for y := uint64(px0[1] / 256.0); y <= uint64(px1[1]/256.0); y++ {
				c <- TileCoord{x, y, z, false, ""}
			}
		}
	}
	close(c)
	for i := 0; i < g.Threads; i++ {
		<-q
	}
}
