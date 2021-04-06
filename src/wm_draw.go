package main

/*
#include "./c_wm_x_access.h"
*/
import "C"

import "github.com/ungerik/go-cairo"
import "unsafe"

func (host *WmHost) wm_host_draw_client(address WmClientAddress){
	/*
	clt := host.client[address]


	var maximized_val int = 0
	if clt.maximize_mode == WM_CLIENT_MAXIMIZE_MODE_REVERSE { maximized_val = 1 }

	C.c_wm_transparent_draw_box(clt.mask.surface, surface_w, surface_h,
				C.int(host.setting.client_border_overall_width),
				C.int(host.setting.client_border_shadow_width),
				C.int(host.setting.client_button_width),
				C.int(host.setting.client_button_margin_width),
				C.int(host.mask_button),
				C.int(host.setting.client_text_margin_width),
				C.CString(clt.title),
				C.int(len(clt.title)),
				C.int(maximized_val))
	*/
	clt := host.client[address]
	cr := clt.mask.surface.wm_cairo_create_go_surface()
	
	attr := host.wm_host_get_window_attributes(clt.mask.window)
	surface_w := float64(attr.width)
	surface_h := float64(attr.height)

	cr.SetOperator(cairo.OPERATOR_CLEAR)
	cr.Paint()
	cr.SetOperator(cairo.OPERATOR_OVER)
	border_width := float64(host.setting.client_border_overall_width)
	shadow_width := float64(host.setting.client_border_shadow_width)
	
	button_width := float64(host.setting.client_button_width)
	button_margin_width := float64(host.setting.client_button_margin_width)
	mask_button := host.mask_button
	window_maximized := clt.maximize_mode == WM_CLIENT_MAXIMIZE_MODE_REVERSE
	
	text_margin_width := float64(host.setting.client_text_margin_width)

	shadow_roughness := shadow_width/5.0
	if shadow_roughness < 1 { shadow_roughness = 1 }

	var pattern_l [4]*cairo.Pattern
	for i := 0.0; i<4.0; i++{
		var linear cairo.Linear
		linear.X0 = 0
		linear.Y0 = 0
		linear.X1 = surface_w-shadow_width*2
		linear.Y1 = surface_h-shadow_width*2
		adj := i * 0.05;
		id := int(i)
		pattern_l[id] = cairo.NewPatternLinear(linear)
		pattern_l[id].AddColorStopRGB(0.0, 0.5-adj, 0.4-adj, 0.4-adj)
		pattern_l[id].AddColorStopRGB(1.0, 0.46-adj, 0.36-adj, 0.36-adj)

	}

	var rectangle_shadow = func(x float64, y float64, w float64, h float64){

		cr.SetLineWidth(shadow_roughness*0.25);

		for sh := 0.0; sh < shadow_width; sh += shadow_roughness{
			strength := (1.0 - sh/shadow_width)
			cr.SetSourceRGBA(0.0, 0.0, 0.0, 0.25 * strength * strength)
			cr.Rectangle(x-sh, y-sh, w+sh*2, h+sh*2)
			cr.Stroke()
		}
	}

	cr.SetSource(pattern_l[0])

	cr.Rectangle(shadow_width, shadow_width, surface_w-shadow_width*2, border_width-shadow_width-1)
	cr.Rectangle(shadow_width, surface_h-border_width-1, surface_w-shadow_width*2, border_width-shadow_width+1)
    cr.Rectangle(shadow_width, border_width-1, border_width-shadow_width-1, surface_h-border_width*2+1)
	cr.Rectangle(surface_w-border_width-1, border_width-1, border_width-shadow_width+1, surface_h-border_width*2+1)

	cr.Fill()
	rectangle_shadow(shadow_width, shadow_width,
		surface_w-shadow_width*2, surface_h-shadow_width*2)

    for i := 1.0; i <= 3.0; i++{

        button_y := border_width - button_width - button_margin_width
        button_x := surface_w - border_width - (button_width + button_margin_width)*i
        button_w := button_width;
        button_x_limit := border_width + button_margin_width

        if button_x < button_x_limit {
            button_w -= button_x_limit-button_x
            button_x = button_x_limit
        }

        if button_w < 0 {
			button_w = 0
		} else {
            rectangle_shadow(button_x, button_y, button_w, button_width)
        }

		cr.Rectangle(button_x, button_y, button_w, button_width)
		ad := 1+int(i)%2
		if mask_button == int(i) { ad++ }
		cr.SetSource(pattern_l[ad])
		cr.Fill()

		if button_w < button_width { continue }
        icon_margin := button_width*0.3
        icon_sx := button_x+icon_margin
        icon_sy := button_y+icon_margin
        icon_ex := button_x+button_width-icon_margin
        icon_ey := button_y+button_width-icon_margin
        icon_3_adjust := 0.7

        cr.SetLineWidth(button_width*0.1)
        cr.SetSourceRGB(0.9, 0.9, 0.9)

		switch i {
		case 1:
			cr.MoveTo(icon_sx, icon_sy)
			cr.LineTo(icon_ex, icon_ey)
			cr.MoveTo(icon_sx, icon_ey)
			cr.LineTo(icon_ex, icon_sy)
			cr.Stroke()
		case 2:
			cr.MoveTo(icon_sx, icon_ey)
			cr.LineTo(icon_ex, icon_ey)
			cr.Stroke()
		case 3:
			icon_w := icon_ex-icon_sx
			icon_h := icon_ey-icon_sy
			icon_margin := icon_w*(1.0-icon_3_adjust)
			if window_maximized {
				cr.SetSourceRGB(0.5, 0.5, 0.5)
				cr.Rectangle(icon_sx+icon_margin, icon_sy,
						icon_w*icon_3_adjust, icon_h*icon_3_adjust)
				cr.Stroke()
				cr.SetSourceRGB(0.9, 0.9, 0.9)
				cr.Rectangle(icon_sx, icon_sy+icon_margin,
						icon_w*icon_3_adjust, icon_h*icon_3_adjust)
				cr.Stroke()
			} else{
				cr.Rectangle(icon_sx, icon_sy, icon_w, icon_h)
				cr.Stroke()
			}
		}

	
	}

	cr.SelectFontFace("Serif",
		cairo.FONT_SLANT_NORMAL,
		cairo.FONT_WEIGHT_NORMAL)
	
	cr.SetFontSize(button_width)

	title := clt.title

	var extents *cairo.TextExtents
	extents = cr.TextExtents(title)
	textbox_x := border_width
    textbox_y := border_width - button_width - text_margin_width
    textbox_w := extents.Width + text_margin_width*2
    textbox_h := button_width + text_margin_width
    textbox_ex_limit := surface_w - border_width - button_margin_width - (button_width+button_margin_width)*3

    if textbox_x + textbox_w > textbox_ex_limit {
        textbox_w = textbox_ex_limit - textbox_x
        
        for i := len(title)-1 ; i >= 0 && textbox_w < extents.Width + text_margin_width*2; i--{
            title = title[:len(title)-1]
            extents = cr.TextExtents(title)
        }
    }

    if textbox_w < 0 {
		textbox_w = 0
	} else {
        rectangle_shadow(textbox_x, textbox_y,
                         textbox_w, textbox_h)
    }
    cr.Rectangle(textbox_x, textbox_y, textbox_w, textbox_h)
    cr.SetSource(pattern_l[0])
    cr.Fill()
    
    cr.MoveTo(textbox_x + text_margin_width, textbox_y + textbox_h - text_margin_width)
    cr.SetSourceRGB(0.9, 0.9, 0.9)
    cr.ShowText(title)

}

func (host *WmHost) wm_host_draw_background(w int, h int){

	cr := host.background.surface.wm_cairo_create_go_surface()
	img, status := cairo.NewSurfaceFromPNG(host.setting.background_file)
	if status != cairo.STATUS_SUCCESS {
		wm_debug_log("Couldn't Open PNG File :",host.setting.background_file)
		return
	}
	cr.Save()
	cr.Scale(float64(w)/float64(img.GetWidth()), float64(h)/float64(img.GetHeight()))
	cr.SetSourceSurface(img, 0, 0)
	cr.Paint()
	cr.Restore()
}

func (surface *CairoSfc) wm_cairo_create_context() *CairoCtx{
	return C.cairo_create(surface)
}

func (surface *CairoSfc) wm_cairo_create_go_surface() *cairo.Surface{
	raw_sfc := (*C.cairo_surface_t)(surface)
	raw_ctx := (*C.cairo_t)(surface.wm_cairo_create_context())
	return cairo.NewSurfaceFromC(
				cairo.Cairo_surface(unsafe.Pointer(raw_sfc)),
				cairo.Cairo_context(unsafe.Pointer(raw_ctx)))
}