package main


/*
#cgo pkg-config: x11 cairo
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <cairo/cairo-xlib.h>
#include <X11/cursorfont.h>
#include "./c_wm_x_access.h"
*/
import "C"

type (
	XEvent = C.XEvent
	XDisplay = C.Display
	XWindowID = C.Window
	XWindowAttributes = C.XWindowAttributes
	XKeyEvent = C.XKeyEvent
	XButtonEvent = C.XButtonEvent
	XMotionEvent = C.XMotionEvent

	CairoSfc = C.cairo_surface_t
	CairoCtx = C.cairo_t
)

const (
	XNone = C.None
	XKeyPress = int(C.KeyPress)
	XButtonPress = int(C.ButtonPress)
	XButtonRelease = int(C.ButtonRelease)
	XMotionNotify = int(C.MotionNotify)
	XMapNotify = int(C.MapNotify)
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

func wm_x11_map_window(display *XDisplay, window XWindowID){
	C.XMapWindow(display, window)
} 

func wm_x11_reparent_window(display *XDisplay, window XWindowID, parent XWindowID, x int, y int){
	C.XReparentWindow(display, window, parent ,C.int(x), C.int(y))
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
		//C.CWColormap|C.CWBackPixel|C.CWOverrideRedirect,
		&attr,
	)
	
	
/*
	transparent.window = C.XCreateSimpleWindow(
		display,
		parent,
		C.int(x),
		C.int(y),
		C.uint(w),
		C.uint(h),
		C.uint(1),
		C.ulong(0),
		C.XBlackPixel(display, 0),
	)
*/

	transparent.surface = C.cairo_xlib_surface_create(
		display,
		transparent.window,
		vinfo.visual,
		C.int(w),
		C.int(h),
	)

}

func wm_x11_draw_transparent(display *XDisplay, transparent WmTransparent){

	transparent.ctx = C.cairo_create(transparent.surface)
	C.cairo_set_source_rgba(transparent.ctx, C.double(0.5), C.double(0.5), C.double(1), C.double(0.5))
	C.cairo_paint(transparent.ctx)
	C.cairo_surface_flush(transparent.surface)

}