package mapnik
// #cgo CXXFLAGS: -I/usr/include -I/usr/include/mapnik/agg -I/usr/include/mapnik -I/usr/include -I/usr/include/freetype2 -I/usr/include/libpng16 -I/usr/include/libxml2 -I/usr/include/gdal -I/usr/include/postgresql -I/usr/include/cairo -I/usr/include/glib-2.0 -I/usr/lib/x86_64-linux-gnu/glib-2.0/include -I/usr/include/pixman-1 -I/usr/include/uuid -DMAPNIK_MEMORY_MAPPED_FILE -DMAPNIK_HAS_DLCFN -DBIGINT -DBOOST_REGEX_HAS_ICU -DHAVE_JPEG -DMAPNIK_USE_PROJ4 -DHAVE_PNG -DHAVE_WEBP -DHAVE_TIFF -DLINUX -DMAPNIK_THREADSAFE -DBOOST_SPIRIT_NO_PREDEFINED_TERMINALS=1 -DBOOST_PHOENIX_NO_PREDEFINED_TERMINALS=1 -DBOOST_SPIRIT_USE_PHOENIX_V3=1 -DNDEBUG -DHAVE_CAIRO -DGRID_RENDERER -DHAVE_LIBXML2 -std=c++11 -DU_USING_ICU_NAMESPACE=0 -g -O2  -fstack-protector-strong -Wformat -Werror=format-security -DACCEPT_USE_OF_DEPRECATED_PROJ_API_H -Wdate-time -D_FORTIFY_SOURCE=2 -g0 -fvisibility=hidden -fvisibility-inlines-hidden -Wall -pthread -ftemplate-depth-300 -Wsign-compare -Wshadow -O2
// #cgo LDFLAGS: -L/usr/lib -lmapnik -lboost_system
import "C"

const (
	fontPath = "/usr/share/fonts"
	pluginPath = "/usr/lib/mapnik/3.0/input"
)

