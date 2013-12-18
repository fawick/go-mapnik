package mapnik

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	//"net/http"
	//"strconv"
)

// MBTiles 1.2-compatible Tile Db with multi-layer support.
// Was named Mbtiles before, hence the use of *m in methods.
type TileDb struct {
	db          *sql.DB
	requestChan chan TileFetchRequest
	insertChan  chan TileFetchResult
	layerIds    map[string]int
	qc          chan bool
}

func NewTileDb(path string) *TileDb {
	m := TileDb{}
	var err error
	m.db, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Println("Error opening db", err.Error())
		return nil
	}
	queries := []string{
		"PRAGMA journal_mode = OFF",
		"CREATE TABLE IF NOT EXISTS layers(layer_name text PRIMARY KEY NOT NULL)",
		"CREATE TABLE IF NOT EXISTS metadata (name text PRIMARY KEY NOT NULL, value text NOT NULL)",
		"CREATE TABLE IF NOT EXISTS layered_tiles (layer_id integer, zoom_level integer, tile_column integer, tile_row integer, checksum text, PRIMARY KEY (layer_id, zoom_level, tile_column, tile_row) FOREIGN KEY(checksum) REFERENCES tile_blobs(checksum))",
		"CREATE TABLE IF NOT EXISTS tile_blobs (checksum text, tile_data blob)",
		"CREATE VIEW IF NOT EXISTS tiles AS SELECT layered_tiles.zoom_level as zoom_level, layered_tiles.tile_column as tile_column, layered_tiles.tile_row as tile_row, (SELECT tile_data FROM tile_blobs WHERE checksum=layered_tiles.checksum) as tile_data FROM layered_tiles WHERE layered_tiles.layer_id = (SELECT rowid FROM layers WHERE layer_name='default')",
		"REPLACE INTO metadata VALUES('name', 'go-mapnik cache file')",
		"REPLACE INTO metadata VALUES('type', 'overlay')",
		"REPLACE INTO metadata VALUES('version', '0')",
		"REPLACE INTO metadata VALUES('description', 'Compatible with MBTiles spec 1.2. However, this file may contain multiple overlay layers, but only the layer called default is exported as MBtiles')",
		"REPLACE INTO metadata VALUES('format', 'png')",
		"REPLACE INTO metadata VALUES('bounds', '-180.0,-85,180,85')",
		"INSERT OR IGNORE INTO layers(layer_name) VALUES('default')",
	}

	for _, query := range queries {
		_, err = m.db.Exec(query)
		if err != nil {
			log.Println("Error setting up db", err.Error())
			return nil
		}
	}

	m.readLayers()

	m.insertChan = make(chan TileFetchResult)
	m.requestChan = make(chan TileFetchRequest)
	go m.Run()
	return &m
}

func (m *TileDb) readLayers() {
	m.layerIds = make(map[string]int)
	rows, err := m.db.Query("SELECT rowid, layer_name FROM layers")
	if err != nil {
		log.Fatal("Error fetching layer definitions", err.Error())
	}
	var s string
	var i int
	for rows.Next() {
		if err := rows.Scan(&i, &s); err != nil {
			log.Fatal(err)
		}
		m.layerIds[s] = i
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func (m *TileDb) ensureLayer(layer string) {
	if _, ok := m.layerIds[layer]; !ok {
		if _, err := m.db.Exec("INSERT OR IGNORE INTO layers(layer_name) VALUES(?)", layer); err != nil {
			log.Println(err)
		}
		m.readLayers()
	}
}

func (m *TileDb) Close() {
	close(m.insertChan)
	close(m.requestChan)
	if m.qc != nil {
		<-m.qc // block until channel qc is closed (meaning Run() is finished)
	}
	if err := m.db.Close(); err != nil {
		log.Print(err)
	}

}

func (m TileDb) InsertQueue() chan<- TileFetchResult {
	return m.insertChan
}

func (m TileDb) RequestQueue() chan<- TileFetchRequest {
	return m.requestChan
}

// Best executed in a dedicated go routine.
func (m *TileDb) Run() {
	m.qc = make(chan bool)
	for {
		select {
		case r := <-m.requestChan:
			m.fetch(r)
		case i := <-m.insertChan:
			m.insert(i)
		}
	}
	m.qc <- true
}

func (m *TileDb) insert(i TileFetchResult) {
	i.Coord.setTMS(true)
	x, y, z, l := i.Coord.X, i.Coord.Y, i.Coord.Zoom, i.Coord.Layer
	if l == "" {
		l = "default"
	}
	h := md5.New()
	_, err := h.Write(i.BlobPNG)
	if err != nil {
		log.Println(err)
		return
	}
	s := fmt.Sprintf("%x", h.Sum(nil))
	row := m.db.QueryRow("SELECT 1 FROM tile_blobs WHERE checksum=?", s)
	var dummy uint64
	err = row.Scan(&dummy)
	switch {
	case err == sql.ErrNoRows:
		if _, err = m.db.Exec("REPLACE INTO tile_blobs VALUES(?,?)", s, i.BlobPNG); err != nil {
			log.Println("error during insert", err)
			return
		}
	case err != nil:
		log.Println("error during test", err)
		return
	default:
		//log.Println("Reusing blob", s)
	}
	m.ensureLayer(l)
	sql := "REPLACE INTO layered_tiles VALUES(?, ?, ?, ?, ?)"
	if _, err = m.db.Exec(sql, m.layerIds[l], z, x, y, s); err != nil {
		log.Println(err)
	}
}

func (m *TileDb) fetch(r TileFetchRequest) {
	r.Coord.setTMS(true)
	zoom, x, y, l := r.Coord.Zoom, r.Coord.X, r.Coord.Y, r.Coord.Layer
	if l == "" {
		l = "default"
	}
	result := TileFetchResult{r.Coord, nil}
	queryString := `
		SELECT tile_data 
		FROM tile_blobs 
		WHERE checksum=(
			SELECT checksum 
			FROM layered_tiles 
			WHERE zoom_level=? 
				AND tile_column=? 
				AND tile_row=?
				AND layer_id=(SELECT rowid FROM layers WHERE layer_name=?)
		)`
	var blob []byte
	row := m.db.QueryRow(queryString, zoom, x, y, l)
	err := row.Scan(&blob)
	switch {
	case err == sql.ErrNoRows:
		result.BlobPNG = nil
	case err != nil:
		log.Println(err)
	default:
		result.BlobPNG = blob
	}
	r.OutChan <- result
}
