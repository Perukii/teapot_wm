#include <cairo/cairo-xlib.h>

void c_wm_transparent_draw_type_box(cairo_surface_t* surface, int w, int h){
    cairo_t* ctx = cairo_create(surface);
    cairo_set_source_rgba(ctx, 1.0, 0.5, 0.5, 0.5);
    cairo_paint(ctx);
}

void c_wm_transparent_draw_type_mask(cairo_surface_t* surface, int w, int h){
    cairo_t* ctx = cairo_create(surface);
    
    cairo_set_source_rgba(ctx, 1.0, 1.0, 1.0, 0.4);
    cairo_set_operator(ctx, CAIRO_OPERATOR_SOURCE);
    cairo_paint(ctx);
}