#include <cairo/cairo-xlib.h>

void c_wm_transparent_draw_type_box(cairo_surface_t* surface, int w, int h, int border_width){
    cairo_t* ctx = cairo_create(surface);
    cairo_set_operator(ctx, CAIRO_OPERATOR_SOURCE);

    int shadow_width = border_width/5;
    int shadow_roughness = shadow_width/3;
    if(shadow_roughness < 1) shadow_roughness = 1;

    cairo_pattern_t *pattern;
    pattern = cairo_pattern_create_linear(0, 0, w-shadow_width*2, h-shadow_width*2);

    cairo_pattern_add_color_stop_rgb(pattern, 0.0, 0.5, 0.4, 0.4);
    cairo_pattern_add_color_stop_rgb(pattern, 1.0, 0.46, 0.36, 0.36);
    cairo_set_source(ctx, pattern);

    cairo_rectangle(ctx, shadow_width, shadow_width, w-shadow_width*2, h-shadow_width*2);
    cairo_fill(ctx);

    cairo_set_line_width(ctx, shadow_roughness);

    for(int sh = 0; sh < shadow_width; sh += shadow_roughness){
        cairo_set_source_rgba(ctx, 0.0, 0.0, 0.0, 0.1 * (1.0 - (double)sh/(double)shadow_width) );
        cairo_rectangle(ctx, shadow_width-sh, shadow_width-sh, w-shadow_width*2+sh*2, h-shadow_width*2+sh*2);
        cairo_stroke(ctx);
    }
}

void c_wm_transparent_draw_type_mask(cairo_surface_t* surface, int w, int h){
    cairo_t* ctx = cairo_create(surface);
    cairo_set_operator(ctx, CAIRO_OPERATOR_SOURCE);
    cairo_set_source_rgba(ctx, 1.0, 1.0, 1.0, 0.4);
    cairo_paint(ctx);
}