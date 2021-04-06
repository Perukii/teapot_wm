
package main

/*
#include "./c_wm_x_access.h"
*/
import "C"

import "unsafe"

func (host *WmHost) wm_host_init(){
	host.display = C.XOpenDisplay(nil)
	host.root_window = C.XDefaultRootWindow(host.display)
	host.client = []WmClient{}
	var clt WmClient
	host.client = append(host.client, clt)
}

func (host *WmHost) wm_host_init_log_file(){
	wm_debug_log_file_init(host.log_file)
}

func (host *WmHost) wm_host_close_log_file(){
	host.log_file.Close()
}

func (host *WmHost) wm_host_move_window(window XWindowID, x int, y int){
	C.XMoveWindow(host.display, window, C.int(x), C.int(y))
}

func (host *WmHost) wm_host_resize_window(window XWindowID, w int, h int){
	if w < 1 { w = 1 }
	if h < 1 { h = 1 }
	C.XResizeWindow(host.display, window, C.uint(w), C.uint(h))
}

func (host *WmHost) wm_host_move_resize_window(window XWindowID, x int, y int, w int, h int){
	if w < 1 { w = 1 }
	if h < 1 { h = 1 }
	C.XMoveResizeWindow(host.display, window, C.int(x), C.int(y), C.uint(w), C.uint(h))
}

func (host *WmHost) wm_host_resize_surface(surface *CairoSfc, w int, h int){
	if w < 1 { w = 1 }
	if h < 1 { h = 1 }
    C.cairo_xlib_surface_set_size(surface, C.int(w), C.int(h))
}

func (host *WmHost) wm_host_select_input(window XWindowID, mask CLong){
	C.XSelectInput(host.display, window, mask)
}

func (host *WmHost) wm_host_get_window_attributes(window XWindowID) XWindowAttributes{
	var attr XWindowAttributes
	C.XGetWindowAttributes(host.display, window, &attr)
	return attr
}

func (host *WmHost) wm_host_raise_window(window XWindowID){
	C.XRaiseWindow(host.display, window)
}

func (host *WmHost) wm_host_lower_window(window XWindowID){
	C.XLowerWindow(host.display, window)
}

func (host *WmHost) wm_host_setup_transparent(transparent *WmTransparent, parent XWindowID,
											  x int, y int, w int, h int,
											  ){
	var vinfo C.XVisualInfo
	C.XMatchVisualInfo(host.display, C.XDefaultScreen(host.display), 32, C.TrueColor, &vinfo)

	var attr C.XSetWindowAttributes
	attr.colormap = C.XCreateColormap(host.display, parent, vinfo.visual, C.AllocNone)
	attr.border_pixel = 0
	attr.background_pixel = 0
	attr.override_redirect = 1
	
	transparent.window = C.XCreateWindow(
		host.display,
		parent,
		C.int(x),
		C.int(y),
		C.uint(w),
		C.uint(h),
		C.uint(1),
		vinfo.depth,
		C.InputOutput,
		vinfo.visual,
		C.CWColormap|C.CWBorderPixel|C.CWBackPixel|C.CWOverrideRedirect,
		&attr,
	)
	
	transparent.surface = C.cairo_xlib_surface_create(
		host.display,
		transparent.window,
		vinfo.visual,
		C.int(w),
		C.int(h),
	)

}

func (host *WmHost) wm_host_remove_transparent(transparent WmTransparent){
	C.XDestroyWindow(host.display, transparent.window)
	C.cairo_surface_destroy(transparent.surface)
}

func (host *WmHost) wm_host_map_window(window XWindowID){
	C.XMapWindow(host.display, window)
}

func (host *WmHost) wm_host_unmap_window(window XWindowID){
	C.XUnmapWindow(host.display, window)
}

func (host *WmHost) wm_host_reparent_window(window XWindowID, parent XWindowID, x int, y int){
	C.XReparentWindow(host.display, window, parent ,C.int(x), C.int(y))
}

func (host *WmHost) wm_host_define_cursor(cursor int){
	if host.cursor == cursor { return }
	host.cursor = cursor
    C.XDefineCursor(host.display, host.root_window,
		C.XCreateFontCursor(host.display, C.uint(cursor)));

}

func (host *WmHost) wm_host_query_parent(window XWindowID) XWindowID{
	relation := wm_x11_query_tree(host.display, window)
	return relation.parent
}

func (host *WmHost) wm_host_send_delete_event(window XWindowID){
	C.c_wm_x11_send_event_destroy(host.display, window)
}

func (host *WmHost) wm_host_check_n_of_queued_event() int{
	return int(C.XEventsQueued(host.display, C.QueuedAlready))
}

func (host *WmHost) wm_host_set_focus_to_client(address WmClientAddress){

	host.wm_client_raise_mask(len(host.client)-1)

	clt := host.client[address]
	for i := address; i < len(host.client)-1; i++{
		host.client[i] = host.client[i+1]
	}
	host.client[len(host.client)-1] = clt

	host.wm_client_raise_app(len(host.client)-1)

}

func (host *WmHost) wm_host_restack_clients(){
	for i := 1; i < len(host.client)-1; i++{
		host.wm_client_raise_mask(i)
	}
	host.wm_client_raise_app(len(host.client)-1)
}

func (host *WmHost) wm_host_get_size_hints(window XWindowID) XSizeHints{
	var hints XSizeHints
	var supplied C.long
	C.XGetWMNormalHints(host.display, window, &hints, &supplied)
	return hints
}

func (host *WmHost) wm_host_get_window_title(window XWindowID) string{
	name := C.c_wm_x11_get_window_title(host.display, window)
	name_res := C.GoString(name)
	if name != nil { C.free(unsafe.Pointer(name)) }
	return name_res
}

func (host *WmHost) wm_host_intern_atom(name string) XAtom{
	return C.XInternAtom(host.display, C.CString(name), 1)
}

func (host *WmHost) wm_host_update_grab_mode(point_x int, point_y int, mask_x int, mask_y int, mask_w int, mask_h int){

	resize_area_width := host.setting.client_border_shadow_width

	grab_rx := point_x-mask_x
	grab_ry := point_y-mask_y
	
	host.grab_mode_1 = WM_RESIZE_MODE_NONE

	if grab_rx < resize_area_width  {
		host.grab_mode_1 = WM_RESIZE_MODE_LEFT
	}
	if grab_rx > mask_w-resize_area_width {
		host.grab_mode_1 = WM_RESIZE_MODE_RIGHT
	}

	host.grab_mode_2 = WM_RESIZE_MODE_NONE

	if grab_ry < resize_area_width  {
		host.grab_mode_2 = WM_RESIZE_MODE_TOP
	}
	if grab_ry > mask_h-resize_area_width {
		host.grab_mode_2 = WM_RESIZE_MODE_BOTTOM
	}
}

func (host *WmHost) wm_host_update_button_mode(point_x int, point_y int, mask_x int, mask_y int, mask_w int, mask_h int){
	var in_rect = func(px, py, x, y, w, h int) bool{
		return (px >= x) && (px <= x+w) && (py >= y) && (py <= y+h)
	}

	host.mask_button = WM_BUTTON_NONE

	border_width := host.setting.client_border_overall_width
	button_width := host.setting.client_button_width
	button_margin_width := host.setting.client_button_margin_width

	for i := 1; i <= 3; i++ {
        
        button_x := mask_w - border_width - (button_width + button_margin_width)*i;
		button_y := border_width - button_width - button_margin_width;

		if in_rect(point_x-mask_x, point_y-mask_y, button_x, button_y, button_width, button_width){
			switch i {
			case 1:
				host.mask_button = WM_BUTTON_EXIT
			case 2:
				host.mask_button = WM_BUTTON_MINIMIZE
			case 3:
				host.mask_button = WM_BUTTON_MAXIMIZE
			}
			break
		}
	}
}

func (host *WmHost) wm_host_update_cursor(){
	if host.grab_mode_1 == WM_RESIZE_MODE_NONE && host.grab_mode_2 == WM_RESIZE_MODE_NONE {
		host.wm_host_define_cursor(XCLeftPtr)
	}
	mode_l := host.grab_mode_1 == WM_RESIZE_MODE_LEFT
	mode_r := host.grab_mode_1 == WM_RESIZE_MODE_RIGHT
	mode_t := host.grab_mode_2 == WM_RESIZE_MODE_TOP
	mode_b := host.grab_mode_2 == WM_RESIZE_MODE_BOTTOM

	if mode_l && mode_t { host.wm_host_define_cursor(XCSideTL);return }
	if mode_r && mode_t { host.wm_host_define_cursor(XCSideTR);return }
	if mode_l && mode_b { host.wm_host_define_cursor(XCSideBL);return }
	if mode_r && mode_b { host.wm_host_define_cursor(XCSideBR);return }

	if mode_t { host.wm_host_define_cursor(XCSideT);return } 
	if mode_b { host.wm_host_define_cursor(XCSideB);return } 
	if mode_l { host.wm_host_define_cursor(XCSideL);return } 
	if mode_r { host.wm_host_define_cursor(XCSideR);return } 
}


func (host *WmHost) wm_host_run(){
	for{
		C.XNextEvent(host.display, &host.event)
		switch host.event.wm_event_get_type(){
		case XMapNotify:
			host.wm_event_loop_map_notify()
		case XUnmapNotify:
			host.wm_event_loop_unmap_notify()
		case XMapRequest:
			host.wm_event_loop_map_request()
		case XDestroyNotify:
			host.wm_event_loop_destroy_notify()
		case XConfigureNotify:
			host.wm_event_loop_configure_notify()
		case XConfigureRequest:
			host.wm_event_loop_configure_request()
		case XKeyRelease:
			fallthrough;
		case XKeyPress:
			host.wm_event_loop_key_press()
		case XButtonPress:
			host.wm_event_loop_button_press()
		case XButtonRelease:
			host.wm_event_loop_button_release()
		case XMotionNotify:
			host.wm_event_loop_motion_notify()
		case XEnterNotify:
			host.wm_event_loop_enter_notify()
		case XLeaveNotify:
			host.wm_event_loop_leave_notify()
		case XPropertyNotify:
			host.wm_event_loop_property_notify()
		}
	}
}

