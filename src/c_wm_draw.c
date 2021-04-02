#include <cairo/cairo-xlib.h>

void rectangle_shadow(cairo_t* ctx, int x, int y, int w, int h,
                        int shadow_width, int shadow_roughness){

    cairo_set_line_width(ctx, shadow_roughness*0.25);

    for(int sh = 0; sh < shadow_width; sh += shadow_roughness){
        double strength = (1.0 - (double)(sh)/(double)(shadow_width));
        cairo_set_source_rgba(ctx, 0.0, 0.0, 0.0, 0.25 * strength * strength );
        cairo_rectangle(ctx, x-sh, y-sh, w+sh*2, h+sh*2);
        cairo_stroke(ctx);
    }
}

void c_wm_transparent_draw_type_box(cairo_surface_t* surface, int w, int h,
                                    int border_width, int shadow_width, int button_width,
                                    int button_margin_width, int mask_button,
                                    char* title){
    cairo_t* ctx = cairo_create(surface);
    cairo_set_operator(ctx, CAIRO_OPERATOR_CLEAR);
    cairo_paint(ctx);
    cairo_set_operator(ctx, CAIRO_OPERATOR_OVER);

    double shadow_roughness = (double)(shadow_width)/5.0;
    if(shadow_roughness < 1) shadow_roughness = 1;

    cairo_pattern_t* pattern_l[4];
    for(int i=0; i<4; i++){
        pattern_l[i] = cairo_pattern_create_linear(0, 0, w-shadow_width*2, h-shadow_width*2);
        double adj = (double)i * 0.05;
        cairo_pattern_add_color_stop_rgb(pattern_l[i], 0.0, 0.5-adj, 0.4-adj, 0.4-adj);
        cairo_pattern_add_color_stop_rgb(pattern_l[i], 1.0, 0.46-adj, 0.36-adj, 0.36-adj);
        
    }

    cairo_set_source(ctx, pattern_l[0]);

    cairo_rectangle(ctx, shadow_width, shadow_width, w-shadow_width*2, border_width-shadow_width); // top
    cairo_rectangle(ctx, shadow_width, h-border_width, w-shadow_width*2, border_width-shadow_width); // bottom
    cairo_rectangle(ctx, shadow_width, border_width, border_width-shadow_width, h-border_width*2); // left
    cairo_rectangle(ctx, w-border_width, border_width, border_width-shadow_width, h-border_width*2); // right

    cairo_fill(ctx);

    rectangle_shadow(ctx, shadow_width, shadow_width, w-shadow_width*2, h-shadow_width*2,
                        shadow_width, shadow_roughness);

    //title
    
    cairo_select_font_face(ctx, "Serif",
        CAIRO_FONT_SLANT_NORMAL,
        CAIRO_FONT_WEIGHT_NORMAL);

    cairo_set_font_size(ctx, button_width);

    cairo_text_extents_t extents;
    cairo_text_extents(ctx, title, &extents);

    double textbox_x = border_width;
    double textbox_y = border_width - button_width - button_margin_width;
    double textbox_w = extents.width + button_margin_width*2;
    double textbox_h = button_width + button_margin_width;

    rectangle_shadow(ctx, textbox_x,
                          textbox_y,
                          textbox_w,
                          textbox_h,
                          button_margin_width,
                          shadow_roughness);

    cairo_rectangle(ctx, textbox_x, textbox_y, textbox_w, textbox_h);
    cairo_set_source(ctx, pattern_l[0]);
    cairo_fill(ctx);
    
    cairo_move_to(ctx, textbox_x + button_margin_width, textbox_y + textbox_h - button_margin_width);
    cairo_set_source_rgb(ctx, 0.9, 0.9, 0.9);
    cairo_show_text(ctx, title);  
    

    // button

    for(double i=1; i<=3; i++){
        double button_y = border_width - button_width - button_margin_width;

        double button_x = w - border_width - (button_width + button_margin_width)*i;

        rectangle_shadow(ctx, button_x, button_y, button_width, button_width,
                button_margin_width, shadow_roughness);

        cairo_rectangle(ctx, button_x, button_y, button_width, button_width);
        cairo_set_source(ctx, pattern_l[1+(int)i % 2 + (int)(mask_button == i) ]);
        cairo_fill(ctx);

        double icon_margin = button_width*0.3;
        double icon_sx = button_x+icon_margin;
        double icon_sy = button_y+icon_margin;
        double icon_ex = button_x+button_width-icon_margin;
        double icon_ey = button_y+button_width-icon_margin;
        cairo_set_line_width(ctx, button_width*0.1);
        cairo_set_source_rgb(ctx, 0.9, 0.9, 0.9);

        switch((int)i){
            case 1:{
                cairo_move_to(ctx, icon_sx, icon_sy);
                cairo_line_to(ctx, icon_ex, icon_ey);
                cairo_move_to(ctx, icon_sx, icon_ey);
                cairo_line_to(ctx, icon_ex, icon_sy);
                cairo_stroke(ctx);
                break;
            }
            case 2:{
                cairo_move_to(ctx, icon_sx, icon_ey);
                cairo_line_to(ctx, icon_ex, icon_ey);
                cairo_stroke(ctx);
                break;
            }
            case 3:{
                cairo_rectangle(ctx, icon_sx, icon_sy, icon_ex-icon_sx, icon_ey-icon_sy);
                cairo_stroke(ctx);
                break;
            }
        }
    }



}

void c_wm_transparent_draw_type_mask(cairo_surface_t* surface, int w, int h){
    //cairo_t* ctx = cairo_create(surface);
    //cairo_set_operator(ctx, CAIRO_OPERATOR_SOURCE);
    //cairo_set_source_rgba(ctx, 1.0, 1.0, 1.0, 0.4);
    //cairo_paint(ctx);
}