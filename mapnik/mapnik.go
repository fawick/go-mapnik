package mapnik

// #include "mapnik_c_api.h"
import "C"

import (
	"unsafe"
)

func init() {
	// register default datasources path and fonts path like the python bindings do
	RegisterDatasources(pluginPath)
	RegisterFonts(fontPath)
}

func RegisterDatasources(path string) {
	C.mapnik_register_datasources(C.CString(path))
}

func RegisterFonts(path string) {
	C.mapnik_register_fonts(C.CString(path))
}

// Point in 2D space
type Coord struct {
	X, Y float64
}

// Projection from one reference system to the other
type Projection struct {
	p *C.struct_mapnik_projection_t
}

func (p *Projection) Free() {
	C.mapnik_projection_free(p.p)
	p.p = nil
}

func (p Projection) Forward(coord Coord) Coord {
	c := C.mapnik_coord_t{C.double(coord.X), C.double(coord.Y)}
	c = C.mapnik_projection_forward(p.p, c)
	return Coord{float64(c.x), float64(c.y)}
}

// Map base type
type Map struct {
	m *C.struct_mapnik_map_t
}

func NewMap(width, height uint32) *Map {
	return &Map{C.mapnik_map(C.uint(width), C.uint(height))}
}

func (m *Map) Load(stylesheet string) {
	C.mapnik_map_load(m.m, C.CString(stylesheet))
}

func (m *Map) Resize(width, height uint32) {
	C.mapnik_map_resize(m.m, C.uint(width), C.uint(height))
}

func (m *Map) Free() {
	C.mapnik_map_free(m.m)
	m.m = nil
}

func (m *Map) SRS() string {
	return C.GoString(C.mapnik_map_get_srs(m.m))
}

func (m *Map) SetSRS(srs string) {
	C.mapnik_map_set_srs(m.m, C.CString(srs))
}

func (m *Map) ZoomAll() {
	C.mapnik_map_zoom_all(m.m)
}

func (m *Map) ZoomToMinMax(minx, miny, maxx, maxy float64) {
	bbox := C.mapnik_bbox(C.double(minx), C.double(miny), C.double(maxx), C.double(maxy))
	defer C.mapnik_bbox_free(bbox)
	C.mapnik_map_zoom_to_box(m.m, bbox)
}

func (m *Map) RenderToFile(path string) {
	C.mapnik_map_render_to_file(m.m, C.CString(path))
}

func (m *Map) RenderToMemoryPng() []byte {
	i := C.mapnik_map_render_to_image(m.m)
	defer C.mapnik_image_free(i)
	b := C.mapnik_image_to_png_blob(i)
	defer C.mapnik_image_blob_free(b)
	return C.GoBytes(unsafe.Pointer(b.ptr), C.int(b.len))
}
