package main


/*
#cgo pkg-config: x11 cairo
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <cairo/cairo-xlib.h>
#include <X11/cursorfont.h>
#include <stdlib.h>
#include "./c_wm_draw.h"
#include "./c_wm_x_access.h"
*/
import "C"

import "unsafe"

type (
	XEvent = C.XEvent
	XDisplay = C.Display
	XWindowID = C.Window
	XWindowAttributes = C.XWindowAttributes
	XKeyEvent = C.XKeyEvent
	XButtonEvent = C.XButtonEvent
	XMotionEvent = C.XMotionEvent
	XMapEvent = C.XMapEvent
	XUnmapEvent = C.XUnmapEvent
	XMapRequestEvent = C.XMapRequestEvent
	XDestroyWindowEvent = C.XDestroyWindowEvent
	XConfigureEvent = C.XConfigureEvent
	XConfigureRequestEvent = C.XConfigureRequestEvent
	XCrossingEvent = C.XCrossingEvent
	XSizeHints = C.XSizeHints
	XPropertyEvent = C.XPropertyEvent

	XAtom = C.Atom

	CairoSfc = C.cairo_surface_t
	CairoCtx = C.cairo_t

	CLong = C.long
)

const (
	XNone = C.None
	XKeyPress = int(C.KeyPress)
	XKeyRelease = int(C.KeyRelease)
	XButtonPress = int(C.ButtonPress)
	XButtonRelease = int(C.ButtonRelease)
	XMotionNotify = int(C.MotionNotify)
	XMapNotify = int(C.MapNotify)
	XUnmapNotify = int(C.UnmapNotify)
	XMapRequest = int(C.MapRequest)
	XDestroyNotify = int(C.DestroyNotify)
	XConfigureNotify = int(C.ConfigureNotify)
	XConfigureRequest = int(C.ConfigureRequest)
	XEnterNotify = int(C.EnterNotify)
	XLeaveNotify = int(C.LeaveNotify)
	XPropertyNotify = int(C.PropertyNotify)

	XSubstructureNotifyMask = CLong(C.SubstructureNotifyMask)
	XSubstructureRedirectMask = CLong(C.SubstructureRedirectMask)
	XButtonPressMask = CLong(C.ButtonPressMask)
	XButtonReleaseMask = CLong(C.ButtonReleaseMask)
	XPointerMotionMask = CLong(C.PointerMotionMask)
	XEnterWindowMask = CLong(C.EnterWindowMask)
	XLeaveWindowMask = CLong(C.LeaveWindowMask)
	XPropertyChangeMask = CLong(C.PropertyChangeMask)
	XKeyPressMask = CLong(C.KeyPressMask)
	XKeyReleaseMask = CLong(C.KeyReleaseMask)

	XCLeftPtr = int(C.XC_left_ptr)
	XCSideT = int(C.XC_top_side)
	XCSideB = int(C.XC_bottom_side)
	XCSideL = int(C.XC_left_side)
	XCSideR = int(C.XC_right_side)
	XCSideTL = int(C.XC_top_left_corner)
	XCSideTR = int(C.XC_top_right_corner)
	XCSideBL = int(C.XC_bottom_left_corner)
	XCSideBR = int(C.XC_bottom_right_corner)

)

func wm_x11_open_display() *XDisplay{
	return C.XOpenDisplay(nil)
}

func wm_x11_get_root(display *XDisplay) XWindowID{
	return C.XDefaultRootWindow(display)
}

func wm_x11_peek_event(display *XDisplay) XEvent{
	var event XEvent
	C.XNextEvent(display, &event)
	return event
}

func wm_x11_select_input(display *XDisplay, window XWindowID, mask CLong){
	C.XSelectInput(display, window, mask)
}

func wm_x11_grab_button(display *XDisplay, window XWindowID){
    C.XGrabButton(display, C.AnyButton, C.Mod1Mask, window, C.True,
		C.ButtonPressMask|C.ButtonReleaseMask|C.PointerMotionMask,
		C.GrabModeAsync, C.GrabModeAsync, C.None, C.None)
}

func wm_x11_grab_key(display *XDisplay, window XWindowID, keycode int){
	C.XGrabKey(display, C.int(keycode), C.Mod1Mask,
				window, C.True, C.GrabModeAsync, C.GrabModeAsync)
}

func wm_x11_get_type_of_event(event *XEvent) int {
	return int(C.c_wm_x11_get_type_of_event(*event))
}

func wm_x11_raise_window(display *XDisplay, window XWindowID){
	C.XRaiseWindow(display, window)
}

func wm_x11_lower_window(display *XDisplay, window XWindowID){
	C.XLowerWindow(display, window)
}

func wm_x11_set_input_focus(display *XDisplay, window XWindowID){
	C.XSetInputFocus(display, window, C.RevertToNone, C.CurrentTime)
}

func wm_x11_get_window_attributes(display *XDisplay, window XWindowID) XWindowAttributes{
	var attr XWindowAttributes
	C.XGetWindowAttributes(display, window, &attr)
	return attr
}

func wm_x11_move_window(display *XDisplay, window XWindowID, x int, y int){
	C.XMoveWindow(display, window, C.int(x), C.int(y))
}

func wm_x11_resize_window(display *XDisplay, window XWindowID, w int, h int){
	if w < 1 { w = 1 }
	if h < 1 { h = 1 }
	C.XResizeWindow(display, window, C.uint(w), C.uint(h))
}

func wm_x11_move_resize_window(display *XDisplay, window XWindowID, x int, y int, w int, h int){
	if w < 1 { w = 1 }
	if h < 1 { h = 1 }
	C.XMoveResizeWindow(display, window, C.int(x), C.int(y), C.uint(w), C.uint(h))
}


func wm_x11_resize_surface(surface *CairoSfc, w int, h int){
	if w < 1 { w = 1 }
	if h < 1 { h = 1 }
    C.cairo_xlib_surface_set_size(surface, C.int(w), C.int(h))
}

func wm_x11_map_window(display *XDisplay, window XWindowID){
	C.XMapWindow(display, window)
} 

func wm_x11_unmap_window(display *XDisplay, window XWindowID){
	C.XUnmapWindow(display, window)
} 

func wm_x11_reparent_window(display *XDisplay, window XWindowID, parent XWindowID, x int, y int){
	C.XReparentWindow(display, window, parent ,C.int(x), C.int(y))
}

func wm_x11_destroy_window(display *XDisplay, window XWindowID){
	C.XDestroyWindow(display, window)
}

func wm_x11_send_delete_event(display *XDisplay, window XWindowID){
	C.c_wm_x11_send_event_destroy(display, window)
}

func wm_x11_destroy_cairo_surface(display *XDisplay, surface *CairoSfc){
	C.cairo_surface_destroy(surface)
}

func wm_x11_intern_atom(display *XDisplay, name string) XAtom{
	return C.XInternAtom(display, C.CString(name), 1)
}

func wm_x11_create_transparent_window(display *XDisplay, parent XWindowID,
									  x int, y int, w int, h int,
									  transparent *WmTransparent,
									  ){
	
	var vinfo C.XVisualInfo
	C.XMatchVisualInfo(display, C.XDefaultScreen(display), 32, C.TrueColor, &vinfo)

	var attr C.XSetWindowAttributes
	attr.colormap = C.XCreateColormap(display, parent, vinfo.visual, C.AllocNone)
    attr.border_pixel = 0
    attr.background_pixel = 0
	attr.override_redirect = 1
	
	transparent.window = C.XCreateWindow(
		display,
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
		display,
		transparent.window,
		vinfo.visual,
		C.int(w),
		C.int(h),
	)

}

func wm_x11_draw_box(display *XDisplay, transparent WmTransparent,
					setting WmSetting, mask_button int, title string, maximized bool){

	attr := wm_x11_get_window_attributes(display, transparent.window)
	surface_w := attr.width
	surface_h := attr.height

	var maximized_val int = 0
	if maximized { maximized_val = 1 }

	C.c_wm_transparent_draw_type_box(transparent.surface, surface_w, surface_h,
										C.int(setting.client_border_overall_width),
										C.int(setting.client_border_shadow_width),
										C.int(setting.client_button_width),
										C.int(setting.client_button_margin_width),
										C.int(mask_button),
										C.int(setting.client_text_margin_width),
										C.CString(title),
										C.int(len(title)),
										C.int(maximized_val))
}

func wm_x11_define_cursor(display *XDisplay, window XWindowID, cursor int){
    C.XDefineCursor(display, window,
		C.XCreateFontCursor(display, C.uint(cursor)));
}

func wm_x11_query_tree(display *XDisplay, window XWindowID) WmWindowRelation{
	var target WmWindowRelation
	var children *XWindowID
	var nc C.uint
	C.XQueryTree(display, window, &target.root_window, &target.parent, &children, &nc)

	target.children = make([]XWindowID, int(nc), int(nc))

	for i := 0; i < int(nc); i++ {
		target.children[i] = C.c_wm_x11_query_window_from_array(children, C.int(i))
	}

	if children != nil {
		C.XFree(unsafe.Pointer(children))
	}

	return target
}

func wm_x11_check_n_of_queued_event(display *XDisplay) int{
	return int(C.XEventsQueued(display, C.QueuedAlready))
}

func wm_x11_get_size_hints(display *XDisplay, window XWindowID) XSizeHints{
	var hints XSizeHints
	var supplied C.long
	C.XGetWMNormalHints(display, window, &hints, &supplied)
	return hints
}

func wm_x11_get_window_title(display *XDisplay, window XWindowID) string{
	
	name := C.c_wm_x11_get_window_title(display, window)
	name_res := C.GoString(name)
	if name != nil { C.free(unsafe.Pointer(name)) }
	return name_res

}