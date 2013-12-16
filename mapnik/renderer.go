package mapnik

import (
	"log"
	"sync"
)

type tileCoord struct {
	x, y, zoom uint64
	tms        bool
}

func (c *tileCoord) setTMS(tms bool) {
	if c.tms != tms {
		c.y = (1 << c.zoom) - c.y - 1
		c.tms = tms
	}
}

type TileRenderer struct {
	m  *Map
	mp Projection
	mu *sync.Mutex
}

func NewTileRenderer(stylesheet string) *TileRenderer {
	t := new(TileRenderer)
	var err error
	if err != nil {
		log.Fatal(err)
	}
	t.m = NewMap(256, 256)
	t.m.Load(stylesheet)
	t.mp = t.m.Projection()
	t.mu = new(sync.Mutex)

	return t
}

func (t *TileRenderer) RenderTile(c tileCoord) ([]byte, error) {
	c.setTMS(false)
	return t.RenderTileZXY(c.zoom, c.x, c.y)
}

// Render a tile with coordinates in Google tile format.
// Most upper left tile is always 0,0. Method is not thread-safe,
// so wrap in a channel communication or with a mutex when accessing
// the same renderer by multiple threads.
func (t *TileRenderer) RenderTileZXY(zoom, x, y uint64) ([]byte, error) {
	// Calculate pixel positions of bottom left & top right
	p0 := [2]float64{float64(x) * 256, (float64(y) + 1) * 256}
	p1 := [2]float64{(float64(x) + 1) * 256, float64(y) * 256}

	// Convert to LatLong(EPSG:4326)
	l0 := fromPixelToLL(p0, zoom)
	l1 := fromPixelToLL(p1, zoom)

	// Convert to map projection (e.g. mercartor co-ords EPSG:3857)
	c0 := t.mp.Forward(Coord{l0[0], l0[1]})
	c1 := t.mp.Forward(Coord{l1[0], l1[1]})

	// Bounding box for the Tile
	t.m.Resize(256, 256)
	t.m.ZoomToMinMax(c0.X, c0.Y, c1.X, c1.Y)
	t.m.SetBufferSize(128)

	blob := t.m.RenderToMemoryPng()
	return blob, nil
}
