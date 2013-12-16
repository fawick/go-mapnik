package mapnik

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
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

// Generates tile files as a <zoom>/<x>/<y>.png file hierarchie in the current
// work directory.
func (g *Generator) Run(lowLeft, upRight Coord, minZ, maxZ uint64, name string) {
	c := make(chan tileCoord)
	q := make(chan bool)

	fmt.Println("starting job", name)

	ensureDirExists(g.TileDir)

	for i := 0; i < g.Threads; i++ {
		go func(id int, ctc <-chan tileCoord, q chan bool) {
			r := NewTileRenderer(g.MapFile)
			for t := range ctc {
				start := time.Now()
				r.RenderTile(t)
				bytes, _ := r.RenderTile(t)
				ioutil.WriteFile(fmt.Sprintf("%d/%d/%d.png", t.zoom, t.x, t.y), bytes, 0644)
				log.Println(id, t.x, t.y, t.zoom, time.Since(start))
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
			fmt.Println(uint64(px0[1]/256.0), uint64(px1[1]/256.0))
			for y := uint64(px0[1] / 256.0); y <= uint64(px1[1]/256.0); y++ {
				c <- tileCoord{x, y, z, false}
			}
		}
	}
	close(c)
	for i := 0; i < g.Threads; i++ {
		<-q
	}
}
