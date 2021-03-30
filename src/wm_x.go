package main


/*
#cgo pkg-config: x11 cairo
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <cairo/cairo-xlib.h>
#include <X11/cursorfont.h>
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

	CairoSfc = C.cairo_surface_t
	CairoCtx = C.cairo_t

	CLong = C.long
)

const (
	XNone = C.None
	XKeyPress = int(C.KeyPress)
	XButtonPress = int(C.ButtonPress)
	XButtonRelease = int(C.ButtonRelease)
	XMotionNotify = int(C.MotionNotify)
	XMapNotify = int(C.MapNotify)
	XUnmapNotify = int(C.UnmapNotify)
	XMapRequest = int(C.MapRequest)
	XDestroyNotify = int(C.DestroyNotify)
	XConfigureNotify = int(C.ConfigureNotify)
	XConfigureRequest = int(C.ConfigureRequest)

	XSubstructureNotifyMask = CLong(C.SubstructureNotifyMask)
	XSubstructureRedirectMask = CLong(C.SubstructureRedirectMask)
	XButtonPressMask = CLong(C.ButtonPressMask)
	XButtonReleaseMask = CLong(C.ButtonReleaseMask)
	XPointerMotionMask = CLong(C.PointerMotionMask)

	XCLeftPtr = int(C.XC_left_ptr)
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
	//C.c_wm_x11_send_event(display, window, C.CString("WM_DELETE_WINDOW"))
}

func wm_x11_destroy_cairo_surface(display *XDisplay, surface *CairoSfc){
	C.cairo_surface_destroy(surface)
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

func wm_x11_draw_transparent(display *XDisplay, transparent WmTransparent){

	surface_w := C.cairo_image_surface_get_width(transparent.surface)
	surface_h := C.cairo_image_surface_get_height(transparent.surface)

	switch transparent.drawtype{
	case WM_DRAW_TYPE_BOX:
		C.c_wm_transparent_draw_type_box(transparent.surface, surface_w, surface_h)
	case WM_DRAW_TYPE_MASK:
		C.c_wm_transparent_draw_type_mask(transparent.surface, surface_w, surface_h)
	}

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