package maptiles

import (
	"log"
	"net/http"
	"regexp"
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

var pathRegex = regexp.MustCompile(`/([A-Za-z0-9]+)/([0-9]+)/([0-9]+)/([0-9]+)\.png`)

func (t *TileServer) ServeTileRequest(w http.ResponseWriter, r *http.Request, tc TileCoord) {
	ch := make(chan TileFetchResult)

	tr := TileFetchRequest{tc, ch}
	t.m.RequestQueue() <- tr

	result := <-ch
	needsInsert := false

	if result.BlobPNG == nil {
		// Tile was not provided by DB, so submit the tile request to the renderer
		t.lmp.SubmitRequest(tr)
		result = <-ch
		if result.BlobPNG == nil {
			// The tile could not be rendered, now we need to bail out.
			http.NotFound(w, r)
			return
		}
		needsInsert = true
	}

	w.Header().Set("Content-Type", "image/png")
	_, err := w.Write(result.BlobPNG)
	if err != nil {
		log.Println(err)
	}
	if needsInsert {
		t.m.InsertQueue() <- result // insert newly rendered tile into cache db
	}
}

func (t *TileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := pathRegex.FindStringSubmatch(r.URL.Path)

	if path == nil {
		http.NotFound(w, r)
		return
	}

	l := path[1]
	z, _ := strconv.ParseUint(path[2], 10, 64)
	x, _ := strconv.ParseUint(path[3], 10, 64)
	y, _ := strconv.ParseUint(path[4], 10, 64)

	t.ServeTileRequest(w, r, TileCoord{x, y, z, t.TmsSchema, l})
}
