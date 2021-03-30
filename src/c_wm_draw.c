#include <cairo/cairo-xlib.h>

void c_wm_transparent_draw_type_box(cairo_surface_t* surface, int w, int h){
    cairo_t* ctx = cairo_create(surface);
    cairo_set_source_rgba(ctx, 1.0, 0.5, 0.5, 0.5);
    cairo_paint(ctx);
}