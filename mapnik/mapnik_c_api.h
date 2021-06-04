#ifndef MAPNIK_C_API_H
#define MAPNIK_C_API_H

#if defined(WIN32) || defined(WINDOWS) || defined(_WIN32) || defined(_WINDOWS)
#  define MAPNIKCAPICALL __declspec(dllexport)
#else
#  define MAPNIKCAPICALL __attribute__ ((visibility ("default")))
#endif

#ifdef __cplusplus
extern "C"
{
#endif

MAPNIKCAPICALL int mapnik_register_datasources(const char* path, char** err);
MAPNIKCAPICALL int mapnik_register_fonts(const char* path, char** err);
MAPNIKCAPICALL const char * mapnik_version_string();


// Coord
typedef struct _mapnik_coord_t {
    double x;
    double y;
} mapnik_coord_t;

// Projection
typedef struct _mapnik_projection_t mapnik_projection_t;

MAPNIKCAPICALL void mapnik_projection_free(mapnik_projection_t *p);

MAPNIKCAPICALL mapnik_coord_t mapnik_projection_forward(mapnik_projection_t *p, mapnik_coord_t c);


// Bbox
typedef struct _mapnik_bbox_t mapnik_bbox_t;

MAPNIKCAPICALL mapnik_bbox_t * mapnik_bbox(double minx, double miny, double maxx, double maxy);

MAPNIKCAPICALL void mapnik_bbox_free(mapnik_bbox_t * b);


// Image
typedef struct _mapnik_image_t mapnik_image_t;

MAPNIKCAPICALL void mapnik_image_free(mapnik_image_t * i);

typedef struct _mapnik_image_blob_t {
    char *ptr;
    unsigned int len;
} mapnik_image_blob_t;

MAPNIKCAPICALL void mapnik_image_blob_free(mapnik_image_blob_t * b);

MAPNIKCAPICALL mapnik_image_blob_t * mapnik_image_to_png_blob(mapnik_image_t * i);



//  Map
typedef struct _mapnik_map_t mapnik_map_t;

MAPNIKCAPICALL mapnik_map_t * mapnik_map( unsigned int width, unsigned int height );

MAPNIKCAPICALL void mapnik_map_free(mapnik_map_t * m);

MAPNIKCAPICALL const char * mapnik_map_last_error(mapnik_map_t * m);

MAPNIKCAPICALL const char * mapnik_map_get_srs(mapnik_map_t * m);

MAPNIKCAPICALL int mapnik_map_set_srs(mapnik_map_t * m, const char* srs);

MAPNIKCAPICALL int mapnik_map_load(mapnik_map_t * m, const char* stylesheet);

MAPNIKCAPICALL int mapnik_map_load_string(mapnik_map_t * m, const char* stylesheet_string);

MAPNIKCAPICALL int mapnik_map_zoom_all(mapnik_map_t * m);

MAPNIKCAPICALL int mapnik_map_render_to_file(mapnik_map_t * m, const char* filepath);

MAPNIKCAPICALL void mapnik_map_resize(mapnik_map_t * m, unsigned int width, unsigned int height);

MAPNIKCAPICALL void mapnik_map_set_buffer_size(mapnik_map_t * m, int buffer_size);

MAPNIKCAPICALL void mapnik_map_zoom_to_box(mapnik_map_t * m, mapnik_bbox_t * b);

MAPNIKCAPICALL mapnik_projection_t * mapnik_map_projection(mapnik_map_t *m);

MAPNIKCAPICALL mapnik_image_t * mapnik_map_render_to_image(mapnik_map_t * m);

#ifdef __cplusplus
}
#endif


#endif // MAPNIK_C_API_H

