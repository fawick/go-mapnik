#!/bin/bash

[ -f mapnik_c_api.cpp ] || curl -Lo mapnik_c_api.cpp https://raw.github.com/springmeyer/mapnik-c-api/master/mapnik_c_api.c
[ -f mapnik_c_api.h ] || curl -LO https://raw.github.com/springmeyer/mapnik-c-api/master/mapnik_c_api.h

cat > gen_import.go <<EOF
package mapnik
// #cgo CXXFLAGS: $(mapnik-config --cflags)
// #cgo LDFLAGS: $(mapnik-config --libs) -lboost_system
import "C"

const (
	fontPath = "$(mapnik-config --fonts)"
	pluginPath = "$(mapnik-config --input-plugins)"
)

EOF

go install -x
