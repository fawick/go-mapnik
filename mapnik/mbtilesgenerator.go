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

type Mbtiles struct {
	db          *sql.DB
	requestChan chan TileFetchRequest
	insertChan  chan TileFetchResult
}

func NewMbtiles(path string) *Mbtiles {
	m := Mbtiles{}
	var err error
	m.db, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Println("Error opening db", err.Error())
		return nil
	}
	queries := []string{
		"PRAGMA journal_mode = OFF",
		"CREATE TABLE IF NOT EXISTS layers(layer_name text PRIMARY KEY NOT NULL);",
		"CREATE TABLE IF NOT EXISTS metadata (name text PRIMARY KEY NOT NULL, value text NOT NULL);",
		"CREATE TABLE IF NOT EXISTS layered_tiles (layer_id integer, zoom_level integer, tile_column integer, tile_row integer, checksum text, PRIMARY KEY (layer_id, zoom_level, tile_column, tile_row) FOREIGN KEY(checksum) REFERENCES tile_blobs(checksum));",
		"CREATE TABLE IF NOT EXISTS tile_blobs (checksum text, tile_data blob);",
		"CREATE VIEW IF NOT EXISTS tiles AS SELECT layered_tiles.zoom_level as zoom_level, layered_tiles.tile_column as tile_column, layered_tiles.tile_row as tile_row, (SELECT tile_data FROM tile_blobs WHERE checksum=layered_tiles.checksum) as tile_data FROM layered_tiles WHERE layered_tiles.layer_id = (SELECT rowid FROM layers WHERE layer_name='default');",
		"REPLACE INTO metadata VALUES('name', 'go-mapnik cache file');",
		"REPLACE INTO metadata VALUES('type', 'overlay');",
		"REPLACE INTO metadata VALUES('version', '0');",
		"REPLACE INTO metadata VALUES('description', 'Compatible with MBTiles spec 1.2. However, this file may contain multiple overlay layers, but only the layer called default is exported as MBtiles');",
		"REPLACE INTO metadata VALUES('format', 'png');",
		"REPLACE INTO metadata VALUES('bounds', '-180.0,-85,180,85');",
		"INSERT OR IGNORE INTO layers(layer_name) VALUES('default');",
	}

	for _, query := range queries {
		_, err = m.db.Exec(query)
		if err != nil {
			log.Println("Error setting up db", err.Error())
			return nil
		}
	}
	m.insertChan = make(chan TileFetchResult)
	m.requestChan = make(chan TileFetchRequest)
	go m.Run()
	return &m
}

func (m *Mbtiles) Close() {
	close(m.insertChan)
	close(m.requestChan)
}

func (m Mbtiles) InsertQueue() chan<- TileFetchResult {
	return m.insertChan
}

func (m Mbtiles) RequestQueue() chan<- TileFetchRequest {
	return m.requestChan
}

// Best executed in a dedicated go routine.
func (m *Mbtiles) Run() {
	for {
		select {
		case r := <-m.requestChan:
			m.fetch(r)
		case i := <-m.insertChan:
			m.insert(i)
		}
	}
}

func (m *Mbtiles) insert(i TileFetchResult) {
	i.Coord.setTMS(true)
	x, y, z := i.Coord.X, i.Coord.Y, i.Coord.Zoom
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
	sql := "REPLACE INTO layered_tiles VALUES( (SELECT rowid FROM layers WHERE layer_name='default'), ?, ?, ?, ?)"
	if _, err = m.db.Exec(sql, z, x, y, s); err != nil {
		log.Println(err)
	}
}

func (m *Mbtiles) fetch(r TileFetchRequest) {
	r.Coord.setTMS(true)
	zoom, x, y := r.Coord.Zoom, r.Coord.X, r.Coord.Y
	result := TileFetchResult{r.Coord, nil}
	queryString := fmt.Sprintf("select tile_data from tiles where zoom_level=%d and tile_column=%d and tile_row=%d",
		zoom, x, y)
	var blob []byte
	row := m.db.QueryRow(queryString)
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
