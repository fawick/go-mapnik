package mapnik

// #include "mapnik_c_api.h"
import "C"

// TODO func init() for running RegisterDatasources (like in the python bindings)

func RegisterDatasources(path string) {
	C.mapnik_register_datasources(C.CString(path))
}

type Coord struct {
	X, Y float64
}

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

type Map struct {
	m *C.struct_mapnik_map_t
}

func NewMap(width, height uint32) *Map {
	//m := new(Map)
	//m.m = C.mapnik_map(C.uint(width), C.uint(height))
	//return m
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

func (m *Map) RenderToFile(path string) {
	C.mapnik_map_render_to_file(m.m, C.CString(path))
}
