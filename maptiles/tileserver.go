package maptiles

import (
	"log"
	"net/http"
	"strconv"
)

// TODO serve list of registered layers per HTTP (preferably leafletjs-compatible js-array)

// Handles HTTP requests for map tiles, caching any produced tiles
// in an MBtiles 1.2 compatible sqlite db.
type TileServer struct {
	m         *TileDb
	lmp       *LayerMultiplex
	TmsSchema bool
}

func NewTileServer(cacheFile string) *TileServer {
	t := TileServer{}
	t.lmp = NewLayerMultiplex()
	t.m = NewTileDb(cacheFile)

	return &t
}

func (t *TileServer) AddMapnikLayer(layerName string, stylesheet string) {
	t.lmp.AddRenderer(layerName, stylesheet)
}

func (t *TileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}

	// TODO use uri scheme <z>/<x>/<y>.png instead of form values
	x, _ := strconv.ParseUint(r.FormValue("x"), 10, 64)
	y, _ := strconv.ParseUint(r.FormValue("y"), 10, 64)
	z, _ := strconv.ParseUint(r.FormValue("z"), 10, 64)
	l := r.FormValue("l")
	ch := make(chan TileFetchResult)

	request := TileFetchRequest{TileCoord{x, y, z, t.TmsSchema, l}, ch}
	t.m.RequestQueue() <- request

	result := <-ch
	needsInsert := false

	if result.BlobPNG == nil {
		// Tile was not provided by DB, so submit the tile request to the renderer
		log.Println("tile cache miss", z, x, y)
		t.lmp.SubmitRequest(request)
		result = <-ch
		if result.BlobPNG == nil {
			// The tile could not be rendered, now we need to bail out.
			http.NotFound(w, r)
			return
		}
		needsInsert = true
	} else {
		log.Println("tile cache hit", z, x, y)
	}

	w.Header().Set("Content-Type", "image/png")
	_, err = w.Write(result.BlobPNG)
	if err != nil {
		log.Println(err)
	}
	if needsInsert {
		t.m.InsertQueue() <- result // insert newly rendered tile into cache db
	}
}
