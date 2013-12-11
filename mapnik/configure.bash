#!/bin/bash

[ -f mapnik_c_api.cpp ] || curl -o mapnik_c_api.cpp https://raw.github.com/fawick/mapnik-c-api/master/mapnik_c_api.c
[ -f mapnik_c_api.h ] || curl -O https://raw.github.com/fawick/mapnik-c-api/master/mapnik_c_api.h

cat > gen_import.go <<EOF
package mapnik
// #cgo CXXFLAGS: $(mapnik-config --cflags)
// #cgo LDFLAGS: $(mapnik-config --libs)
import "C"

const (
	fontPath = "$(mapnik-config --fonts)"
	pluginPath = "$(mapnik-config --input-plugins)"
)

EOF

